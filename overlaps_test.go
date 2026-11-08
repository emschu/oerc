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
	"sync"
	"testing"
	"time"
)

func createOverlapTestEntry(t *testing.T, channel *Channel, start time.Time, duration time.Duration) *ProgramEntry {
	t.Helper()
	end := start.Add(duration)
	entry := ProgramEntry{
		ManagedRecord: ManagedRecord{
			Title: "overlap test entry",
			Hash:  buildHash([]string{channel.Hash, start.Format(time.RFC3339Nano), duration.String()}),
		},
		StartDateTime: &start,
		EndDateTime:   &end,
		ChannelID:     channel.ID,
	}
	if err := getDb().Create(&entry).Error; err != nil {
		t.Fatalf("could not create test program entry: %v", err)
	}
	return &entry
}

func getOverlapTestChannel(t *testing.T) Channel {
	t.Helper()
	channels := *mustGetChannels(t)
	if len(channels) == 0 {
		t.Fatalf("no channels present after setup")
	}
	return channels[0]
}

func TestHandleOverlapsByDayStoresAndClearsCollisions(t *testing.T) {
	setupInMemoryDbForTesting()
	// the raw overlap query requires the lowercase DbType
	appConf.DbType = "sqlite"
	db := getDb()

	channel := getOverlapTestChannel(t)
	day := time.Date(2026, 3, 2, 0, 0, 0, 0, time.UTC)

	first := createOverlapTestEntry(t, &channel, day.Add(20*time.Hour), time.Hour)
	second := createOverlapTestEntry(t, &channel, day.Add(20*time.Hour+30*time.Minute), time.Hour)
	unrelated := createOverlapTestEntry(t, &channel, day.Add(6*time.Hour), time.Hour)

	var wg sync.WaitGroup
	wg.Add(1)
	handleOverlapsByDay(&wg, &channel, day)
	wg.Wait()

	var deprecatedFirst ProgramEntry
	db.Preload("CollisionEntries").First(&deprecatedFirst, first.ID)
	if !deprecatedFirst.IsDeprecated {
		t.Fatalf("the overlapping entry must be deprecated")
	}
	if len(deprecatedFirst.CollisionEntries) != 1 || deprecatedFirst.CollisionEntries[0].ID != second.ID {
		t.Fatalf("expected one collision entry, got: %+v", deprecatedFirst.CollisionEntries)
	}

	var checkUnrelated ProgramEntry
	db.First(&checkUnrelated, unrelated.ID)
	if checkUnrelated.IsDeprecated {
		t.Fatalf("the non overlapping entry must not be deprecated")
	}

	// move the overlapping entry, so that the overlap disappears
	if err := db.Model(&ProgramEntry{}).Where("id = ?", second.ID).
		Update("start_date_time", day.Add(22*time.Hour+30*time.Minute)).Error; err != nil {
		t.Fatalf("could not move the overlapping entry: %v", err)
	}

	wg.Add(1)
	handleOverlapsByDay(&wg, &channel, day)
	wg.Wait()

	db.Preload("CollisionEntries").First(&deprecatedFirst, first.ID)
	if deprecatedFirst.IsDeprecated {
		t.Fatalf("the deprecation status of an entry without overlaps must be cleared")
	}
	if len(deprecatedFirst.CollisionEntries) != 0 {
		t.Fatalf("stale collision entries must be removed, got: %+v", deprecatedFirst.CollisionEntries)
	}

	var joinCount int64
	db.Table("collision_entries").Where("program_entry_id = ?", first.ID).Count(&joinCount)
	if joinCount != 0 {
		t.Fatalf("join table rows of removed collisions must be deleted, got %d", joinCount)
	}
}

func TestClearDeprecationsResetsCollisionCheck(t *testing.T) {
	setupInMemoryDbForTesting()
	db := getDb()

	channel := getOverlapTestChannel(t)

	now := time.Now().Truncate(time.Second)
	entry := createOverlapTestEntry(t, &channel, now, time.Hour)

	lastCheck := now.Add(-time.Hour)
	if err := db.Model(&ProgramEntry{}).Where("id = ?", entry.ID).
		Updates(map[string]any{"is_deprecated": true, "last_collision_check": &lastCheck}).Error; err != nil {
		t.Fatalf("could not deprecate the test entry: %v", err)
	}

	ClearDeprecations()

	var cleared ProgramEntry
	db.First(&cleared, entry.ID)
	if cleared.IsDeprecated {
		t.Fatalf("the deprecation status must be cleared")
	}
	if cleared.LastCollisionCheck != nil {
		t.Fatalf("the last collision check must be reset to null, was: %v", cleared.LastCollisionCheck)
	}
}
