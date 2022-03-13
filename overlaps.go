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
	"gorm.io/gorm"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

// FindOverlaps central method to find, store and process overlapping program entries with same range as "fetch" command
func FindOverlaps(handler dateRangeHandler) {
	log.Println("Start processing overlaps")
	var wg sync.WaitGroup
	for _, channel := range *getChannels() {
		if isChannelFamilyExcluded(&channel.ChannelFamily) {
			continue
		}
		dateRange := handler.getDateRange()
		for _, dayToCheck := range *dateRange {
			wg.Add(1)
			handleOverlapsByDay(&wg, &channel, dayToCheck)
		}
	}
	wg.Wait()
	log.Println("End processing overlaps")
}

// FindOverlapsGlobal recalculate all overlaps in database, could take veeeeery long
func FindOverlapsGlobal() {
	location, _ := time.LoadLocation(GetAppConf().TimeZone)
	if len(getStatusObject().DataStartTime) == 0 || len(getStatusObject().DataEndTime) == 0 {
		return
	}
	dataStartTime, err := parseDate(getStatusObject().DataStartTime, location)
	if err != false {
		return
	}
	dataEndTime, err := parseDate(getStatusObject().DataEndTime, location)
	if err != false {
		return
	}

	times := generateDateRangeBetweenDates(dataStartTime.Add(-24*time.Hour), dataEndTime.Add(24*time.Hour))
	log.Printf("Check %d days for each channel\n", len(*times))

	var wg sync.WaitGroup
	for _, channel := range *getChannels() {
		chn := channel

		for _, dayToCheck := range *times {
			day := dayToCheck
			wg.Add(1)
			go handleOverlapsByDay(&wg, &chn, day)
		}
	}
	wg.Wait()
}

// central method to find, store and process overlapping program entries
func handleOverlapsByDay(wg *sync.WaitGroup, channel *Channel, day time.Time) {
	defer wg.Done()
	var dailyProgramEntries []ProgramEntry

	db := getDb()
	collisionMap := findOverlaps(db, channel, day, &dailyProgramEntries)
	if len(*collisionMap) == 0 {
		return
	}

	overlapIDs := storeOverlaps(db, collisionMap)

	log.Printf("Processing %d program entry collisions of channel '%s' on day '%s'.\n", len(overlapIDs), channel.Title, day.Format("2006-01-02"))

	for _, programEntryID := range overlapIDs {
		processOverlaps(db, &programEntryID)
	}
}

func processOverlaps(db *gorm.DB, programEntryID *uint) {
	var programEntry ProgramEntry
	db.Model(&ProgramEntry{}).Preload("CollisionEntries").Preload("Channel").Find(&programEntry, programEntryID)
	if programEntry.ID == 0 {
		log.Printf("Warning: Could not find program entry with id #%d. Skipping...\n", *programEntryID)
		return
	}
	if len(programEntry.CollisionEntries) == 0 {
		return
	}
	// ensure the correct items only are marked as deprecated here

	// assume the current entry is the most recent one, until we know it better from collision entries
	isEntryDeprecated := isProgramEntryPossiblyDeprecated(&programEntry)
	// update entry - if needed
	if isEntryDeprecated != programEntry.IsDeprecated {
		programEntry.IsDeprecated = isEntryDeprecated
		if verboseGlobal {
			if isEntryDeprecated {
				log.Printf("Set program entry #%d as deprecated.\n", programEntry.ID)
			} else {
				log.Printf("Set program entry #%d as NOT deprecated.\n", programEntry.ID)
			}
		}
		now := time.Now()
		programEntry.LastCollisionCheck = &now
		db.Save(&programEntry)
	}
}

func isProgramEntryPossiblyDeprecated(programEntry *ProgramEntry) bool {
	if programEntry.IsDeprecated {
		// TODO what to do?
		return true
	}

	const changeTimeTolerance = 3600 // one hour
	isEntryDeprecated := false
	for _, collisionEntry := range programEntry.CollisionEntries {
		if collisionEntry.IsDeprecated {
			// skip overlapping entries which are already deprecated
			continue
		}
		createdAtDiffSecs := int64(programEntry.CreatedAt.Sub(collisionEntry.CreatedAt) / time.Second)
		var lastCheckAtDiffSecs int64
		if programEntry.LastCheck == nil || programEntry.LastCheck.IsZero() {
			lastCheckAtDiffSecs = 0
		} else {
			lastCheckAtDiffSecs = int64(programEntry.LastCheck.Sub(*collisionEntry.LastCheck) / time.Second)
		}
		if math.Abs(float64(createdAtDiffSecs)) < changeTimeTolerance && math.Abs(float64(lastCheckAtDiffSecs)) < changeTimeTolerance {
			// we cannot decide something here
			continue
		}
		if math.Abs(float64(lastCheckAtDiffSecs)) < changeTimeTolerance {
			// decide by created at diff only
			if programEntry.CreatedAt.Before(collisionEntry.CreatedAt) {
				isEntryDeprecated = true
				break
			}
		}
		if math.Abs(float64(createdAtDiffSecs)) < changeTimeTolerance {
			// decide by last checked at diff only
			if programEntry.LastCheck.Before(*collisionEntry.LastCheck) {
				isEntryDeprecated = true
				break
			}
		}
		if programEntry.LastCheck.Before(*collisionEntry.LastCheck) {
			isEntryDeprecated = true
			break
		}
	}
	return isEntryDeprecated
}

