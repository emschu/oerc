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
	"github.com/gocolly/colly/v2"
	"log"
	"math"
	"regexp"
	"strconv"
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

// SRFParser struct for srf parsing code
type SRFParser struct {
	ParserInterface
	Parser
}

func (s *SRFParser) postProcess() {}

func (s *SRFParser) preProcess() bool {
	s.parallelWorkersCount = 4
	return true
}

// fetchTvShowSRF: This method checks all the tv shows
func (s *SRFParser) fetchTVShows() {
	if !GetAppConf().EnableTVShowCollection || isRecentlyFetched() {
		s.logRecentFetch("Skip update of srf tv shows")
		return
	}

	c := s.newSrfCollector()
	c.OnHTML("section", func(element *colly.HTMLElement) {
		title := trimAndSanitizeString(element.DOM.Find("h2").Text())
		url, _ := element.DOM.Find("a").Attr("href")

		if !srfTvShowURLMatcher.Match([]byte(url)) {
			appLog(fmt.Sprintf("Unexpected url detected: '%s'", url))
			return
		}

		var hash = buildHash([]string{
			fmt.Sprintf("%d", int(s.ChannelFamily.ID)),
			title,
			"tv-show",
		})
		var tvShow = &TvShow{
			ManagedRecord: ManagedRecord{
				Title:           title,
				URL:             srfHostWithPrefix + url,
				Hash:            hash,
				Homepage:        srfHostWithPrefix + url,
				ChannelFamily:   s.ChannelFamily,
				ChannelFamilyID: s.ChannelFamily.ID, // 3 = srf
			},
		}

		show := TvShow{}
		s.db.Where("hash = ?", hash).Find(&show)
		if show.ID != 0 {
			tvShow.ID = show.ID
		}
		tvShow.saveTvShowRecord(s.db)
	})

	err := c.Visit("https://www.srf.ch/play/tv/sendungen")
	if err != nil {
		appLog(fmt.Sprintf("Problem fetching URL 'https://www.srf.ch/play/tv/sendungen' %v.", err))
	}
	c.Wait()
}

// handleDay: handle a single day to fetch SRF program data
func (s *SRFParser) handleDay(channel Channel, day time.Time) {
	queryURL := fmt.Sprintf("https://www.srf.ch/play/v3/api/srf/production/tv-program-guide?date=%02d-%02d-%d&businessUnits=srf&reduced=false&channelId=%s", day.Year(), day.Month(), day.Day(), channel.TechnicalID)

	response, err := doGetRequest(queryURL, map[string]string{
		"accept":          "application/json, text/plain, */*",
		"accept-language": "de-DE,de;q=0.9,en;q=0.8",
	}, 1)
	if response == nil || err != nil {
		errorMessage := fmt.Sprintf("Problem fetching SRF URL '%s' %v.\n", queryURL, err)
		appLog(errorMessage)
		log.Printf(errorMessage)
		return
	}
	var apiObject SrfApiResponse
	jsonErr := json.Unmarshal([]byte(*response), &apiObject)
	if jsonErr != nil {
		errorMessage := fmt.Sprintf("Cannot decode SRF API response to JSON\n")
		appLog(errorMessage)
		log.Printf(errorMessage)
		return
	}

	if len(apiObject.Data) > 0 {
		for _, srfEntry := range apiObject.Data[0].ProgramList {
			var programEntry ProgramEntry
			programEntry.StartDateTime = &srfEntry.StartTime
			programEntry.EndDateTime = &srfEntry.EndTime

			durationInMinutes := int16(math.Abs(programEntry.StartDateTime.Sub(*programEntry.EndDateTime).Minutes()))
			programEntry.DurationMinutes = durationInMinutes

			title := trimAndSanitizeString(fmt.Sprintf("%s", srfEntry.Title))
			if srfEntry.Subtitle != "" {
				title = trimAndSanitizeString(fmt.Sprintf("%s - %s", srfEntry.Title, srfEntry.Subtitle))
			}
			programEntry.Title = title
			programEntry.URL = queryURL
			programEntry.Description = trimAndSanitizeString(fmt.Sprintf(srfEntry.Description))

			hash := buildHash([]string{
				title,
				programEntry.StartDateTime.String(),
				programEntry.EndDateTime.String(),
				strconv.Itoa(int(s.ChannelFamily.ID)),
				strconv.Itoa(int(channel.ID)),
				"program-entry",
			})
			programEntry.Hash = hash
			programEntry.TechnicalID = hash

			programEntry.ChannelFamily = channel.ChannelFamily
			programEntry.ChannelFamilyID = channel.ChannelFamilyID
			programEntry.Channel = channel
			programEntry.ChannelID = channel.ID

			entry := ProgramEntry{}
			s.db.Model(&entry).Where("hash = ?", programEntry.Hash).Where("channel_id = ?", channel.ID).Find(&entry)
			if entry.ID != 0 {
				if entry.isRecentlyUpdated() {
					atomic.AddUint64(&status.TotalSkippedPE, 1)
					return
				}
				programEntry.ID = entry.ID
			}

			if programEntry.EndDateTime == nil || programEntry.EndDateTime.Before(*programEntry.StartDateTime) {
				appLog(fmt.Sprintf("Problem with end or start date in program entry with hash '%s'.", programEntry.Hash))
				return
			}

			if len(srfEntry.Genre) > 0 {
				programEntry.considerTagExists(&srfEntry.Genre)
			}

			if !programEntry.doesImageLinkExist(srfEntry.ImageURL) {
				programEntry.ImageLinks = append(programEntry.ImageLinks, ImageLink{URL: srfEntry.ImageURL})
			}

			programEntry.saveProgramEntryRecord(s.db)
		}
	}
}

