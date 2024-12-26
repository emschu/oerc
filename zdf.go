// oerc, alias oer-collector
// Copyright (C) 2021-2024 emschu[aet]mailbox.org
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public
// License along with this program.
// If not, see <https://www.gnu.org/licenses/>.
package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/alitto/pond"
	"log"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	zdfHost              = "https://www.zdf.de"
	zdfAPIHost           = "https://api.zdf.de"
	zdfAPIKeyPath        = "/live-tv"
	zdfAPIProgramEntries = zdfAPIHost + "/cmdm/epg/broadcasts?limit=100&order=asc"
)

var (
	pendingHashes              sync.Map // to avoid duplicate entries
	zdfTvShowLinkMatcher       = regexp.MustCompile(`^https?://(www\.)?zdf.de/.*`)
	zdfTvShowExternalIDMatcher = regexp.MustCompile(`[a-zA-Z0-9_-]+`)
)

// ZDFParser struct to group zdf parsing code
type ZDFParser struct {
	ParserInterface
	Parser
	zdfAPIKey string
}

func (z *ZDFParser) postProcess() {}

func (z *ZDFParser) preProcess() bool {
	var apiKeyErr error
	zdfAPIKey, apiKeyErr := z.getZdfAPIKey()
	if apiKeyErr != nil {
		log.Printf("Error fetching zdf api key: %v\n", apiKeyErr)
		return false
	}
	if verboseGlobal {
		log.Printf("Using ZDF API key: %s\n", *zdfAPIKey)
	}
	z.zdfAPIKey = *zdfAPIKey

	z.parallelWorkersCount = 3
	return true
}

// getZdfAPIKey method to retrieve the api key we need to connect to the zdf api
func (z *ZDFParser) getZdfAPIKey() (*string, error) {
	apiURL := fmt.Sprintf("%s%s", zdfHost, zdfAPIKeyPath)

	doc, err := getDocument(apiURL)
	if doc == nil || err != nil {
		return nil, fmt.Errorf("problem fetching zdf url '%s'", apiURL)
	}
	var apiToken string
	doc.Find("script").Each(func(i int, selection *goquery.Selection) {
		html, _ := selection.Html()
		if strings.Contains(html, "IMPORTANT CONFIGURATION!") {
			extractionPattern := regexp.MustCompile(`apiToken: &#39;(.*)&#39;`)
			findString := extractionPattern.FindAllStringSubmatch(html, 1)
			if len(findString) > 0 && len(findString[0]) > 1 {
				apiToken = trimAndSanitizeString(findString[0][1])
			}
		}
	})
	if len(apiToken) > 0 {
		return &apiToken, nil
	}
	return nil, fmt.Errorf("can't fetch ZDF api key")
}

