// oerc, alias oer-collector
// Copyright (C) 2021-2026 emschu[aet]mailbox.org
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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"regexp"
	"strings"
	"sync/atomic"
	"time"
)

const (
	ardHost                   = "programm-api.ard.de"
	ardHostWithPrefix         = "https://" + ardHost
	ardMediathekAPIHost       = "api.ardmediathek.de"
	ardMediathekAPITvShowPath = "https://" + ardMediathekAPIHost + "/page-gateway/widgets/ard/editorials/"
	ardDefaultImageWidth      = "600"
)

var (
	ardProgramPageMatcher = regexp.MustCompile(`^https:\/\/programm\-api\.ard\.de\/program\/api\/teaser\?teaserId=.+`)
	ardMediathekMatcher   = regexp.MustCompile(`^https:\/\/www.ardmediathek\.de\/video\/.+`)
	ardImageLinkMatcher   = regexp.MustCompile(`^https:\/\/api\.ardmediathek\.de\/image-service\/image.+`)
	ardImageLinkMatcher2  = regexp.MustCompile(`^https:\/\/programm-api.ard.de\/images\/.+`)
	ardTvShowCategories   = []string{
		"a",
		"b",
		"c",
		"d",
		"e",
		"f",
		"g",
		"h",
		"i",
		//"j", empty 26/04
		"k",
		"l",
		"m",
		"n",
		"o",
		"p",
		"q",
		"r",
		"s",
		"t",
		"u",
		"v",
		"w",
		//"x", empty 26/04
		"y",
		"z",
		"#",
	}
)

// ARDParser struct of ard parser code
type ARDParser struct {
	ParserInterface
	Parser
}

