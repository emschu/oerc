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
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/alitto/pond"
	"github.com/gocolly/colly/v2"
	"gorm.io/gorm"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

const (
	srfHost           = "www.srf.ch"
	srfHostWithPrefix = "https://" + srfHost
)

var (
	srfTvShowURLMatcher = regexp.MustCompile(`^(/play/tv/sendung/.*|(https?://www\.srgssr\.ch/?.*))$`)
)

// ParseSRF central method to parse SRF tv show and program data
func ParseSRF() {
	db := getDb()

	// get channel family db record
	var channelFamily = getChannelFamily(db, "SRF")
	if channelFamily.ID == 0 {
		log.Fatalln("SRF channelFamily was not found!")
		return
	}

	// import tv shows
	if GetAppConf().EnableTVShowCollection {
		fetchTvShowsSRF(db, channelFamily)
	}

	// import program entries for the configured date range
	if GetAppConf().EnableProgramEntryCollection {
		daysInPast := GetAppConf().DaysInPast
		daysInFuture := GetAppConf().DaysInFuture

		if daysInPast > 15 {
			warnMsg := "Maximum for days in past for SRF is 15!\n"
			log.Printf(warnMsg)
			appLog(warnMsg)
			daysInPast = 15
		}
		times := *generateDateRangeInPastAndFuture(daysInPast, daysInFuture)
		pool := pond.New(4, 100, getWorkerPoolIdleTimeout())
		for _, channel := range getChannelsOfFamily(db, channelFamily) {
			for _, day := range times {
				if int(time.Since(day).Hours()/24) <= 30 {
					family := *channelFamily
					chn := channel
					dayToFetch := day

					// srf specific limits!
					pool.Submit(func() {
						handleDaySRF(db, family, chn, dayToFetch)
					})
				}
			}
		}
		// wait for finish
		pool.StopAndWait()
	}

	if verboseGlobal {
		log.Println("SRF parsed successfully")
	}
}

// fetchTvShowSRF: This method checks all the tv shows
func fetchTvShowsSRF(db *gorm.DB, family *ChannelFamily) {
	if !GetAppConf().EnableTVShowCollection || isRecentlyFetched() {
		log.Printf("Skip update of tv shows, due to recent fetch. Use 'forceUpdate' = true to ignore this.")
		return
	}

	c := srfCollector()
	c.OnHTML("section", func(element *colly.HTMLElement) {
		title := trimAndSanitizeString(element.DOM.Find("h2").Text())
		url, _ := element.DOM.Find("a").Attr("href")

		if !srfTvShowURLMatcher.Match([]byte(url)) {
			appLog(fmt.Sprintf("Unexpected url detected: '%s'", url))
			return
		}

		var hash = buildHash([]string{
			fmt.Sprintf("%d", int(family.ID)),
			title,
			"tv-show",
		})
		var tvShow = &TvShow{
			ManagedRecord: ManagedRecord{
				Title:           title,
				URL:             srfHostWithPrefix + url,
				Hash:            hash,
				Homepage:        srfHostWithPrefix + url,
				ChannelFamily:   *family,
				ChannelFamilyID: family.ID, // 3 = srf
			},
		}

		show := TvShow{}
		db.Where("hash = ?", hash).Find(&show)
		if show.ID != 0 {
			tvShow.ID = show.ID
		}
		saveTvShowRecord(db, tvShow)
	})

	err := c.Visit("https://www.srf.ch/play/tv/sendungen")
	if err != nil {
		appLog(fmt.Sprintf("Problem fetching URL 'https://www.srf.ch/play/tv/sendungen' %v.\n", err))
	}
	c.Wait()
}