// method to process a single day of a single channel
func (z *ZDFParser) handleDay(channel Channel, day time.Time) {
	var startDateStr = day.Format(time.RFC3339)
	var endDate = day.AddDate(0, 0, 1)
	var endDateStr = endDate.Format(time.RFC3339)

	// TODO pagination?
	apiURL := fmt.Sprintf(
		"%s&tvServices=%s&from=%s&to=%s&page=1",
		zdfAPIProgramEntries,
		url.QueryEscape(channel.TechnicalID),
		url.QueryEscape(startDateStr),
		url.QueryEscape(endDateStr),
	)
	if verboseGlobal {
		log.Printf("Visit: %s\n", apiURL)
	}

	zdfProgramRequest, apiErr := z.doZDFApiBroadcastRequest(apiURL)
	if apiErr != nil {
		log.Printf("%v\n", apiErr)
		return
	}

	for _, broadcast := range zdfProgramRequest.BroadCasts {
		var programEntry ProgramEntry

		hash := buildHash([]string{
			broadcast.PosID,
			broadcast.AirtimeBegin.Format(time.RFC3339),
			broadcast.PlayoutID,
			broadcast.Title,
			broadcast.TvServiceID,
		})
		z.db.Model(ProgramEntry{}).
			Where("hash = ?", hash).
			Where("channel_id = ?", channel.ID).
			Preload("ImageLinks").
			Find(&programEntry)
		programEntry.Hash = hash
		if programEntry.ID >= 0 && programEntry.isRecentlyUpdated() {
			atomic.AddUint64(&status.TotalSkippedPE, 1)
			continue
		}
		_, ok := pendingHashes.Load(hash)
		if ok {
			// the same hash is already getting added
			continue
		}
		pendingHashes.Store(hash, true)

		programEntry.TechnicalID = broadcast.PosID
		begin := broadcast.EffectiveAirtimeBegin
		end := broadcast.EffectiveAirtimeEnd

		// Use them as fallback...
		if begin.IsZero() {
			begin = broadcast.AirtimeBegin
		}
		if end.IsZero() {
			end = broadcast.AirtimeEnd
		}
		if begin.IsZero() || end.IsZero() {
			appLog(fmt.Sprintf("Error fetching zdf program entry at POSID '%s'. Invalid airtime begin or end date.", broadcast.PosID))
			continue
		}

		programEntry.StartDateTime = &begin
		programEntry.EndDateTime = &end
		programEntry.DurationMinutes = int16(programEntry.EndDateTime.Sub(*programEntry.StartDateTime).Minutes())

		programEntry.URL = fmt.Sprintf("%s%s", zdfAPIHost, broadcast.Self)
		programEntry.Title = trimAndSanitizeString(broadcast.Title)

		var subTitle = trimAndSanitizeString(broadcast.SubTitle)
		if len(broadcast.SubTitle) > 0 {
			// append to title
			programEntry.Title += " - " + subTitle
		}

		programEntry.Description = trimAndSanitizeString(broadcast.Text)
		programEntry.Channel = channel
		programEntry.ChannelFamily = z.ChannelFamily

		// save image links
		programEntry.handleProgramImageLinks(&broadcast)

		programItemAPIUrl := zdfAPIHost + broadcast.ProgrammeItem
		programEntry.Homepage = programItemAPIUrl
		apiResponse, apiErr := z.doZDFApiProgramItemRequest(programItemAPIUrl)
		if apiErr != nil {
			log.Printf("Problem fetching zdf api program item data: '%s'\n", programItemAPIUrl)
		} else {
			if len(apiResponse.Category) > 0 {
				programEntry.considerTagExists(&apiResponse.Category)
			}
			if len(apiResponse.Genre) > 0 {
				programEntry.considerTagExists(&apiResponse.Genre)
			}
		}

		programEntry.saveProgramEntryRecord(z.db)
		pendingHashes.Delete(programEntry.Hash)
	}
}

// atm only useful in zdf context
func (p *ProgramEntry) handleProgramImageLinks(broadcast *ZdfBroadcast) {
	if broadcast == nil || p == nil {
		// nothing to do here
		return
	}
	if len(broadcast.Images.Layouts.W2400) > 0 && !p.doesImageLinkExist(broadcast.Images.Layouts.W2400) {
		p.ImageLinks = append(p.ImageLinks, ImageLink{URL: broadcast.Images.Layouts.W2400})
	}
	if len(broadcast.Images.Layouts.W1920) > 0 && !p.doesImageLinkExist(broadcast.Images.Layouts.W1920) {
		p.ImageLinks = append(p.ImageLinks, ImageLink{URL: broadcast.Images.Layouts.W1920})
	}
	if len(broadcast.Images.Layouts.W1280) > 0 && !p.doesImageLinkExist(broadcast.Images.Layouts.W1280) {
		p.ImageLinks = append(p.ImageLinks, ImageLink{URL: broadcast.Images.Layouts.W1280})
	}
	if len(broadcast.Images.Layouts.W768) > 0 && !p.doesImageLinkExist(broadcast.Images.Layouts.W768) {
		p.ImageLinks = append(p.ImageLinks, ImageLink{URL: broadcast.Images.Layouts.W768})
	}
	if len(broadcast.Images.Layouts.W640) > 0 && !p.doesImageLinkExist(broadcast.Images.Layouts.W640) {
		p.ImageLinks = append(p.ImageLinks, ImageLink{URL: broadcast.Images.Layouts.W640})
	}
	if len(broadcast.Images.Layouts.W384) > 0 && !p.doesImageLinkExist(broadcast.Images.Layouts.W384) {
		p.ImageLinks = append(p.ImageLinks, ImageLink{URL: broadcast.Images.Layouts.W384})
	}
	if len(broadcast.Images.Layouts.W276) > 0 && !p.doesImageLinkExist(broadcast.Images.Layouts.W276) {
		p.ImageLinks = append(p.ImageLinks, ImageLink{URL: broadcast.Images.Layouts.W276})
	}
	if len(broadcast.Images.Layouts.W240) > 0 && !p.doesImageLinkExist(broadcast.Images.Layouts.W240) {
		p.ImageLinks = append(p.ImageLinks, ImageLink{URL: broadcast.Images.Layouts.W240})
	}
}