// method to process a single day of a single channel
func (a *ARDParser) handleDay(channel Channel, day time.Time) {
	db := a.db
	flattenedProgramItems, err := a.fetchProgramItemsOfDay(channel, day)
	if err != nil {
		return
	}

	var programEntryList = &[]ProgramEntry{}
	// fill program entry list
	for _, item := range flattenedProgramItems {
		programEntry := ProgramEntry{}

		eid := strings.TrimSpace(item.Id)
		programEntry.Hash = buildHash([]string{
			eid,
			item.NumericId,
			fmt.Sprintf("%d", int(channel.ID)),
			fmt.Sprintf("%d", int(a.ChannelFamily.ID)),
			"program-entry",
		})

		programEntry.TechnicalID = eid

		// if there already is a program entry with this technical_id, use the original record
		entry := ProgramEntry{}
		db.Model(&entry).Where("hash = ?", programEntry.Hash).Where("channel_id = ?", channel.ID).Preload("ImageLinks").Find(&entry)
		if entry.ID != 0 {
			if entry.isRecentlyUpdated() {
				atomic.AddUint64(&status.TotalSkippedPE, 1)
				continue
			}
			programEntry = entry
		}
		// else create a new record

		var entryTitle = trimAndSanitizeString(item.CoreTitle)
		if item.CoreSubline != "" {
			entryTitle += " - " + trimAndSanitizeString(item.CoreSubline)
		}
		programEntry.Title = entryTitle
		programEntry.Description = trimAndSanitizeString(item.Synopsis)

		var startDate, endDate = item.BroadcastedOn, item.BroadcastEnd
		// atm it is not clear which information "BeginNet" contains - sometimes it seems to be the real broadcasting _end_ time.
		//if !item.BeginNet.IsZero() {
		//	startDate = item.BeginNet
		//}
		if startDate.IsZero() || endDate.IsZero() || startDate.After(endDate) {
			appLog(fmt.Sprintf("Invalid start date '%s' or end date '%s' for program entry with hash '%s'", startDate, endDate, programEntry.Hash))
			atomic.AddUint64(&status.TotalSkippedPE, 1)
			continue
		}

		programEntry.StartDateTime = &startDate
		programEntry.EndDateTime = &endDate
		programEntry.DurationMinutes = int16(programEntry.EndDateTime.Sub(*programEntry.StartDateTime).Minutes())

		programEntry.URL = trimAndSanitizeString(item.Links.Self.Href)
		if programEntry.URL != "" && !ardProgramPageMatcher.Match([]byte(programEntry.URL)) {
			appLog(fmt.Sprintf("Invalid url '%s' for program entry with hash '%s'. Skipping.", programEntry.URL, programEntry.Hash))
			atomic.AddUint64(&status.TotalSkippedPE, 1)
			continue
		}
		var details, err = getArdAPIResponse[ardAPIProgramItemDetail](programEntry.URL)
		if err != nil {
			appLog(fmt.Sprintf("error in call to ard API URL '%s': %v", programEntry.URL, err))
			return
		}

		programEntry.Homepage = trimAndSanitizeString(details.Video.WebUrl)
		if programEntry.Homepage != "" && !ardMediathekMatcher.Match([]byte(programEntry.Homepage)) {
			appLog(fmt.Sprintf("Invalid mediathek page '%s' for program entry with hash '%s'. Skipping.", programEntry.Homepage, programEntry.Hash))
			atomic.AddUint64(&status.TotalSkippedPE, 1)
			continue
		}

		// link channel and channel family
		programEntry.ChannelID = channel.ID
		programEntry.ChannelFamilyID = a.ChannelFamily.ID

		tag := trimAndSanitizeString(details.Grouping.Title)
		programEntry.considerTagExists(&tag)

		// image links
		imageURL := details.Video.ImageUrl
		if imageURL != "" && len(imageURL) > 5 {
			programEntry.considerImageLinkExists(imageURL)
		} else {
			x169Image := details.Images.Aspect16X9.Src
			x11Image := details.Images.Aspect1X1.Src
			x167Image := details.Images.Aspect16X7.Src
			if x169Image != "" && len(x169Image) > 5 {
				programEntry.considerImageLinkExists(strings.Replace(x169Image, "{width}", ardDefaultImageWidth, 1))
			} else if x11Image != "" && len(x11Image) > 5 {
				programEntry.considerImageLinkExists(strings.Replace(x11Image, "{width}", ardDefaultImageWidth, 1))
			} else if x167Image != "" && len(x167Image) > 5 {
				programEntry.considerImageLinkExists(strings.Replace(x167Image, "{width}", ardDefaultImageWidth, 1))
			}
		}

		if len(programEntry.ImageLinks) > 0 {
			for _, img := range programEntry.ImageLinks {
				if !ardImageLinkMatcher.MatchString(img.URL) && !ardImageLinkMatcher2.MatchString(img.URL) {
					appLog(fmt.Sprintf("Found invalid image link '%s' for program entry with hash '%s'. Skipping.", img.URL, programEntry.Hash))
					atomic.AddUint64(&status.TotalSkippedPE, 1)
					continue
				}
			}
		}

		programEntry.saveProgramEntryRecord(db)

		*programEntryList = append(*programEntryList, programEntry)
	}

	if verboseGlobal && len(*programEntryList) > 0 {
		log.Printf("ard channel program list has %d entries\n", len(*programEntryList))
	}
}

func (a *ARDParser) fetchProgramItemsOfDay(channel Channel, day time.Time) ([]ardAPIChannelProgramItem, error) {
	formattedDate := day.Format("2006-01-02")
	// the following line generated the URL we fetch the program entries of a single channel of a single day
	url := fmt.Sprintf("%s/program/api/program?day=%s", ardHostWithPrefix, formattedDate)

	response, err := getArdAPIResponse[ardDailyProgramOfChannelResponse](url)
	if err != nil {
		appLog(fmt.Sprintf("error in call to ard url '%s': %v", url, err))
		return nil, err
	}
	var flattenedProgramItems []ardAPIChannelProgramItem
	for _, singleChannel := range response.Channels {
		for _, slot := range singleChannel.TimeSlots {
			for _, item := range slot {
				// filter for items of the channel we are interested in
				if item.Channel.MainChannelId == channel.TechnicalID {
					flattenedProgramItems = append(flattenedProgramItems, item)
				}
			}
		}
	}
	if verboseGlobal {
		log.Printf("Received response from url '%s': %v", url, response)
	}
	return flattenedProgramItems, nil
}

func getArdAPIResponse[T any](url string) (*T, error) {
	headers := map[string]string{}
	resp, err := doGetRequest(url, headers, 3)
	if resp == nil || err != nil {
		errMsg := fmt.Sprintf("Problem fetching URL '%s' with error '%v'", url, err)
		appLog(errMsg)
		log.Println(errMsg)
		return nil, err
	}
	var response T
	jsonErr := json.Unmarshal([]byte(*resp), &response)
	if jsonErr != nil {
		errMsg := fmt.Sprintf("Invalid json format in ard api response. url: '%s'", url)
		appLog(errMsg)
		log.Println(errMsg)
		return nil, jsonErr
	}
	return &response, nil
}