// handleDaySRF: handle a single day to fetch SRF program data
func handleDaySRF(db *gorm.DB, family ChannelFamily, channel Channel, day time.Time) {
	c := srfCollector()
	location, _ := time.LoadLocation(GetAppConf().TimeZone)

	queryURL := fmt.Sprintf("https://www.srf.ch/programm/tv/sender/%s/%02d-%02d-%d", channel.TechnicalID, day.Day(), day.Month(), day.Year())

	c.OnHTML("#content .channel-show", func(element *colly.HTMLElement) {
		startDateElement := trimAndSanitizeString(element.DOM.Find(".channel-show__airtime .channel-show__begin").Text())
		endDateElement := trimAndSanitizeString(element.DOM.Find(".channel-show__airtime .channel-show__stop").Text())
		linkElement := trimAndSanitizeString(element.DOM.Find(".channel-show__link").AttrOr("href", ""))

		var programEntry ProgramEntry
		if len(startDateElement) > 0 {
			programEntry.StartDateTime = getDateInterpretation(false, day, startDateElement)
		}
		if len(endDateElement) > 0 {
			programEntry.EndDateTime = getDateInterpretation(false, day, endDateElement)
		}
		title := trimAndSanitizeString(element.DOM.Find("h3.channel-show__title").Text())
		// in pre-processing we cannot decide on which day the given entry is televised. #
		// we do this in the next step
		if programEntry.StartDateTime == nil || programEntry.EndDateTime == nil {
			appLog(fmt.Sprintf("Problem parsing start date time (%s) or end date time (%s). Skipping program entry of page '%s' with title '%s'.",
				programEntry.StartDateTime, programEntry.EndDateTime, queryURL, title))
			return
		}
		subTitle := trimAndSanitizeString(element.DOM.Find(".channel-show__subtitle").Text())
		programEntry.Title = title

		if len(linkElement) > 0 {
			programEntry.URL = strings.Replace(linkElement, "//", "https://", 1)
		}

		hash := buildHash([]string{
			title,
			programEntry.StartDateTime.String(),
			linkElement,
			strconv.Itoa(int(family.ID)),
			strconv.Itoa(int(channel.ID)),
			"program-entry",
		})
		programEntry.Hash = hash
		programEntry.TechnicalID = hash
		programEntry.ChannelFamily = family
		programEntry.ChannelFamilyID = family.ID
		programEntry.Channel = channel
		programEntry.ChannelID = channel.ID

		entry := ProgramEntry{}
		db.Model(&entry).Where("hash = ?", programEntry.Hash).Where("channel_id = ?", channel.ID).Find(&entry)
		if entry.ID != 0 {
			if isRecentlyUpdated(&entry) {
				atomic.AddUint64(&status.TotalSkippedPE, 1)
				return
			}
			programEntry.ID = entry.ID
		}

		if len(programEntry.URL) == 0 {
			// leave here
			appLog(fmt.Sprintf("Empty URL for program entry."))
			return
		}

		requestHeaders := map[string]string{
			"Origin": "srf.ch",
			"Accept": "text/html",
		}
		response, err := doGetRequest(programEntry.URL, requestHeaders, 3)
		if response == nil || err != nil {
			log.Printf("Problem fetching URL '%s' %v.\n", programEntry.URL, err)
			return
		}
		reader := strings.NewReader(*response)
		doc, err := goquery.NewDocumentFromReader(reader)
		if err != nil {
			appLog(fmt.Sprintf("Could not create reader with url '%s'. Error was: %v", programEntry.URL, err))
			return
		}
		// start and end date parsing
		infoBox := doc.Find(".detail--content .infobox .left p")
		if len(infoBox.Nodes) > 2 {
			first := trimAndSanitizeString(infoBox.First().Text())
			sec := trimAndSanitizeString(infoBox.Next().First().Text())
			dateStr := first + " " + sec
			if len(dateStr) > 50 {
				// there is something wrong...
				appLog(fmt.Sprintf("Invalid date range for srf program entry detected. Hash: %s", programEntry.Hash))
				return
			}
			dateOfEntry := getDateFromStringSrf(dateStr)
			if dateOfEntry == nil {
				appLog("Could not parse date in srf program entry: #")
				return
			}
			assumedStartDateTime := programEntry.StartDateTime
			if assumedStartDateTime.Day() != dateOfEntry.Day() {
				newStartTime := time.Date(dateOfEntry.Year(), dateOfEntry.Month(), dateOfEntry.Day(),
					programEntry.StartDateTime.Hour(), programEntry.StartDateTime.Minute(), 0, 0, location)
				programEntry.StartDateTime = &newStartTime
				durationInMinutes := int16(math.Abs(assumedStartDateTime.Sub(*programEntry.EndDateTime).Minutes()))
				newEndTime := newStartTime.Add(time.Duration(durationInMinutes) * time.Minute)
				programEntry.EndDateTime = &newEndTime
			}
		}
		if programEntry.EndDateTime.Before(*programEntry.StartDateTime) {
			newEndDateTime := programEntry.EndDateTime.Add(24 * time.Hour)
			programEntry.EndDateTime = &newEndDateTime
		}

		if programEntry.EndDateTime == nil || programEntry.EndDateTime.Before(*programEntry.StartDateTime) {
			appLog(fmt.Sprintf("Problem with end or start date in program entry with hash '%s'.", programEntry.Hash))
			return
		}

		durationInMinutes := int16(math.Abs(programEntry.StartDateTime.Sub(*programEntry.EndDateTime).Minutes()))
		programEntry.DurationMinutes = durationInMinutes

		// handle description
		if len(subTitle) > 0 {
			programEntry.Description = subTitle + "\n"
		}
		var desc string
		leadElements := trimAndSanitizeString(doc.Find("p.lead").Text())
		descElements := trimAndSanitizeString(doc.Find("p.description").Text())

		var actorsList string
		doc.Find("ul.actors li").Each(func(i int, selection *goquery.Selection) {
			text := trimAndSanitizeString(selection.Text())
			if text != "" {
				actorsList += text
				// try to improve formatting a little
				if text != "Mit" && text != "Gastgeber" && text != "Moderation" && text != "Kommentar" {
					actorsList += "<br/>"
				} else {
					actorsList += " "
				}
			}
		})
		actorElements := strings.TrimSuffix(trimAndSanitizeString(actorsList), "<br/>")
		if len(leadElements) > 0 {
			desc += leadElements + "<br/>"
		}
		if len(descElements) > 0 {
			desc += descElements + "<br/>"
		}
		if len(actorElements) > 0 {
			desc += actorElements + "<br/>"
		}
		programEntry.Description += desc

		// handle tags
		if len(infoBox.Nodes) > 4 {
			var genre string
			if infoBox.Next() != nil && infoBox.Next().Next() != nil && infoBox.Next().Next().Next() != nil && infoBox.Next().Next().Next().First() != nil {
				selection := infoBox.Next().Next().First().Text()
				if strings.Contains(selection, "Wiederholung") {
					genre = infoBox.Next().Next().Next().First().Text()
				} else {
					genre = selection
				}

				if len(genre) > 0 {
					considerTagExists(&programEntry, &genre)
				}
			}
		}

		var homepage string
		doc.Find(".detail--content a").Each(func(i int, selection *goquery.Selection) {
			attr, exists := selection.Attr("title")
			attr = trimAndSanitizeString(attr)
			if exists && strings.Contains(attr, "Sendungsseite") {
				homepage = trimAndSanitizeString(selection.AttrOr("href", ""))
			}
		})
		programEntry.Homepage = strings.Replace(homepage, "//", "https://", 1)

		saveProgramEntryRecord(db, &programEntry)
	})

	err := c.Visit(queryURL)
	if err != nil {
		appLog(fmt.Sprintf("Error of collector: %v\n", err))
	}
	c.Wait()
}

