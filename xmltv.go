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
	"encoding/xml"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// XMLTV Document structure
type Document struct {
	Channel   []ChannelItem   `xml:"channel"`
	Programme []ProgrammeItem `xml:"programme"`
}

type ChannelItem struct {
	ID          string              `xml:"id,attr"`
	DisplayName []LocalizedTextNode `xml:"display-name"`
	Icon        *Icon               `xml:"icon"`
	URL         []URLNode           `xml:"url"`
}

type ProgrammeItem struct {
	Channel   string  `xml:"channel,attr"`
	Clumpidx  *string `xml:"clumpidx,attr"`
	PdcStart  *string `xml:"pdc-start,attr"`
	Showview  *int    `xml:"showview,attr"`
	Start     string  `xml:"start,attr"`
	Stop      *string `xml:"stop,attr"`
	Videoplus *int    `xml:"videoplus,attr"`
	VpsStart  *string `xml:"vps-start,attr"`
	Audio     *struct {
		Present string `xml:"present"`
		Stereo  string `xml:"stereo"`
	} `xml:"audio"`
	Category []LocalizedTextNode `xml:"category"`
	Country  *string             `xml:"country"`
	Credits  *struct {
		Actor       []ProgrammeItemActor `xml:"actor"`
		Adapter     string               `xml:"adapter"`
		Commentator string               `xml:"commentator"`
		Composer    string               `xml:"composer"`
		Director    string               `xml:"director"`
		Editor      string               `xml:"editor"`
		Guest       []string             `xml:"guest"`
		Presenter   string               `xml:"presenter"`
		Producer    string               `xml:"producer"`
		Writer      string               `xml:"writer"`
	} `xml:"credits"`
	Date       *int                `xml:"date"`
	Desc       []LocalizedTextNode `xml:"desc"`
	EpisodeNum *URLNode            `xml:"episode-num"`
	Icon       *Icon               `xml:"icon"`
	Image      []struct {
		Orient   string `xml:"orient,attr"`
		Size     int    `xml:"size,attr"`
		System   string `xml:"system,attr"`
		Type     string `xml:"type,attr"`
		CharData string `xml:",chardata"`
	} `xml:"image"`
	Keyword    []LocalizedTextNode `xml:"keyword"`
	Language   *string             `xml:"language"`
	LastChance *LocalizedTextNode  `xml:"last-chance"`
	Length     *struct {
		Units    string `xml:"units,attr"`
		CharData string `xml:",chardata"`
	} `xml:"length"`
	New             *struct{}          `xml:"new"`
	OrigLanguage    *LocalizedTextNode `xml:"orig-language"`
	Premiere        *string            `xml:"premiere"`
	PreviouslyShown *struct {
		Channel string `xml:"channel,attr"`
		Start   string `xml:"start,attr"`
	} `xml:"previously-shown"`
	Rating []Rating `xml:"rating"`
	Review []struct {
		Lang     string `xml:"lang,attr"`
		Reviewer string `xml:"reviewer,attr"`
		Source   string `xml:"source,attr"`
		Type     string `xml:"type,attr"`
		CharData string `xml:",chardata"`
	} `xml:"review"`
	StarRating []Rating           `xml:"star-rating"`
	SubTitle   *LocalizedTextNode `xml:"sub-title"`
	Subtitles  []struct {
		Type     string            `xml:"type,attr"`
		Language LocalizedTextNode `xml:"language"`
	} `xml:"subtitles"`
	Title LocalizedTextNode `xml:"title"`
	URL   []URLNode         `xml:"url"`
	Video *struct {
		Aspect  string `xml:"aspect"`
		Colour  string `xml:"colour"`
		Present string `xml:"present"`
		Quality string `xml:"quality"`
	} `xml:"video"`
}

type ProgrammeItemActor struct {
	Guest    string `xml:"guest,attr"`
	Role     string `xml:"role,attr"`
	CharData string `xml:",chardata"`
	Image    struct {
		Type     string `xml:"type,attr"`
		CharData string `xml:",chardata"`
	} `xml:"image"`
	URL URLNode `xml:"url"`
}

type LocalizedTextNode struct {
	Lang     *string `xml:"lang,attr"`
	CharData string  `xml:",chardata"`
}

type Icon struct {
	Height int    `xml:"height,attr,omitempty"`
	Src    string `xml:"src,attr"`
	Width  int    `xml:"width,attr,omitempty"`
}

type URLNode struct {
	System   string `xml:"system,attr,omitempty"`
	CharData string `xml:",chardata"`
}