func getArdAPIResponseForTvShows[T any](url string) (*T, error) {
	headers := map[string]string{
		"Host": ardMediathekAPIHost,
	}
	resp, err := doGetRequest(url, headers, 3)
	if resp == nil || err != nil {
		errMsg := fmt.Sprintf("Problem fetching URL '%s' with error '%v'", url, err)
		appLog(errMsg)
		log.Println(errMsg)
		return nil, err
	}
	var response T
	jsonErr := json.Unmarshal([]byte(*resp), &response)
	if jsonErr != nil {
		errMsg := fmt.Sprintf("Invalid json format in ard api response. url: '%s'", url)
		appLog(errMsg)
		log.Println(errMsg)
		return nil, jsonErr
	}
	return &response, nil
}

// method to fetch all tv show data
func (a *ARDParser) fetchTVShows() {
	if !appConf.EnableTVShowCollection || isRecentlyFetched() {
		a.logRecentFetch("Skip update of ard tv shows")
		return
	}

	// build set of urls to fetch tv shows from
	var tvShowURLs = make([]string, 0)
	for _, category := range ardTvShowCategories {
		if category == "" {
		}
		categoryString := strings.TrimSuffix(base64.StdEncoding.EncodeToString([]byte("Das Erste."+category)), "=")
		tvShowURLs = append(tvShowURLs, fmt.Sprintf("%s%s", ardMediathekAPITvShowPath, categoryString))
	}
	for _, apiURL := range tvShowURLs {
		// fetch sizes first
		response, err := getArdAPIResponseForTvShows[ardAPITvShowResponse](fmt.Sprintf("%s?pageSize=5", apiURL))
		if err != nil || response == nil {
			appLog(fmt.Sprintf("Problem fetching URL for tv show:'%s'", apiURL))
			continue
		}
		var pageSizeLimit = 200
		var totalElementsOfCategory = response.Pagination.TotalElements
		var pageCount = int(math.Ceil(float64(totalElementsOfCategory / pageSizeLimit)))

		for page := -1; page < pageCount; page++ {
			singlePageURL := fmt.Sprintf("%s?pageSize=%d&pageNumber=%d", apiURL, pageSizeLimit, page+1)
			response, err = getArdAPIResponseForTvShows[ardAPITvShowResponse](singlePageURL)
			if err != nil || response == nil {
				appLog(fmt.Sprintf("Problem fetching URL for tv show:'%s'", singlePageURL))
				continue
			}
			for _, teaser := range response.Teasers {
				var hash = buildHash([]string{
					fmt.Sprintf("%d", int(a.ChannelFamily.ID)),
					trimAndSanitizeString(teaser.LongTitle),
					trimAndSanitizeString(teaser.Id),
					"tv-show",
				})
				url := teaser.Links.Self.Href
				var tvShow = &TvShow{
					ManagedRecord: ManagedRecord{
						Title:           trimAndSanitizeString(teaser.LongTitle),
						URL:             url,
						Hash:            hash,
						Homepage:        url,
						TechnicalID:     teaser.Id,
						ChannelFamily:   a.ChannelFamily,
						ChannelFamilyID: a.ChannelFamily.ID,
					},
				}

				show := TvShow{}
				a.db.Model(&TvShow{}).Where("hash = ?", hash).Find(&show)
				if show.ID != 0 {
					tvShow.ID = show.ID
				}
				tvShow.saveTvShowRecord(a.db)
			}
		}
	}
}

// getTimeOfNextUpdate this function returns the next date time a fetch will take place considering the refresh interval of the configuration
func getTimeOfNextUpdate() time.Time {
	now := time.Now()
	if !isRecentlyFetched() || GetAppConf().ForceUpdate {
		return now
	}
	// it is recently fetched, return last update + refresh interval
	set := getSetting(settingKeyLastFetch)
	if set != nil && set.ID != 0 && len(set.Value) > 0 {
		lastUpdateTime, err := time.Parse(time.RFC3339, set.Value)
		if err != nil {
			log.Printf("Could not parse '%s' as date", set.Value)
			return now
		}
		location, _ := time.LoadLocation(GetAppConf().TimeZone)
		lastUpdateTime = lastUpdateTime.In(location)

		return lastUpdateTime.Add(time.Duration(GetAppConf().TimeToRefreshInMinutes) * time.Minute)
	}
	return now
}