// getDateFromStringSrf: Extract a date object of a given string format of the page
func getDateFromStringSrf(dateStr string) *time.Time {
	if len(dateStr) == 0 {
		return nil
	}
	firstPointIndex := strings.Index(dateStr, ".")
	if firstPointIndex < 0 {
		return nil
	}
	// two points '.' in the string are needed for this function to work
	secondPointIndex := firstPointIndex + strings.Index(dateStr[firstPointIndex+1:], ".")
	if secondPointIndex < 0 {
		return nil
	}
	if len(dateStr) < secondPointIndex+5 {
		// not plausible
		return nil
	}
	var day, month, year int64
	var err error

	day, err = strconv.ParseInt(dateStr[firstPointIndex-2:firstPointIndex], 10, 64)
	if err != nil {
		log.Printf("Problem parsing day: %v\n", err)
		return nil
	}
	month, err = strconv.ParseInt(dateStr[firstPointIndex+1:secondPointIndex+1], 10, 64)
	if err != nil {
		log.Printf("Problem parsing month: %v\n", err)
		return nil
	}
	s := dateStr[secondPointIndex+2:]
	year, err = strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Printf("Problem parsing year: %v\n", err)
		return nil
	}
	location, _ := time.LoadLocation(GetAppConf().TimeZone)
	date := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, location)
	return &date
}

// helper method to get a collector instance
func srfCollector() *colly.Collector {
	return baseCollector([]string{srfHost})
}

// getDateInterpretation: method to parse a date string (used by ORF) into a single time object
func getDateInterpretation(isNight bool, day time.Time, timeString string) *time.Time {
	middle := strings.Index(timeString, ":")
	if middle == -1 {
		return nil
	}

	var hour, minutes int64
	var err error

	max := math.Max(float64(middle-2), 0)
	hour, err = strconv.ParseInt(timeString[int(max):middle], 10, 64)
	if err != nil {
		return nil
	}
	minutes, err = strconv.ParseInt(timeString[int(middle+1):middle+3], 10, 64)
	if err != nil {
		return nil
	}
	if len(timeString) >= middle+4 && byte(timeString[middle+4]) > 0 {
		return nil
	}
	if hour > 60 || minutes > 60 {
		// invalid!
		return nil
	}
	var entryDateTime time.Time
	location, _ := time.LoadLocation(GetAppConf().TimeZone)
	if isNight && hour < 8 {
		// add one day
		targetDay := day.Add(24 * time.Hour)
		entryDateTime = time.Date(targetDay.Year(), targetDay.Month(), targetDay.Day(), int(hour), int(minutes), 0, 0, location)
	} else {
		entryDateTime = time.Date(day.Year(), day.Month(), day.Day(), int(hour), int(minutes), 0, 0, location)
	}
	return &entryDateTime
}