type Rating struct {
	System string `xml:"system,attr"`
	Icon   struct {
		Src string `xml:"src,attr"`
	} `xml:"icon"`
	Value string `xml:"value"`
}

// newXMLTvDocument exports program data to XMLTV format
func newXMLTvDocument(channelList *[]Channel, programEntryList *[]ProgramEntry) (*Document, error) {
	doc := Document{
		// Pre-allocate slices to avoid reallocations
		Channel:   make([]ChannelItem, 0, len(*channelList)),
		Programme: make([]ProgrammeItem, 0, len(*programEntryList)),
	}

	if len(*channelList) > 5 {
		channelItems := make(chan ChannelItem, len(*channelList))

		for _, ch := range *channelList {
			go func(c Channel) {
				channelItems <- mapChannelToXMLTV(c)
			}(ch)
		}

		for i := 0; i < len(*channelList); i++ {
			doc.Channel = append(doc.Channel, <-channelItems)
		}
	} else {
		for _, ch := range *channelList {
			doc.Channel = append(doc.Channel, mapChannelToXMLTV(ch))
		}
	}

	for _, pe := range *programEntryList {
		programmeItem := mapProgramEntryToXMLTV(&pe)
		doc.Programme = append(doc.Programme, programmeItem)
	}
	return &doc, nil
}

func mapChannelToXMLTV(ch Channel) ChannelItem {
	// Create the ChannelItem
	channelItem := ChannelItem{
		ID: buildChannelID(ch),
		DisplayName: []LocalizedTextNode{
			{
				CharData: ch.Title,
			},
		},
	}

	if ch.Homepage != "" {
		channelItem.URL = append(channelItem.URL, URLNode{
			CharData: ch.Homepage,
		})
	}
	return channelItem
}

func buildChannelID(ch Channel) string {
	return fmt.Sprintf("%s: %s", ch.ChannelFamily.getXMLTvChannelPrefix(), ch.Title)
}

func mapProgramEntryToXMLTV(pe *ProgramEntry) ProgrammeItem {
	startTime := ""
	var stopTime *string

	if pe.StartDateTime != nil {
		startTime = pe.StartDateTime.Format("20060102150405 -0700")
	}

	if pe.EndDateTime != nil {
		formattedStopTime := pe.EndDateTime.Format("20060102150405 -0700")
		stopTime = &formattedStopTime
	}

	programmeItem := ProgrammeItem{
		Channel: buildChannelID(pe.Channel),
		Start:   startTime,
		Stop:    stopTime,
		Title: LocalizedTextNode{
			CharData: pe.Title,
		},
	}

	if pe.Description != "" {
		programmeItem.Desc = append(programmeItem.Desc, LocalizedTextNode{
			CharData: pe.Description,
		})
	}

	if pe.DurationMinutes > 0 {
		lengthStr := strconv.Itoa(int(pe.DurationMinutes))
		programmeItem.Length = &struct {
			Units    string `xml:"units,attr"`
			CharData string `xml:",chardata"`
		}{
			Units:    "minutes",
			CharData: lengthStr,
		}
	}

	if len(pe.ImageLinks) > 0 {
		programmeItem.Icon = &Icon{
			Src: pe.ImageLinks[0].URL,
		}
	}

	if pe.Tags != "" {
		tags := strings.Split(pe.Tags, ";")
		programmeItem.Category = make([]LocalizedTextNode, 0, len(tags))
		for _, tag := range tags {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				programmeItem.Category = append(programmeItem.Category, LocalizedTextNode{
					CharData: tag,
				})
			}
		}
	}

	return programmeItem
}

func writeXMLTVToString(doc *Document) (string, error) {
	var stringWriter strings.Builder
	encoder := xml.NewEncoder(&stringWriter)
	encoder.Indent("", "  ")
	if err := encoder.Encode(doc); err != nil {
		return "", fmt.Errorf("error encoding XMLTV document: %v", err)
	}
	return stringWriter.String(), nil
}

// exportToXMLTV exports program data to XMLTV format
func exportToXMLTV(rangeStart, rangeEnd time.Time) (*Document, error) {
	channels := getChannels()
	if channels == nil {
		return nil, fmt.Errorf("error fetching channels")
	}

	if verboseGlobal {
		log.Printf("Fetching program entries for %d channels from %s to %s", len(*channels),
			rangeStart.Format("2006-01-02"), rangeEnd.Format("2006-01-02"))
	}

	programEntries, err := getAllProgramEntriesOf(rangeStart, rangeEnd)
	if err != nil {
		return nil, fmt.Errorf("error fetching program entries: %v", err)
	}

	return newXMLTvDocument(channels, programEntries)
}
