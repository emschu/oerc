// oerc, alias oer-collector
// Copyright (C) 2021-2025 emschu[aet]mailbox.org
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
	"gorm.io/gorm"
	"log"
	"strings"
	"time"
)

// uniqueStrSlice: function to make a slice of program entries unique
func uniqueStrSlice(stringSlice *[]string) []string {
	keys := make(map[string]bool)
	var res = make([]string, 0)
	for _, entry := range *stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			res = append(res, entry)
		}
	}
	return res
}

// uniqueProgramEntryList: function to make a slice of program entries unique
func uniqueProgramEntryList(programEntrySlice *[]ProgramEntry) []ProgramEntry {
	keys := make(map[uint]bool)
	var res = make([]ProgramEntry, 0)
	for _, entry := range *programEntrySlice {
		if _, value := keys[entry.ID]; !value {
			keys[entry.ID] = true
			res = append(res, entry)
		}
	}
	return res
}

// SearchProgram method to check the program data of the next days for user-defined keywords
func SearchProgram() {
	now := time.Now()
	tStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endDate := now.Add(time.Duration(GetAppConf().SearchDaysInFuture+1) * 24 * time.Hour)
	tEnd := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 3, 0, 0, 0, now.Location())

	db := getDb()
	var excludedChannelIDs = getExcludedChannelsFromSearch(db)

	var programEntryList = make([]ProgramEntry, 0)
	programResponse := getProgramOf(&tStart, &tEnd, nil)

	if programResponse == nil {
		log.Fatalf("Could not retrieve program for date range '%s'-'%s'. Please fetch the program at first.\n", tStart, tEnd)
		return
	}

	skipCounter := 0
	for _, programEntry := range programEntryList {
		// exclude the channels found above
		if isChannelExcluded(excludedChannelIDs, &programEntry) {
			skipCounter++
			continue
		}
		var keywords = make([]string, 0)
		for _, searchWord := range GetAppConf().SearchKeywords {
			if strings.Contains(strings.ToLower(programEntry.Title), strings.ToLower(searchWord)) {
				programEntryList = append(programEntryList, programEntry)
				keywords = append(keywords, searchWord)
			}
			if strings.Contains(strings.ToLower(programEntry.Description), strings.ToLower(searchWord)) {
				programEntryList = append(programEntryList, programEntry)
				keywords = append(keywords, searchWord)
			}
			if strings.Contains(strings.ToLower(programEntry.Tags), strings.ToLower(searchWord)) {
				programEntryList = append(programEntryList, programEntry)
				keywords = append(keywords, searchWord)
			}
		}

		programEntryList = uniqueProgramEntryList(&programEntryList)

		if len(keywords) > 0 {
			keywords = uniqueStrSlice(&keywords)

			var reccEntry Recommendation
			db.Model(&Recommendation{}).Where("program_entry_id = ?", programEntry.ID).First(&reccEntry)

			now := time.Now()
			reccEntry.CreatedAt = &now
			reccEntry.ProgramEntry = programEntry
			reccEntry.ProgramEntryID = programEntry.ID
			reccEntry.ChannelID = programEntry.ChannelID
			reccEntry.StartDateTime = programEntry.StartDateTime
			reccEntry.Keywords = strings.Join(keywords, ",")

			if reccEntry.ID > 0 {
				db.Model(&reccEntry).Updates(map[string]interface{}{
					"keywords":        reccEntry.Keywords,
					"start_date_time": programEntry.StartDateTime,
					"channel_id":      programEntry.ChannelID,
				})
			} else {
				db.Create(&reccEntry)
			}
		}
	}

	log.Printf("Found %v search results in tv program of today + %d days. Searched in %d program entries.\n", len(programEntryList), GetAppConf().SearchDaysInFuture, len(*programResponse)-skipCounter)
}

func getExcludedChannelsFromSearch(db *gorm.DB) *[]uint {
	var excludedChannelIds = make([]uint, len(GetAppConf().SearchSkipChannels))
	for _, channel := range GetAppConf().SearchSkipChannels {
		var channelEntry Channel
		db.Model(&Channel{}).Where("is_deprecated is false AND title = ?", channel).First(&channelEntry)
		if channelEntry.ID > 0 {
			excludedChannelIds = append(excludedChannelIds, channelEntry.ID)
		} else {
			log.Printf("Problem with channel name '%s'. Could not find it in database. Skipping this entry.\n", channel)
		}
	}
	return &excludedChannelIds
}

func isChannelFamilyExcluded(family *ChannelFamily) bool {
	if !appConf.EnableARD && family.Title == "ARD" {
		return true
	}
	if !appConf.EnableZDF && family.Title == "ZDF" {
		return true
	}
	if !appConf.EnableORF && family.Title == "ORF" {
		return true
	}
	if !appConf.EnableSRF && family.Title == "SRF" {
		return true
	}
	return false
}

// method to check if a single program entry should be skipped because its channel id is contained in the first parameter
func isChannelExcluded(excludedChannelIds *[]uint, programEntry *ProgramEntry) bool {
	for _, skippedChannelID := range *excludedChannelIds {
		if programEntry.ChannelID == skippedChannelID {
			return true
		}
	}
	return false
}