func (z *ZDFParser) doZDFApiBroadcastRequest(apiURL string) (*ZdfBroadcastResponse, error) {
	headers := map[string]string{
		"Host":     "api.zdf.de",
		"Accept":   "application/vnd.de.zdf.v1.0+json",
		"Origin":   zdfHost,
		"Api-Auth": "Bearer " + z.zdfAPIKey,
	}
	resp, err := doGetRequest(apiURL, headers, 3)
	if resp == nil || err != nil {
		errMsg := fmt.Sprintf("Problem fetching URL '%s' with error '%v'", apiURL, err)
		appLog(errMsg)
		log.Println(errMsg)
		return nil, err
	}
	var response ZdfBroadcastResponse
	jsonErr := json.Unmarshal([]byte(*resp), &response)
	if jsonErr != nil {
		errMsg := fmt.Sprintf("Invalid json format in zdf api response. url: '%s'", apiURL)
		appLog(errMsg)
		log.Println(errMsg)
		return nil, jsonErr
	}
	return &response, nil
}

func (z *ZDFParser) doZDFApiProgramItemRequest(apiURL string) (*ZdfProgramItemResponse, error) {
	headers := map[string]string{
		"Host":     "api.zdf.de",
		"Accept":   "application/vnd.de.zdf.v1.0+json",
		"Origin":   zdfHost,
		"Api-Auth": "Bearer " + z.zdfAPIKey,
	}
	resp, err := doGetRequest(apiURL, headers, 3)
	if resp == nil || err != nil {
		errMsg := fmt.Sprintf("Problem fetching URL '%s' with error '%v'", apiURL, err)
		appLog(errMsg)
		log.Println(errMsg)
		return nil, err
	}
	var response ZdfProgramItemResponse
	jsonErr := json.Unmarshal([]byte(*resp), &response)
	if jsonErr != nil {
		errMsg := fmt.Sprintf("Invalid json format: %s", jsonErr.Error())
		appLog(errMsg)
		log.Println(errMsg)
		return nil, jsonErr
	}
	return &response, nil
}

// fetchTvShows method to fetch zdf tv shows
func (z *ZDFParser) fetchTVShows() {
	if !GetAppConf().EnableTVShowCollection || isRecentlyFetched() {
		z.logRecentFetch("Skip update of zdf tv shows")
		return
	}
	var tvShowGroups = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "0+-+9"}

	// fetch all links to tv shows
	if verboseGlobal {
		log.Printf("Collecting zdf tv shows ...\n")
	}
	var tvShowLinks = make([]string, 0)
	for _, group := range tvShowGroups {
		apiURL := fmt.Sprintf("%s/sendungen-a-z?group=%s", zdfHost, group)

		document, err := getDocument(apiURL)
		if document == nil || err != nil {
			appLog(fmt.Sprintf("Problem with http call to zdf %v", err))
			continue
		}
		document.Find(".b-content-teaser-item h3 a").Each(func(i int, selection *goquery.Selection) {
			tvShowPage := selection.AttrOr("href", "")
			if len(tvShowPage) > 0 {
				tvShowLinks = append(tvShowLinks, zdfHost+tvShowPage)
			}
		})
	}

	// and now process them
	if verboseGlobal {
		log.Printf("Processing %d zdf tv shows ...\n", len(tvShowLinks))
	}

	pool := pond.New(4, 100, getWorkerPoolIdleTimeout())
	for _, singleTvShowPage := range tvShowLinks {
		singlePage := singleTvShowPage

		pool.Submit(func() {
			z.processSingleTvShow(singlePage)
		})
	}
	pool.StopAndWait()
}

