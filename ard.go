//
// oerc, alias oer-collector
// Copyright (C) 2021 emschu[aet]mailbox.org
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
	"bufio"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"log"
	"math"
	url2 "net/url"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	ardHost                   = "programm.ard.de"
	ardHostWithPrefix         = "https://" + ardHost
	ardHostWithPrefixInsecure = "http://" + ardHost
	ardTagHost                = "/TV/Programm/Load/Similar35?eid="
	ardMainTagPage            = "/TV/Themenschwerpunkte/"
)

var (
	ardEidMatcher           = regexp.MustCompile(`eid[0-9]+`)
	ardProgramURLMatcher    = regexp.MustCompile(`^/TV/Programm/Sender/.*`)
	ardTvShowLinkMatcher    = regexp.MustCompile(`^/TV/Sendungen-von-A-bis-Z/.*/.{0,16}`)
	ardIcalLinkMatcher      = regexp.MustCompile(`^/ICalendar/iCal---Sendung\?sendung=[0-9]+`)
	ardImageLinkAttrMatcher = regexp.MustCompile(`^((/sendungsbilder/original/[0-9]+/[a-zA-Z0-9_-]+\.(jpe?g|png|JPE?G|PNG))|((https?://programm.ard.de)?/files/.*\.(jpe?g|png|JPE?G|PNG)))`)
	ardMainTags             = map[string]string{
		"Film":            "Film/Alle-Filme/Alle-Filme",
		"Dokumentation":   "Dokus--Reportagen/Alle-Dokumentationen/Startseite",
		"Kultur":          "Musik-und-Kultur/Alle-Kultursendungen/Startseite",
		"Ratgeber":        "Ratgeber-der-ARD/Alle-Ratgeber/Alle-Ratgeber",
		"Magazin":         "Ratgeber-der-ARD/Magazine/Startseite",
		"Serie":           "Serien--Soaps/Serien-von-A-bis-Z/Startseite/Serien-von-A-bis-Z",
		"Unterhaltung":    "Unterhaltung/Alle-Unterhaltungssendungen/Startseite",
		"Show/Quiz":       "Unterhaltung/Show--Quiz/Startseite",
		"Kabarett/Comedy": "Unterhaltung/Kabarett--Comedy/Startseite",
		"Zoogeschichten":  "Unterhaltung/Zoogeschichten/Startseite",
	}
	ardSubTags = map[string]string{
		"Herzgefühl":            "Film/Herzgefuehl/Startseite",
		"Komödie":               "Film/Komoedie/Startseite",
		"Klassiker":             "Film/Klassiker/Startseite",
		"Heimatfilme":           "Film/Heimatfilme/Startseite",
		"Krimi":                 "Film/Krimi/Startseite",
		"Tatort":                "Film/Tatort/Startseite",
		"Polizeiruf 110":        "Film/Polizeiruf-110/Startseite",
		"Drama":                 "Film/Drama/Startseite",
		"Western":               "Film/Western/Startseite",
		"Kurzfilm":              "Film/Kurzfilm/Startseite",
		"Polit-Talkshow":        "Politik/Polit-Talkshows/Startseite",
		"Nachrichten":           "Politik/Nachrichten/Startseite",
		"Aktuelle-Reportagen":   "Politik/Aktuelle-Reportagen/Startseite",
		"Polit-Magazine":        "Politik/Politmagazine/Startseite",
		"Geschichte":            "Dokus--Reportagen/Geschichte/Startseite",
		"Kultur":                "Dokus--Reportagen/Kultur/Startseite",
		"Tiere":                 "Dokus--Reportagen/Tiere/Startseite",
		"Gesundheit":            "Dokus--Reportagen/Gesundheit/Startseite",
		"Umwelt/Natur":          "Dokus--Reportagen/Umwelt-und-Natur/Startseite",
		"Reisen":                "Dokus--Reportagen/Reisen/Startseite",
		"Eisenbahn":             "Dokus--Reportagen/Eisenbahn/Startseite",
		"Wissenschaft":          "Dokus--Reportagen/Wissenschaft/Startseite",
		"Wissensmagazin":        "Dokus--Reportagen/Wissensmagazine/Startseite",
		"Klassik/Oper/Tanz":     "Musik-und-Kultur/Klassik-Oper--Tanz/Startseite",
		"Popkultur":             "Musik-und-Kultur/Popkultur/Startseite",
		"Jazz":                  "Musik-und-Kultur/Jazz/Startseite",
		"Literatur":             "Musik-und-Kultur/Literatur/Startseite",
		"Architektur":           "Musik-und-Kultur/Architektur/Startseite",
		"Kultur-Dokumentation":  "Musik-und-Kultur/Kultur-Dokumentationen/Startseite",
		"Kulturmagazine":        "Musik-und-Kultur/Kulturmagazine/Startseite",
		"Heim-/Gartenratgeber":  "Ratgeber-der-ARD/Heim-und-Garten/Startseite",
		"Reiseratgeber":         "Ratgeber-der-ARD/Reisen/Startseite",
		"Gesundheitsratgeber":   "Ratgeber-der-ARD/Gesundheit/Startseite",
		"Natur-/Umweltratgeber": "Ratgeber-der-ARD/Natur-und-Umwelt/Startseite",
		"Magazin":               "Ratgeber-der-ARD/Magazine/Startseite",
		"Coronavirus":           "Ratgeber-der-ARD/Coronavirus/Startseite",
		"Kochen":                "Kochen/Alle-Sendungen/Startseite",
		"Fußball":               "Sport/Fussball-im-TV/Startseite",
		"Sport":                 "Sport/Alle-Sportsendungen/Startseite",
		"Sportmagazin":          "Sport/Sportmagazine/Startseite",
		"Soap/Telenovela":       "Serien--Soaps/Soaps-und-Telenovelas/Startseite/Startseite",
		"Dokusoap":              "Serien--Soaps/Dokusoaps/Startseite",
		"Show/Quiz":             "Unterhaltung/Show--Quiz/Startseite",
		"Kabarett/Comedy":       "Unterhaltung/Kabarett--Comedy/Startseite",
		"Schlager/Volksmusik":   "Unterhaltung/Schlager--Volksmusik/Startseite",
		"Talkshow":              "Unterhaltung/Talkshows/Startseite",
		"Zoogeschichten":        "Unterhaltung/Zoogeschichten/Startseite",
		"Fernsehgottesdienst":   "Kirche-und-Religion/Fernsehgottesdienste/Startseite",
		"Religion":              "Kirche-und-Religion/Religion-Fernsehen/Startseite",
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
	// Create a Collector specifically for ard
	c := a.newArdCollector()

	var programEntryList = &[]ProgramEntry{}
	c.OnHTML(".event-list li[class^=eid]", func(element *colly.HTMLElement) {
		programEntry := ProgramEntry{}

		// get eid
		eid := ardEidMatcher.FindString(element.Attr("class"))
		programEntry.Hash = buildHash([]string{
			eid,
			fmt.Sprintf("%d", int(a.ChannelFamily.ID)),
			"program-entry",
		})

		// this is safe because the query selector for this closure handles this
		eid = strings.Replace(eid, "eid", "", 1)
		if len(eid) > 256 {
			appLog(fmt.Sprint("Invalid eid detected. length > 256"))
			return
		}
		programEntry.TechnicalID = eid

		// if there already is a program entry with this technical_id, use the original record
		entry := ProgramEntry{}
		db.Model(&entry).Where("hash = ?", programEntry.Hash).Where("channel_id = ?", channel.ID).Preload("ImageLinks").Find(&entry)
		if entry.ID != 0 {
			if entry.isRecentlyUpdated() {
				atomic.AddUint64(&status.TotalSkippedPE, 1)
				return
			}
			programEntry = entry
		}

		title := trimAndSanitizeString(element.DOM.Find("span.title").Text())
		subtitle := trimAndSanitizeString(element.DOM.Find("span.subtitle").Text())

		// subtitle (nested span) is removed from title and added to description
		title = trimAndSanitizeString(strings.Replace(title, subtitle, "", 1))
		programEntry.Title = title
		// reset description field
		programEntry.Description = ""
		if len(subtitle) > 0 {
			programEntry.Description = subtitle + ". "
		}

		find := element.DOM.Find("a")
		urlOfEntry, attrExists := find.Attr("href")
		if !attrExists {
			appLog(fmt.Sprintf("No 'href' attribute on program detail page. EID: %s", eid))
			return
		}
		if !ardProgramURLMatcher.Match([]byte(urlOfEntry)) {
			appLog(fmt.Sprintf("Invalid url '%s' on program detail page detected.", urlOfEntry))
			return
		}

		programEntry.URL = ardHostWithPrefix + urlOfEntry
		// link channel and channel family
		programEntry.ChannelID = channel.ID
		programEntry.ChannelFamilyID = a.ChannelFamily.ID

		*programEntryList = append(*programEntryList, programEntry)
	})

	// fire fetching all program entries from all channels for the defined time range
	formattedDate := day.Format("02.01.2006")
	// the following line generated the URL we fetch the program entries of
	url := fmt.Sprintf("%s/TV/Programm/Sender?datum=%s&hour=0&sender=%s", ardHostWithPrefix, formattedDate, channel.Hash)
	err := c.Visit(url)
	if err != nil {
		appLog(fmt.Sprintf("error in call to url '%s': %v", url, err))
	}
	c.Wait()

	if verboseGlobal && len(*programEntryList) > 0 {
		log.Printf("program list has %d entries\n", len(*programEntryList))
	}

	if len(*programEntryList) == 0 {
		return
	}

	c2 := c.Clone()
	var programEntry *ProgramEntry

	// this is called for each single program detail page
	c2.OnHTML("body .program-con", func(element *colly.HTMLElement) {
		if programEntry == nil {
			appLog(fmt.Sprintf("No program entry pointer found. This should never happen."))
			return
		}
		// fetch ical-links for proper datetime information
		icalHref, exists := element.DOM.Find("a[href*=ICalendar]").Attr("href")
		if !exists {
			appLog(fmt.Sprintf("ERROR: No iCal link found for program entry '%s'", programEntry.Hash))
			return
		}
		if !ardIcalLinkMatcher.Match([]byte(icalHref)) {
			appLog(fmt.Sprintf("Invalid iCal link found for program entry hash '%s'", programEntry.Hash))
			return
		}
		icalLink := ardHostWithPrefix + icalHref

		icalContent, err := a.parseStartAndEndDateTimeFromIcal(icalLink)
		if icalContent == nil || err != nil {
			appLog(fmt.Sprintf("Problem fetching ical at link '%s'", icalLink))
			return
		}

		programEntry.StartDateTime = &icalContent.startDate
		programEntry.EndDateTime = &icalContent.endDate
		programEntry.DurationMinutes = int16(icalContent.endDate.Sub(icalContent.startDate).Minutes())

		// add "tags"
		tags, tagErr := tryToFindTags(programEntry.TechnicalID)
		if tagErr != nil {
			tagErrMsg := fmt.Sprintf("Problem fetching tags of eid '%s'. %v.", programEntry.TechnicalID, tagErr)
			appLog(tagErrMsg)
			return
		}
		programEntry.Tags = strings.Join(*tags, ";")

		descrSelector := fmt.Sprintf("#mehr-%s .eventText", programEntry.TechnicalID)
		desc := element.DOM.Find(descrSelector)
		text := desc.Text()
		if !strings.HasPrefix(programEntry.Description, "Keine weiteren Informationen") {
			programEntry.Description += trimAndSanitizeString(text)
		}

		if len(programEntry.Description) == 0 {
			// try an alternative description location
			descr2Selector := fmt.Sprintf("div.detail-top div.eventText")
			desc = element.DOM.Find(descr2Selector)
			programEntry.Description = trimAndSanitizeString(desc.Text())
			if len(programEntry.Description) == 0 {
				programEntry.Description = "Keine weiteren Informationen"
			}
		}

		// add image links
		element.DOM.Find(".media-con img").Each(func(i int, selection *goquery.Selection) {
			attr, srcAttrExists := selection.Attr("src")
			if !ardImageLinkAttrMatcher.Match([]byte(attr)) {
				appLog(fmt.Sprintf("Invalid image link detected! program entry hash: '%s'", programEntry.Hash))
				return
			}
			for _, existingEntry := range programEntry.ImageLinks {
				if existingEntry.URL == (ardHostWithPrefix+attr) || existingEntry.URL == (ardHostWithPrefixInsecure+attr) {
					// already exists
					return
				}
			}
			if srcAttrExists {
				if !strings.HasPrefix(attr, ardHostWithPrefix) && !strings.HasPrefix(attr, ardHostWithPrefixInsecure) {
					attr = ardHostWithPrefix + attr
				}
				programEntry.ImageLinks = append(programEntry.ImageLinks, ImageLink{URL: attr})
			}
		})

		desc = element.DOM.Find(".bcData a").Each(func(i int, s *goquery.Selection) {
			text := trimAndSanitizeString(s.Text())
			if strings.Contains(text, "Sendungsseite im Internet") {
				attr, hrefAttrExists := s.Attr("href")
				if hrefAttrExists {
					u, err := url2.ParseRequestURI(attr)
					if err != nil {
						appLog("Invalid url of program entry's homepage found!")
						return
					}
					programEntry.Homepage = fmt.Sprintf("%s%s", u.Host, u.RequestURI())
				}
			}
		})

		programEntry.saveProgramEntryRecord(db)
	})

	for _, pe := range *programEntryList {
		programEntry = &pe
		err := c2.Visit(pe.URL)
		if err != nil {
			appLog(fmt.Sprintf("Error in ard tv show call to url '%s':%v", ardHostWithPrefix+pe.URL, err))
		}
		c2.Wait()
	}

	a.linkTagsToEntriesDaily(day)
}

// method to fetch all tv show data
func (a *ARDParser) fetchTVShows() {
	if !GetAppConf().EnableTVShowCollection || isRecentlyFetched() {
		log.Printf("Skip update of tv shows, due to recent fetch. Use 'forceUpdate' = true to ignore this.")
		return
	}
	// Create a Collector specifically for Shopify
	collector := a.newArdCollector()

	// Create a callback on the XPath query searching for the URLs
	collector.OnHTML(".az-slick > .box > a", func(e *colly.HTMLElement) {
		var link = e.Attr("href")
		if !ardTvShowLinkMatcher.Match([]byte(link)) {
			appLog(fmt.Sprintf("Invalid link '%s' for ard tv show detected. Skipping entry.", link))
			return
		}

		var title = trimAndSanitizeString(e.ChildAttr("img", "title"))
		if link == "" || title == "" {
			appLog(fmt.Sprintf("ERR: empty link or title in URL '%s'", e.Request.URL.EscapedPath()))
			return
		}
		var hash = buildHash([]string{
			fmt.Sprintf("%d", int(a.ChannelFamily.ID)),
			title,
			"tv-show",
		})
		url := ardHost + link
		var tvShow = &TvShow{
			ManagedRecord: ManagedRecord{
				Title:           title,
				URL:             url,
				Hash:            hash,
				Homepage:        url,
				ChannelFamily:   a.ChannelFamily,
				ChannelFamilyID: a.ChannelFamily.ID, // 0 = ard
			},
		}

		show := TvShow{}
		a.db.Model(&TvShow{}).Where("hash = ?", hash).Find(&show)
		if show.ID != 0 {
			tvShow.ID = show.ID
		}
		tvShow.saveTvShowRecord(a.db)
	})

	// Start the collector
	tvShowURL := ardHostWithPrefix + "/TV/Sendungen-von-A-bis-Z/Startseite?page=&char=all"
	err := collector.Visit(tvShowURL)
	if err != nil {
		appLog(fmt.Sprintf("Problem scraping URL '%s'", tvShowURL))
	}
	collector.Wait()
	// TODO add tv show post processing: image links + tags + related program entries
}

// helper method to get a collector instance
func (a *ARDParser) newArdCollector() *colly.Collector {
	return baseCollector([]string{ardHost})
}

// method to identify the "tags" a program entry has
func tryToFindTags(eid string) (*[]string, error) {
	var tags = &[]string{}

	var url = fmt.Sprintf("%s%s%s", ardHostWithPrefix, ardTagHost, eid)
	doc, err := getDocument(url)
	if doc == nil || err != nil {
		return tags, err
	}
	doesTagExist := func(name string) bool {
		for _, existingTag := range *tags {
			if existingTag == name {
				return true
			}
		}
		return false
	}
	doc.Find("form[id^=bookmark-checks] .row span[class*=similar-events-bookmark]").Each(func(i int, selection *goquery.Selection) {
		tagText := trimAndSanitizeString(selection.Text())
		if !doesTagExist(tagText) {
			*tags = append(*tags, tagText)
		}
	})
	return tags, nil
}

// method which is called after the program entries and tv shows are fetched
func (a *ARDParser) postProcess() {
	a.linkTagsToEntriesGeneral()
}

// method which is called after the program entries and tv shows are fetched: empty for ARD
func (a *ARDParser) preProcess() bool {
	a.parallelWorkersCount = runtime.NumCPU()
	return true
}

func (a *ARDParser) isDateValidToFetch(day *time.Time) (bool, error) {
	if day == nil {
		return false, fmt.Errorf("invalid day")
	}
	if a.isMoreThanXDaysInFuture(day, 43) { // = six weeks in future + today
		return false, fmt.Errorf("maximum for days in future for ARD is 43")
	}
	location, _ := time.LoadLocation(GetAppConf().TimeZone)
	earliestDate := time.Date(2010, 1, 1, 0, 0, 0, 0, location)
	if day.Before(earliestDate) {
		return false, fmt.Errorf("maximum for days in past for ARD is %s", earliestDate.Format(time.RFC822))
	}
	return true, nil
}

// method to link tags to program entries of a single day
func (a *ARDParser) linkTagsToEntriesDaily(day time.Time) {
	if isRecentlyFetched() && !GetAppConf().ForceUpdate {
		return
	}

	// handle main tags
	for mainTagName, tagURLPart := range ardMainTags {
		formattedDate := day.Format("02.01.2006")
		dailyURL := fmt.Sprintf(
			"%s%s%s?datum=%s&hour=0&ajaxPageLoad=1",
			ardHostWithPrefix,
			ardMainTagPage,
			tagURLPart,
			formattedDate,
		)
		eidList := a.getEIDsOfUrls([]string{dailyURL})

		var programEntry ProgramEntry
		if len(eidList) > 0 {
			if len(eidList) == 1 {
				a.db.Model(ProgramEntry{}).Where("technical_id LIKE ?", eidList[0]).Find(&programEntry)
			} else {
				a.db.Model(ProgramEntry{}).Where("technical_id IN(?)", eidList).Find(&programEntry)
			}
			programEntry.considerTagExists(&mainTagName)
		}
	}
}

// method to link tags to program entries
func (a *ARDParser) linkTagsToEntriesGeneral() {
	if isRecentlyFetched() {
		log.Printf("Skip update of ard program entry tag search, due to recent fetch. Use 'forceUpdate' = true to ignore this.")
		return
	}

	// handle sub-tags
	for subTagName, tagURLPart := range ardSubTags {
		previewURL := fmt.Sprintf("%s%s%s?ajaxPageLoad=1", ardHostWithPrefix, ardMainTagPage, tagURLPart)
		archiveURL := fmt.Sprintf("%s&archiv=1", previewURL)
		eidList := a.getEIDsOfUrls([]string{previewURL, archiveURL})

		var programEntry ProgramEntry
		if len(eidList) > 0 {
			if len(eidList) == 1 {
				a.db.Model(&ProgramEntry{}).Where("technical_id LIKE ?", eidList[0]).Find(&programEntry)
			} else {
				// we have to ensure that the IN-operation of the SQL database is has a limited input length
				// typically an eid has 15 chars, allow 15 x 15 chars = 225 chars in IN-query + 14 times ","
				var blockSize = 14

				for len(eidList) > 0 {
					highestIndex := int(math.Min(float64(blockSize), float64(len(eidList)-1)))
					list := eidList[:highestIndex]
					a.db.Model(ProgramEntry{}).Where("technical_id IN (?)", list).Find(&programEntry)
					if (len(eidList)-1) >= highestIndex && highestIndex > 0 {
						eidList = eidList[highestIndex:]
					} else {
						break
					}
				}
			}
			programEntry.considerTagExists(&subTagName)
		}
	}
}

// getEIDsOfUrls get eid of urls, these urls should be checked to be not malicious
func (a *ARDParser) getEIDsOfUrls(urls []string) []string {
	c := a.newArdCollector()
	var eidList []string
	var listMutex sync.Mutex

	c.OnHTML(".event-list li[class^=eid]", func(element *colly.HTMLElement) {
		eid := strings.Replace(ardEidMatcher.FindString(element.Attr("class")), "eid", "", 1)
		listMutex.Lock()
		eidList = append(eidList, eid)
		listMutex.Unlock()
	})

	for _, url := range urls {
		urlErr := c.Visit(url)
		if urlErr != nil {
			errMsg := fmt.Sprintf("Problem fetching URL '%s'. %v.", url, urlErr)
			appLog(errMsg)
		}
	}
	c.Wait()

	return eidList
}

// ICalContent object to wrap retrieved ical content
type ICalContent struct {
	startDate time.Time
	endDate   time.Time
}

// parseStartAndEndDateTimeFromIcal method to parse a plain ical file data just for DTSTART and DTEND. Needed for ARD only atm.
func (a *ARDParser) parseStartAndEndDateTimeFromIcal(targetURL string) (*ICalContent, error) {
	requestHeaders := map[string]string{"Accept": "text/html", "Host": "programm.ard.de"}

	icalContent, err := doGetRequest(targetURL, requestHeaders, 3)
	if icalContent == nil || err != nil {
		return nil, err
	}

	location, _ := time.LoadLocation(GetAppConf().TimeZone)

	scanner := bufio.NewScanner(strings.NewReader(*icalContent))
	scanner.Split(bufio.ScanLines)

	var hasStart = false
	var hasEnd = false
	var content = ICalContent{}
	const iCalDateLayout = "20060102T150405"
	_, zoneOffsetInSecs := time.Now().In(location).Zone()
	timeZoneOffset := time.Duration(-zoneOffsetInSecs) * time.Second

	for scanner.Scan() {
		line := scanner.Text()
		if !hasStart && strings.HasPrefix(line, "DTSTART;") {
			startDate := strings.Replace(line, "DTSTART;TZID=Europe/Berlin:", "", 1)
			content.startDate, err = time.Parse(iCalDateLayout, startDate)
			if err != nil {
				appLog(fmt.Sprintf("Problem with date DTSTART in ical data of '%v': %v.", icalContent, err))
			} else {
				hasStart = true
			}
		}
		if !hasEnd && strings.HasPrefix(line, "DTEND;") {
			endDate := strings.Replace(line, "DTEND;TZID=Europe/Berlin:", "", 1)
			content.endDate, err = time.Parse(iCalDateLayout, endDate)
			if err != nil {
				appLog(fmt.Sprintf("Problem with date DTEND in ical data of '%v': %v.", icalContent, err))
			} else {
				hasEnd = true
			}
		}
		if hasStart && hasEnd {
			// its important to subtract this offset and set the correct time zone here
			content.startDate = content.startDate.Add(timeZoneOffset).In(location)
			content.endDate = content.endDate.Add(timeZoneOffset).In(location)
			break
		}
	}
	if !hasStart || !hasEnd {
		return nil, errors.New("Could not find start and/or end date in supplied ical content")
	}
	if content.startDate.IsZero() || content.endDate.IsZero() {
		return nil, errors.New("Empty dates detected in ical content. Probably a parser error")
	}
	return &content, nil
}
