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
	"log"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"
)

// FindOverlaps central method to find, store and process overlapping program entries with same range as "fetch" command
func FindOverlaps(handler dateRangeHandler) {
	log.Println("Start processing overlaps")

	channels, err := getChannels()
	if err != nil {
		log.Printf("Error fetching channels for the overlap processing: %v\n", err)
		return
	}
	var wg sync.WaitGroup
	for _, channel := range *channels {
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
	setupMaterializedView()
	location, _ := time.LoadLocation(GetAppConf().TimeZone)
	currentStatusData, err := getStatusResponse()
	if err != nil || currentStatusData == nil {
		log.Fatalln("Error fetching application status data from database")
	}
	if len(currentStatusData.DataStartTime) == 0 || len(currentStatusData.DataEndTime) == 0 {
		return
	}
	dataStartTime, fail := parseDate(currentStatusData.DataStartTime, location)
	if fail {
		return
	}
	dataEndTime, fail := parseDate(currentStatusData.DataEndTime, location)
	if fail {
		return
	}

	times := generateDateRangeBetweenDates(dataStartTime.Add(-24*time.Hour), dataEndTime.Add(24*time.Hour))
	log.Printf("Check %d days for each channel\n", len(*times))

	channels, err := getChannels()
	if err != nil {
		log.Printf("Error fetching channels for the overlap processing: %v\n", err)
		return
	}
	var wg sync.WaitGroup
	for _, channel := range *channels {
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

	db := getDb()
	collisionMap := findOverlaps(db, channel, day)
	// clear collision state of program entries of this day without current overlaps
	clearStaleCollisions(db, channel, day, collisionMap)
	if len(*collisionMap) == 0 {
		return
	}

	overlapIDs := storeCollisions(db, collisionMap)

	if isDebug() {
		log.Printf("Processing %d program entry collisions of channel '%s' on day '%s'.\n", len(overlapIDs), channel.Title, day.Format("2006-01-02"))
	}

	for programEntryID := range *collisionMap {
		processOverlaps(db, &programEntryID)
	}

}

// clearStaleCollisions removes collision entries and the deprecation status of
// program entries of the given channel and day which do not have any current
// overlaps anymore, so that the overlap state cannot become stale
func clearStaleCollisions(db *gorm.DB, channel *Channel, day time.Time, currentOverlaps *map[uint][]uint) {
	startOfDay := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	endOfDay := time.Date(day.Year(), day.Month(), day.Day(), 23, 59, 59, 0, day.Location())

	var entries []ProgramEntry
	result := db.Model(&ProgramEntry{}).
		Where("channel_id = ? AND start_date_time BETWEEN ? AND ?", channel.ID, startOfDay, endOfDay).
		Preload("CollisionEntries").
		Find(&entries)
	if result.Error != nil {
		log.Printf("Error fetching program entries of channel %d on day '%s' to clear stale collisions: %v\n",
			channel.ID, day.Format("2006-01-02"), result.Error)
		return
	}

	for _, entry := range entries {
		if _, hasCurrentOverlap := (*currentOverlaps)[entry.ID]; hasCurrentOverlap {
			continue
		}
		// clear join table entries of collisions which do not exist anymore
		if len(entry.CollisionEntries) > 0 {
			if err := db.Model(&entry).Association("CollisionEntries").Clear(); err != nil {
				log.Printf("Error clearing collision entries of program entry #%d: %v\n", entry.ID, err)
				continue
			}
		}
		if entry.IsDeprecated {
			if err := db.Model(&entry).Updates(map[string]any{"is_deprecated": false, "last_collision_check": time.Now()}).Error; err != nil {
				log.Printf("Error clearing deprecation status of program entry #%d: %v\n", entry.ID, err)
			}
		}
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
			// replace the join table entries, so that collisions which do not
			// exist anymore are removed as well
			if err := db.Model(&programEntry).Association("CollisionEntries").Replace(relatedItems); err != nil {
				log.Printf("Error storing collisions of program entry #%d: %v\n", programEntry.ID, err)
				continue
			}
		}
		affectedIDs = append(affectedIDs, programEntry.ID)
	}
	return affectedIDs
}

// this method asks the database for overlapping items and stores the results in a map
func findOverlaps(db *gorm.DB, channel *Channel, day time.Time) *map[uint][]uint {
	startOfDay := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
	endOfDay := time.Date(day.Year(), day.Month(), day.Day(), 23, 59, 59, 0, day.Location())

	// execute raw SQL query to find overlaps with duration calculation
	var overlaps []overlapResult
	var sqlQuery string
	if strings.ToLower(GetAppConf().DbType) == "postgres" {
		sqlQuery = `
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
	}
	if strings.ToLower(GetAppConf().DbType) == "sqlite" {
		sqlQuery = `
		SELECT
		  a.channel_id,
		  a.id AS program1_id,
		  b.id AS program2_id
		FROM program_entries AS a
		JOIN program_entries AS b
		  ON a.channel_id = b.channel_id
		WHERE a.id < b.id
		  AND a.channel_id = ?
		  AND strftime('%Y-%m-%d %H:%M:00', a.start_date_time) BETWEEN ? AND ?
		  AND strftime('%Y-%m-%d %H:%M:00', a.start_date_time) < strftime('%Y-%m-%d %H:%M:00', b.end_date_time)
		  AND strftime('%Y-%m-%d %H:%M:00', a.end_date_time)   > strftime('%Y-%m-%d %H:%M:00', b.start_date_time)
		ORDER BY a.channel_id, a.start_date_time;
`
	}

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