// method which is called after the program entries and tv shows are fetched
func (a *ARDParser) postProcess() {}

// preProcess implementation
func (a *ARDParser) preProcess() bool {
	a.parallelWorkersCount = 10
	return true
}

func (a *ARDParser) isDateValidToFetch(day *time.Time) (bool, error) {
	if day == nil {
		return false, fmt.Errorf("invalid day")
	}
	if a.isMoreThanXDaysInFuture(day, 8) { // = six weeks in future + today
		return false, fmt.Errorf("maximum for days in future for ARD is 8")
	}
	if a.isMoreThanXDaysInPast(day, 8) {
		return false, fmt.Errorf("maximum for days in past for ARD is 8")
	}
	return true, nil
}

type ardDailyProgramOfChannelResponse struct {
	Links struct {
		Self struct {
			Type  string `json:"type"`
			Title string `json:"title"`
			Href  string `json:"href"`
		} `json:"self"`
	} `json:"links"`
	Channels      []ardAPIChannel `json:"channels"`
	TrackingPiano struct {
		PageTitle         string `json:"page_title"`
		PageInstitutionId string `json:"page_institution_id"`
		PageInstitution   string `json:"page_institution"`
		PageId            string `json:"page_id"`
		PageChapter2      string `json:"page_chapter2"`
		PageChapter1      string `json:"page_chapter1"`
	} `json:"trackingPiano"`
	TimeSlots []struct {
		Title       string    `json:"title"`
		HeightUnits int       `json:"heightUnits"`
		EndDate     string    `json:"endDate"`
		StartDate   time.Time `json:"startDate"`
	} `json:"timeSlots"`
	CreationDate time.Time `json:"creationDate"`
}

type ardAPIChannel struct {
	Id            string `json:"id"`
	TrackingPiano struct {
		WidgetType  string `json:"widget_type"`
		WidgetTitle string `json:"widget_title"`
		WidgetId    string `json:"widget_id"`
	} `json:"trackingPiano"`
	TimeSlots          [][]ardAPIChannelProgramItem `json:"timeSlots"`
	PublicationService struct {
		Name    string `json:"name"`
		Partner string `json:"partner"`
	} `json:"publicationService"`
	Crid             string `json:"crid"`
	LocalChannelList []struct {
		Id           string `json:"id"`
		Name         string `json:"name"`
		Crid         string `json:"crid"`
		LocalDefault bool   `json:"localDefault,omitempty"`
	} `json:"localChannelList,omitempty"`
}

