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
	"log"
	"sync"
	"time"

	"gorm.io/gorm"
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
	currentStatusData := getStatusObject()
	if currentStatusData == nil {
		log.Fatalln("Error fetching application status data from database")
	}
	if len(currentStatusData.DataStartTime) == 0 || len(currentStatusData.DataEndTime) == 0 {
		return
	}
	dataStartTime, err := parseDate(currentStatusData.DataStartTime, location)
	if err != false {
		return
	}
	dataEndTime, err := parseDate(currentStatusData.DataEndTime, location)
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

	overlapIDs := storeCollisions(db, collisionMap)

	if isDebug() {
		log.Printf("Processing %d program entry collisions of channel '%s' on day '%s'.\n", len(overlapIDs), channel.Title, day.Format("2006-01-02"))
	}

	for programEntryID, _ := range *collisionMap {
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

	// update entry - if needed
	programEntry.IsDeprecated = true
	if verboseGlobal {
		log.Printf("Set program entry #%d as deprecated.\n", programEntry.ID)
	}
	now := time.Now()
	programEntry.LastCollisionCheck = &now
	db.Save(&programEntry)
}

// overlapResult struct to map SQL query results for overlapping program entries
type overlapResult struct {
	ChannelID  uint `gorm:"column:channel_id"`
	Program1ID uint `gorm:"column:program1_id"`
	Program2ID uint `gorm:"column:program2_id"`
}

// This function stores overlap information (= collision entries) in the ProgramEntry table
func storeCollisions(db *gorm.DB, collisionMap *map[uint][]uint) []uint {
	var affectedIDs = make([]uint, 0)

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
		affectedIDs = append(affectedIDs, programEntry.ID)
	}
	return affectedIDs
}

// this method asks the database for overlapping items and store the results in a map
func findOverlaps(db *gorm.DB, channel *Channel, day time.Time, dailyProgramEntries *[]ProgramEntry) *map[uint][]uint {
	startOfDay := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	endOfDay := time.Date(day.Year(), day.Month(), day.Day(), 23, 59, 59, 0, day.Location())

	// execute raw SQL query to find overlaps with duration calculation
	var overlaps []overlapResult
	sqlQuery := `
		SELECT
			a.channel_id,
			a.id as program1_id,
			b.id as program2_id
		FROM program_entries a
		JOIN program_entries b ON a.channel_id = b.channel_id
		WHERE a.id < b.id -- avoid duplicates
		  AND a.channel_id = ?
		  AND date_trunc('minute', a.start_date_time) BETWEEN ? AND ?
		  AND date_trunc('minute',a.start_date_time) < date_trunc('minute',b.end_date_time)
		  AND date_trunc('minute',a.end_date_time) > date_trunc('minute',b.start_date_time)
		ORDER BY a.channel_id, a.start_date_time
	`

	result := db.Raw(sqlQuery, channel.ID, startOfDay, endOfDay).Scan(&overlaps)
	if result.Error != nil {
		log.Printf("Error executing overlap query for channel %d on day %s: %v\n", channel.ID, day.Format("2006-01-02"), result.Error)
		collisionMap := make(map[uint][]uint)
		return &collisionMap
	}

	// build collision map from query results
	var collisionMap = make(map[uint][]uint)
	for _, overlap := range overlaps {
		collisionMap[overlap.Program1ID] = append(collisionMap[overlap.Program1ID], overlap.Program2ID)
	}

	return &collisionMap
}
