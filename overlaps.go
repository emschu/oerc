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
	"github.com/alitto/pond"
	"log"
	"math"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// FindOverlaps central method to find, store and process overlapping program entries with same range as "fetch" command
func FindOverlaps() {
	for _, channel := range *getChannels() {
		if isChannelFamilyExcluded(&channel.ChannelFamily) {
			continue
		}
		times := generateDateRangeInPastAndFuture(GetAppConf().DaysInPast-1, GetAppConf().DaysInFuture+1)
		for _, dayToCheck := range *times {
			handleOverlaps(&channel, dayToCheck)
		}
	}
}

// FindOverlapsGlobal recalcalate all overlaps in database
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

	var pool = pond.New(int(math.RoundToEven(float64(runtime.NumCPU())*1.5)), 100, getWorkerPoolIdleTimeout(), pond.PanicHandler(func(i interface{}) {
		log.Printf("Problem with goroutine pool: %v\n", i)
	}))
	times := generateDateRangeBetweenDates(dataStartTime.Add(-24*time.Hour), dataEndTime.Add(24*time.Hour))
	log.Printf("Check %d days for each channel\n", len(*times))
	for _, channel := range *getChannels() {
		chn := channel
		pool.Submit(func() {
			for _, dayToCheck := range *times {
				handleOverlaps(&chn, dayToCheck)
			}
		})
	}

	// wait for finish
	pool.StopAndWait()
}

// central method to find, store and process overlapping program entries
func handleOverlaps(channel *Channel, day time.Time) {
	var dailyProgramEntries []ProgramEntry

	collisionMap := findOverlaps(channel, day, &dailyProgramEntries)
	if len(collisionMap) == 0 {
		return
	}

	overlapIDs := storeOverlaps(&collisionMap)

	log.Printf("Processing %d program entry collisions of channel '%s' on day '%s'.\n", len(overlapIDs), channel.Title, day.Format("2006-01-02"))

	for _, programEntryID := range overlapIDs {
		processOverlaps(&programEntryID)
	}
}

func processOverlaps(programEntryID *uint) {
	db := getDb()
	var programEntryWithCollisions ProgramEntry
	db.Model(&ProgramEntry{}).Preload("CollisionEntries").Find(&programEntryWithCollisions, programEntryID)
	if programEntryWithCollisions.ID == 0 {
		log.Printf("Warning: Could not find program entry with id #%d. Skipping...\n", *programEntryID)
		return
	}
	if len(programEntryWithCollisions.CollisionEntries) == 0 {
		return
	}
	// assume the current entry is the most recent one, until we know it better from collision entries
	isEntryDeprecated := isProgramEntryPossiblyDeprecated(&programEntryWithCollisions)
	// update entry - if needed
	if isEntryDeprecated != programEntryWithCollisions.IsDeprecated {
		programEntryWithCollisions.IsDeprecated = isEntryDeprecated
		if verboseGlobal {
			if isEntryDeprecated {
				log.Printf("Set program entry #%d as deprecated.\n", programEntryWithCollisions.ID)
			} else {
				log.Printf("Set program entry #%d as NOT deprecated.\n", programEntryWithCollisions.ID)
			}
		}
		now := time.Now()
		programEntryWithCollisions.LastCollisionCheck = &now
		db.Save(&programEntryWithCollisions)
	}
}

func isProgramEntryPossiblyDeprecated(programEntryWithCollisions *ProgramEntry) bool {
	if programEntryWithCollisions.IsDeprecated {
		// TODO what to do?
		return true
	}

	const changeTimeTolerance = 3600 // one hour
	isEntryDeprecated := false
	for _, collisionEntry := range programEntryWithCollisions.CollisionEntries {
		if collisionEntry.IsDeprecated {
			// skip overlapping entries which are already deprecated
			continue
		}
		createdAtDiffSecs := int64(programEntryWithCollisions.CreatedAt.Sub(collisionEntry.CreatedAt) / time.Second)
		lastCheckAtDiffSecs := int64(programEntryWithCollisions.LastCheck.Sub(*collisionEntry.LastCheck) / time.Second)
		if math.Abs(float64(createdAtDiffSecs)) < changeTimeTolerance && math.Abs(float64(lastCheckAtDiffSecs)) < changeTimeTolerance {
			// we cannot decide something here
			continue
		}
		if math.Abs(float64(lastCheckAtDiffSecs)) < changeTimeTolerance {
			// decide by created at diff only
			if programEntryWithCollisions.CreatedAt.Before(collisionEntry.CreatedAt) {
				isEntryDeprecated = true
				break
			}
			continue
		}
		if math.Abs(float64(createdAtDiffSecs)) < changeTimeTolerance {
			// decide by last checked at diff only
			if programEntryWithCollisions.LastCheck.Before(*collisionEntry.LastCheck) {
				isEntryDeprecated = true
				break
			} else {
				continue
			}
		}
		if programEntryWithCollisions.LastCheck.Before(*collisionEntry.LastCheck) {
			isEntryDeprecated = true
			break
		}
	}
	return isEntryDeprecated
}

func storeOverlaps(collisionMap *map[uint][]uint) []uint {
	var affectedIds = make([]uint, 0)
	db := getDb()

	for programEntryID, collisions := range *collisionMap {
		var programEntry ProgramEntry
		db.Model(ProgramEntry{}).Preload("CollisionEntries").Find(&programEntry, programEntryID)
		if programEntry.ID == 0 {
			// entry not found
			log.Printf("Warning: Could not fetch program entry record #%d\n", programEntryID)
			continue
		}
		// Clean up existing collisions
		if len(programEntry.CollisionEntries) > 0 {
			// clear old collision entries
			err := getDb().Model(&programEntry).Association("CollisionEntries").Clear()
			if err != nil {
				log.Printf("Warning: Problem with deleting collision entries: %v\n", err)
				continue
			}
		}
		programEntry.CollisionEntries = make([]ProgramEntry, 0)
		for _, collisionEntryID := range collisions {
			if collisionEntryID != 0 {
				var linkedProgramEntry ProgramEntry
				linkedProgramEntry.ID = collisionEntryID
				programEntry.CollisionEntries = append(programEntry.CollisionEntries, linkedProgramEntry)
			}
		}
		db.Save(&programEntry)
		affectedIds = append(affectedIds, programEntry.ID)
	}
	return affectedIds
}

func findOverlaps(channel *Channel, day time.Time, dailyProgramEntries *[]ProgramEntry) map[uint][]uint {
	db := getDb()
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
			queries = append(queries, fmt.Sprintf("CASE WHEN (SELECT (timestamptz '%s', timestamptz '%s') "+
				"OVERLAPS (timestamptz '%s', timestamptz '%s'))=TRUE THEN '%d;%d' ELSE '' END as a%d",
				peToCheck.StartDateTime.Format(time.RFC3339), peToCheck.EndDateTime.Format(time.RFC3339),
				pe.StartDateTime.Format(time.RFC3339), pe.EndDateTime.Format(time.RFC3339), peToCheck.ID, pe.ID, counter))
			counter++
		}
	}

	var collisionMap = make(map[uint][]uint)
	chunks := chunkStringSlice(queries, 25)
	for _, singleChunk := range chunks {
		var result []map[string]interface{}
		db.Raw(fmt.Sprintf("SELECT %s", strings.Join(singleChunk, ","))).Scan(&result)

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
					log.Printf("Error")
				}
				collisionMap[uint(programEntryID)] = append(collisionMap[uint(programEntryID)], uint(collisionID))
			}
		}
	}
	return collisionMap
}