type ardAPIChannelProgramItem struct {
	Id    string `json:"id"`
	Links struct {
		Self struct {
			Type  string `json:"type"`
			Title string `json:"title"`
			Href  string `json:"href"`
		} `json:"self"`
		Target struct {
			Type    string `json:"type"`
			Title   string `json:"title"`
			UrlId   string `json:"urlId"`
			Partner string `json:"partner"`
		} `json:"target,omitempty"`
	} `json:"links"`
	Type     string `json:"type"`
	Title    string `json:"title"`
	Duration int    `json:"duration"`
	Images   struct {
		Aspect16X9 struct {
			Title        string  `json:"title"`
			Text         *string `json:"text"`
			Alt          *string `json:"alt"`
			Src          string  `json:"src"`
			ProducerName string  `json:"producerName"`
		} `json:"aspect16x9"`
	} `json:"images,omitempty"`
	Channel struct {
		Id            string `json:"id"`
		Name          string `json:"name"`
		MainChannelId string `json:"main_channel_id"`
		LocalDefault  bool   `json:"localDefault,omitempty"`
	} `json:"channel"`
	TrackingPiano struct {
		WidgetSection       string `json:"widget_section"`
		TeaserTitle         string `json:"teaser_title"`
		TeaserRecommended   bool   `json:"teaser_recommended"`
		TeaserInstitutionId string `json:"teaser_institution_id"`
		TeaserInstitution   string `json:"teaser_institution"`
		TeaserId            string `json:"teaser_id"`
		TeaserContentType   string `json:"teaser_content_type"`
		TeaserRegion        string `json:"teaser_region,omitempty"`
	} `json:"trackingPiano"`
	CreationDate          time.Time `json:"creationDate"`
	HeightUnits           int       `json:"heightUnits"`
	BeginNet              time.Time `json:"beginNet"`
	BinaryFeatures        []string  `json:"binaryFeatures,omitempty"`
	MaturityContentRating string    `json:"maturityContentRating,omitempty"`
	BroadcastEnd          time.Time `json:"broadcastEnd"`
	BroadcastedOn         time.Time `json:"broadcastedOn"`
	CoreSubline           string    `json:"coreSubline"`
	CoreTitle             string    `json:"coreTitle"`
	LastMod               time.Time `json:"lastMod"`
	NumericId             string    `json:"numericId"`
	Subline               string    `json:"subline,omitempty"`
	Synopsis              string    `json:"synopsis,omitempty"`
	Live                  bool      `json:"live,omitempty"`
	BlockDuration         int       `json:"blockDuration,omitempty"`
	TrackingSplit         int       `json:"trackingSplit,omitempty"`
	Split                 bool      `json:"split,omitempty"`
	SplitPart             int       `json:"splitPart,omitempty"`
	IsLocal               bool      `json:"isLocal,omitempty"`
	FollowUps             []struct {
		Id    string `json:"id"`
		Links struct {
			Self struct {
				Type  string `json:"type"`
				Title string `json:"title"`
				Href  string `json:"href"`
			} `json:"self"`
			Target struct {
				Type    string `json:"type"`
				Title   string `json:"title"`
				UrlId   string `json:"urlId"`
				Partner string `json:"partner"`
			} `json:"target,omitempty"`
		} `json:"links"`
		Type     string `json:"type"`
		Title    string `json:"title"`
		Duration int    `json:"duration"`
		Images   struct {
			Aspect16X9 struct {
				Title        string  `json:"title"`
				Text         *string `json:"text"`
				Alt          string  `json:"alt"`
				Src          string  `json:"src"`
				ProducerName string  `json:"producerName"`
			} `json:"aspect16x9"`
		} `json:"images"`
		Channel struct {
			Id            string `json:"id"`
			Name          string `json:"name"`
			MainChannelId string `json:"main_channel_id"`
		} `json:"channel"`
		TrackingPiano struct {
			WidgetSection       string `json:"widget_section"`
			TeaserTitle         string `json:"teaser_title"`
			TeaserRecommended   bool   `json:"teaser_recommended"`
			TeaserInstitutionId string `json:"teaser_institution_id"`
			TeaserInstitution   string `json:"teaser_institution"`
			TeaserId            string `json:"teaser_id"`
			TeaserContentType   string `json:"teaser_content_type"`
		} `json:"trackingPiano"`
		CreationDate          time.Time `json:"creationDate"`
		HeightUnits           int       `json:"heightUnits"`
		BeginNet              time.Time `json:"beginNet"`
		BroadcastEnd          time.Time `json:"broadcastEnd"`
		BroadcastedOn         time.Time `json:"broadcastedOn"`
		CoreSubline           string    `json:"coreSubline"`
		CoreTitle             string    `json:"coreTitle"`
		LastMod               time.Time `json:"lastMod"`
		NumericId             string    `json:"numericId"`
		Subline               string    `json:"subline"`
		Synopsis              string    `json:"synopsis"`
		BinaryFeatures        []string  `json:"binaryFeatures,omitempty"`
		MaturityContentRating string    `json:"maturityContentRating,omitempty"`
	} `json:"followUps,omitempty"`
	GroupDuration int `json:"groupDuration,omitempty"`
}