func (z *ZDFParser) processSingleTvShow(singleTvShowPage string) {
	doc, err := getDocument(singleTvShowPage)
	if doc == nil || err != nil {
		appLog(fmt.Sprintf("Could not fetch tv show detail page at '%s'", singleTvShowPage))
		return
	}
	// var tvShowRecord TvShow
	plusBtn := doc.Find(".b-plus-button")

	var tvShowRecord TvShow
	tvShowTitle, exists := plusBtn.Attr("data-plusbar-title")
	tvShowURL, exists2 := plusBtn.Attr("data-plusbar-url")
	tvShowAPIPath, exists3 := plusBtn.Attr("data-plusbar-path")
	tvShowExternalID, exists4 := plusBtn.Attr("data-plusbar-external-id")
	tvShowID, exists5 := plusBtn.Attr("data-plusbar-id")

	if !exists || !exists2 || !exists3 || !exists4 || !exists5 {
		appLog(fmt.Sprintf("Error finding zdf tv show information"))
		return
	}
	hash := buildHash([]string{tvShowExternalID, tvShowID})
	z.db.Where("hash = ?", hash).Find(&tvShowRecord)

	// do validation steps
	if !zdfTvShowLinkMatcher.Match([]byte(tvShowURL)) {
		appLog(fmt.Sprintf("Unexpected zdf tv show url '%s' detected. URL was set empty.", tvShowURL))
		tvShowURL = ""
	}
	if !zdfTvShowExternalIDMatcher.Match([]byte(tvShowExternalID)) {
		appLog("Unexpected external id received for zdf tv show. Skipping entry.")
		return
	}
	tvShowRecord.Hash = hash
	tvShowRecord.Title = trimAndSanitizeString(tvShowTitle)
	tvShowRecord.Homepage = tvShowURL
	tvShowRecord.URL = zdfHost + tvShowAPIPath
	tvShowRecord.TechnicalID = tvShowExternalID
	tvShowRecord.ChannelFamily = z.ChannelFamily

	tvShowRecord.saveTvShowRecord(z.db)
	return
}

func (z *ZDFParser) isDateValidToFetch(day *time.Time) (bool, error) {
	if day == nil {
		return false, fmt.Errorf("invalid day")
	}

	if z.isMoreThanXDaysInFuture(day, 43) { // = six weeks in future + today
		return false, fmt.Errorf("maximum for days in future for ZDF is 43")
	}
	location, _ := time.LoadLocation(GetAppConf().TimeZone)
	earliestDate := time.Date(2015, 2, 1, 0, 0, 0, 0, location)
	if day.Before(earliestDate) {
		return false, fmt.Errorf("maximum for days in past for ZDF is %s", earliestDate.Format(time.RFC822))
	}
	return true, nil
}

// ZdfBroadcastResponse api response struct definitions
type ZdfBroadcastResponse struct {
	BroadCasts []ZdfBroadcast `json:"http://zdf.de/rels/cmdm/broadcasts"`
}

// ZdfBroadcast zdf api object
type ZdfBroadcast struct {
	PosID                 string    `json:"posId"`
	PlayoutID             string    `json:"playoutId"`
	AirtimeBegin          time.Time `json:"airtimeBegin"`
	AirtimeEnd            time.Time `json:"airtimeEnd"`
	EffectiveAirtimeBegin time.Time `json:"effectiveAirtimeBegin"`
	EffectiveAirtimeEnd   time.Time `json:"effectiveAirtimeEnd"`
	Self                  string    `json:"self"`
	Title                 string    `json:"title"`
	SubTitle              string    `json:"subtitle"`
	Text                  string    `json:"text"`
	Images                ZdfImage  `json:"http://zdf.de/rels/image"`
	TvServiceID           string    `json:"tvServiceId"`
	ProgrammeItem         string    `json:"http://zdf.de/rels/cmdm/programme-item"`
}

// ZdfImage zdf api object
type ZdfImage struct {
	Source  string          `json:"source"`
	Layouts ZdfImageLayouts `json:"layouts"`
}

// ZdfImageLayouts zdf api object
type ZdfImageLayouts struct {
	W2400 string `json:"2400x1350,omitempty"`
	W1920 string `json:"1920x1080,omitempty"`
	W1280 string `json:"1280x720,omitempty"`
	W1152 string `json:"1152x1296,omitempty"`
	W768  string `json:"768x432,omitempty"`
	W640  string `json:"640x720,omitempty"`
	W384  string `json:"384x216,omitempty"`
	W276  string `json:"276x155,omitempty"`
	W240  string `json:"240x270,omitempty"`
}

// ZdfProgramItemResponse zdf api object
type ZdfProgramItemResponse struct {
	Category string `json:"category,omitempty"`
	Genre    string `json:"genre,omitempty"`
}