func storeOverlaps(db *gorm.DB, collisionMap *map[uint][]uint) []uint {
	var affectedIds = make([]uint, 0)

	tx := db.Session(&gorm.Session{PrepareStmt: true})

	for programEntryID, collisions := range *collisionMap {
		var programEntry ProgramEntry
		tx.Model(ProgramEntry{}).Preload("Channel").Preload("CollisionEntries").Find(&programEntry, programEntryID)
		if programEntry.ID == 0 {
			// entry not found
			log.Printf("Warning: Could not fetch program entry record #%d\n", programEntryID)
			continue
		}
		var relatedIDs = make([]uint, 0)
		for _, collisionEntryID := range collisions {
			relatedIDs = append(relatedIDs, collisionEntryID)
		}
		var relatedItems []ProgramEntry
		if len(relatedIDs) > 0 {
			db.Model(&ProgramEntry{}).Where("id IN(?)", relatedIDs).Find(&relatedItems)
		}
		if len(relatedItems) > 0 {
			programEntry.CollisionEntries = relatedItems
			db.Save(&programEntry)
		}
		affectedIds = append(affectedIds, programEntry.ID)
	}
	return affectedIds
}

// this method asks the database for overlapping items and store the results in a map
func findOverlaps(db *gorm.DB, channel *Channel, day time.Time, dailyProgramEntries *[]ProgramEntry) *map[uint][]uint {
	startOfDay := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	endOfDay := time.Date(day.Year(), day.Month(), day.Day(), 23, 59, 59, 0, day.Location())

	// fetch affected program entries
	db.Model(ProgramEntry{}).
		Where("channel_id = ? AND start_date_time between ? AND ?", channel.ID, startOfDay, endOfDay).
		Order("start_date_time ASC").
		Find(&dailyProgramEntries)

	var counter = 0
	var queries []string
	for _, peToCheck := range *dailyProgramEntries {
		for _, pe := range *dailyProgramEntries {
			if peToCheck.ID == pe.ID {
				// avoid checking overlap with itself
				continue
			}
			if peToCheck.EndDateTime.Before(*pe.StartDateTime) || peToCheck.StartDateTime.After(*pe.EndDateTime) {
				// exclude trivial cases from db queries and save a lot of time and memory
				continue
			}
			queries = append(queries, fmt.Sprintf("CASE WHEN (SELECT (timestamptz '%s', timestamptz '%s') "+
				"OVERLAPS (timestamptz '%s', timestamptz '%s'))=TRUE THEN '%d;%d' ELSE '' END as a%d",
				peToCheck.StartDateTime.Format(time.RFC3339), peToCheck.EndDateTime.Format(time.RFC3339),
				pe.StartDateTime.Format(time.RFC3339), pe.EndDateTime.Format(time.RFC3339), peToCheck.ID, pe.ID, counter))
			counter++
		}
	}

	tx := db.Session(&gorm.Session{PrepareStmt: true})
	var collisionMap = make(map[uint][]uint)
	chunks := chunkStringSlice(queries, 50)
	for _, singleChunk := range chunks {
		var result []map[string]interface{}
		tx.Raw(fmt.Sprintf("SELECT %s", strings.Join(singleChunk, ","))).Scan(&result)

		if len(result) == 0 || len(result[0]) == 0 {
			continue
		}
		for _, mapEntry := range result[0] {
			collision := mapEntry.(string)
			if len(collision) > 0 {
				entries := strings.Split(collision, ";")
				if len(entries) != 2 {
					log.Fatalf("Problem with formatting of collision check response\n")
				}
				collisionID, err := strconv.ParseUint(entries[1], 10, 64)
				if err != nil {
					log.Printf("Error parsing collision sql response")
					continue
				}
				programEntryID, err := strconv.ParseUint(entries[0], 10, 64)
				if err != nil {
					log.Printf("Error parsing collision sql response")
					continue
				}
				collisionMap[uint(programEntryID)] = append(collisionMap[uint(programEntryID)], uint(collisionID))
			}
		}
	}
	return &collisionMap
}