type ardAPIProgramItemDetail struct {
	Id    string `json:"id"`
	Links struct {
		Self struct {
			Type  string `json:"type"`
			Title string `json:"title"`
			Href  string `json:"href"`
		} `json:"self"`
		Target struct {
			Type    string `json:"type"`
			Title   string `json:"title"`
			UrlId   string `json:"urlId"`
			Partner string `json:"partner"`
		} `json:"target"`
	} `json:"links"`
	Type     string `json:"type"`
	Title    string `json:"title"`
	Duration int    `json:"duration"`
	Channel  struct {
		Id            string `json:"id"`
		Name          string `json:"name"`
		Crid          string `json:"crid"`
		MainChannelId string `json:"main_channel_id"`
	} `json:"channel"`
	TrackingPiano struct {
		WidgetSection       string `json:"widget_section"`
		TeaserTitle         string `json:"teaser_title"`
		TeaserRecommended   bool   `json:"teaser_recommended"`
		TeaserInstitutionId string `json:"teaser_institution_id"`
		TeaserInstitution   string `json:"teaser_institution"`
		TeaserId            string `json:"teaser_id"`
		TeaserContentType   string `json:"teaser_content_type"`
	} `json:"trackingPiano"`
	CreationDate time.Time `json:"creationDate"`
	Synopsis     string    `json:"synopsis"`
	Subline      string    `json:"subline"`
	NumericId    string    `json:"numericId"`
	Images       struct {
		Aspect16X9 struct {
			Title        string `json:"title"`
			Text         string `json:"text"`
			Alt          string `json:"alt"`
			Src          string `json:"src"`
			ProducerName string `json:"producerName"`
		} `json:"aspect16x9"`
		Aspect16X7 struct {
			Title        string `json:"title"`
			Text         string `json:"text"`
			Alt          string `json:"alt"`
			Src          string `json:"src"`
			ProducerName string `json:"producerName"`
		} `json:"aspect16x7"`
		Aspect1X1 struct {
			Title        string `json:"title"`
			Text         string `json:"text"`
			Alt          string `json:"alt"`
			Src          string `json:"src"`
			ProducerName string `json:"producerName"`
		} `json:"aspect1x1"`
	} `json:"images"`
	HeightUnits int `json:"heightUnits"`
	Grouping    struct {
		Title string `json:"title"`
		Url   string `json:"url"`
	} `json:"grouping"`
	CoreTitle      string    `json:"coreTitle"`
	BroadcastedOn  time.Time `json:"broadcastedOn"`
	BroadcastEnd   time.Time `json:"broadcastEnd"`
	BinaryFeatures []string  `json:"binaryFeatures"`
	BeginNet       time.Time `json:"beginNet"`
	CoreSubline    string    `json:"coreSubline"`
	LastMod        time.Time `json:"lastMod"`
	Video          struct {
		AvailableFrom time.Time   `json:"availableFrom"`
		AvailableTo   time.Time   `json:"availableTo"`
		BroadcastedOn time.Time   `json:"broadcastedOn"`
		CoreAssetType string      `json:"coreAssetType"`
		CreatedAt     time.Time   `json:"createdAt"`
		CreatedBy     interface{} `json:"createdBy"`
		CreationBrid  interface{} `json:"creationBrid"`
		Duration      int         `json:"duration"`
		EpisodeNumber interface{} `json:"episodeNumber"`
		ExternalMedia []struct {
			MediaType string `json:"mediaType"`
			Ratio     string `json:"ratio"`
			Type      string `json:"type"`
			Url       string `json:"url"`
			Versions  []struct {
				Ratio string `json:"ratio"`
				Url   string `json:"url"`
			} `json:"versions"`
		} `json:"externalMedia"`
		Extras []struct {
			Index interface{} `json:"index"`
			Text  string      `json:"text"`
			Type  string      `json:"type"`
		} `json:"extras"`
		Fsk             string        `json:"fsk"`
		GroupingId      string        `json:"groupingId"`
		GroupingTitle   string        `json:"groupingTitle"`
		GroupingWebUrl  string        `json:"groupingWebUrl"`
		Id              string        `json:"id"`
		ImageCredit     string        `json:"imageCredit"`
		ImageUrl        string        `json:"imageUrl"`
		IsTrailer       bool          `json:"isTrailer"`
		SeasonNumber    interface{}   `json:"seasonNumber"`
		SingleReport    bool          `json:"singleReport"`
		Source          string        `json:"source"`
		SourceId        string        `json:"sourceId"`
		SourceUpdatedAt time.Time     `json:"sourceUpdatedAt"`
		TagIds          []interface{} `json:"tagIds"`
		Text            struct {
			Short string `json:"short"`
		} `json:"text"`
		Title     string      `json:"title"`
		UpdatedAt time.Time   `json:"updatedAt"`
		UpdatedBy interface{} `json:"updatedBy"`
		WebUrl    string      `json:"webUrl"`
	} `json:"video"`
}