func (s *SRFParser) isDateValidToFetch(day *time.Time) (bool, error) {
	if day == nil {
		return false, fmt.Errorf("invalid day")
	}

	if s.isMoreThanXDaysInFuture(day, 42) {
		return false, fmt.Errorf("maximum for days in future for SRF is 42")
	}

	if s.isMoreThanXDaysInPast(day, 15) {
		return false, fmt.Errorf("maximum for days in past for SRF is 15")
	}
	return true, nil
}

// helper method to get a collector instance
func (s *SRFParser) newSrfCollector() *colly.Collector {
	return baseCollector([]string{srfHost})
}

// SrfApiResponse Definition of SRF api response object
type SrfApiResponse struct {
	Data []struct {
		Channel struct {
			Livestream struct {
				ID             string `json:"id"`
				Title          string `json:"title"`
				LivestreamUrn  string `json:"livestreamUrn"`
				ChannelID      string `json:"channelId"`
				PlayableAbroad bool   `json:"playableAbroad"`
			} `json:"livestream"`
			ID           string `json:"id"`
			Vendor       string `json:"vendor"`
			Urn          string `json:"urn"`
			Title        string `json:"title"`
			ImageURL     string `json:"imageUrl"`
			ImageURLRaw  string `json:"imageUrlRaw"`
			ImageTitle   string `json:"imageTitle"`
			Transmission string `json:"transmission"`
			BusinessUnit string `json:"businessUnit"`
			Type         string `json:"type"`
		} `json:"channel"`
		ProgramList []struct {
			Title              string    `json:"title"`
			StartTime          time.Time `json:"startTime"`
			EndTime            time.Time `json:"endTime"`
			Lead               string    `json:"lead,omitempty"`
			Description        string    `json:"description,omitempty"`
			ImageURL           string    `json:"imageUrl"`
			ImageIsFallbackURL bool      `json:"imageIsFallbackUrl"`
			ImageTitle         string    `json:"imageTitle,omitempty"`
			ImageCopyright     string    `json:"imageCopyright,omitempty"`
			MediaUrn           string    `json:"mediaUrn,omitempty"`
			Genre              string    `json:"genre"`
			ProductionYear     int       `json:"productionYear"`
			ProductionCountry  string    `json:"productionCountry"`
			Subtitle           string    `json:"subtitle,omitempty"`
			CreditList         []struct {
				RealName string `json:"realName"`
				Role     string `json:"role"`
			} `json:"creditList,omitempty"`
			SubtitlesAvailable    bool   `json:"subtitlesAvailable"`
			IsLive                bool   `json:"isLive"`
			HasTwoLanguages       bool   `json:"hasTwoLanguages"`
			HasSignLanguage       bool   `json:"hasSignLanguage"`
			HasVisualDescription  bool   `json:"hasVisualDescription"`
			IsFollowUp            bool   `json:"isFollowUp"`
			IsDolbyDigital        bool   `json:"isDolbyDigital"`
			IsRepetition          bool   `json:"isRepetition"`
			RepetitionDescription string `json:"repetitionDescription,omitempty"`
			ChannelTitle          string `json:"channelTitle"`
			ChannelUrn            string `json:"channelUrn"`
			SeasonNumber          int    `json:"seasonNumber,omitempty"`
			EpisodeNumber         int    `json:"episodeNumber,omitempty"`
			EpisodesTotal         int    `json:"episodesTotal,omitempty"`
			HeadlineList          []struct {
				Title       string `json:"title"`
				Description string `json:"description"`
			} `json:"headlineList,omitempty"`
			OriginalTitle        string `json:"originalTitle,omitempty"`
			BroadcastInfo        string `json:"broadcastInfo,omitempty"`
			YouthProtectionColor string `json:"youthProtectionColor,omitempty"`
		} `json:"programList"`
	} `json:"data"`
	IlRequests       []string `json:"ilRequests"`
	FailedIlRequests []any    `json:"failedIlRequests"`
}