type ardAPITvShowResponse struct {
	CompilationType string `json:"compilationType"`
	IsChildContent  bool   `json:"isChildContent"`
	Pagination      struct {
		PageNumber    int `json:"pageNumber"`
		PageSize      int `json:"pageSize"`
		TotalElements int `json:"totalElements"`
	} `json:"pagination"`
	Personalized bool `json:"personalized"`
	Teasers      []struct {
		AvailableSeasons []string `json:"availableSeasons"`
		BinaryFeatures   []string `json:"binaryFeatures"`
		CoreAssetType    string   `json:"coreAssetType"`
		Id               string   `json:"id"`
		Images           struct {
			Aspect16X9 struct {
				Alt          string `json:"alt"`
				ProducerName string `json:"producerName"`
				Src          string `json:"src"`
				Title        string `json:"title"`
			} `json:"aspect16x9"`
			Aspect16X7 struct {
				Alt          string `json:"alt"`
				ProducerName string `json:"producerName"`
				Src          string `json:"src"`
				Title        string `json:"title"`
			} `json:"aspect16x7"`
			Aspect1X1 struct {
				Alt          string `json:"alt"`
				ProducerName string `json:"producerName"`
				Src          string `json:"src"`
				Title        string `json:"title"`
			} `json:"aspect1x1"`
			Aspect3X4 struct {
				Alt          string `json:"alt"`
				ProducerName string `json:"producerName"`
				Src          string `json:"src"`
				Title        string `json:"title"`
			} `json:"aspect3x4,omitempty"`
		} `json:"images"`
		IsChildContent     bool   `json:"isChildContent"`
		IsFamilyFriendly   bool   `json:"isFamilyFriendly"`
		LongTitle          string `json:"longTitle"`
		MediumTitle        string `json:"mediumTitle"`
		Personalized       bool   `json:"personalized"`
		PublicationService struct {
			Name string `json:"name"`
			Logo struct {
				Title        string `json:"title"`
				Alt          string `json:"alt"`
				ProducerName string `json:"producerName"`
				Src          string `json:"src"`
				AspectRatio  string `json:"aspectRatio"`
			} `json:"logo"`
			PublisherType string `json:"publisherType"`
			Partner       string `json:"partner"`
			Id            string `json:"id"`
			CoreId        string `json:"coreId"`
		} `json:"publicationService"`
		Links struct {
			Self struct {
				Id      string `json:"id"`
				UrlId   string `json:"urlId"`
				Title   string `json:"title"`
				Href    string `json:"href"`
				Type    string `json:"type"`
				Partner string `json:"partner"`
			} `json:"self"`
			Target struct {
				Id      string `json:"id"`
				UrlId   string `json:"urlId"`
				Title   string `json:"title"`
				Href    string `json:"href"`
				Type    string `json:"type"`
				Partner string `json:"partner"`
			} `json:"target"`
		} `json:"links"`
		ShortTitle    string `json:"shortTitle"`
		TitleVisible  bool   `json:"titleVisible"`
		TrackingPiano struct {
			TeaserContentType   string `json:"teaser_content_type"`
			TeaserId            string `json:"teaser_id"`
			TeaserInstitution   string `json:"teaser_institution"`
			TeaserInstitutionId string `json:"teaser_institution_id"`
			TeaserTitle         string `json:"teaser_title"`
		} `json:"trackingPiano"`
		Type string `json:"type"`
	} `json:"teasers"`
	UserVisibility string `json:"userVisibility"`
	AZContent      bool   `json:"aZContent"`
	Id             string `json:"id"`
	Links          struct {
		Self struct {
			Id      string `json:"id"`
			UrlId   string `json:"urlId"`
			Title   string `json:"title"`
			Href    string `json:"href"`
			Type    string `json:"type"`
			Partner string `json:"partner"`
		} `json:"self"`
	} `json:"links"`
	Size          string `json:"size"`
	Swipeable     bool   `json:"swipeable"`
	Title         string `json:"title"`
	TitleVisible  bool   `json:"titleVisible"`
	TrackingPiano struct {
		TeaserRecommended bool   `json:"teaser_recommended"`
		WidgetId          string `json:"widget_id"`
		WidgetTitle       string `json:"widget_title"`
		WidgetTyp         string `json:"widget_typ"`
	} `json:"trackingPiano"`
	Type string `json:"type"`
}
