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
	"testing"
	"time"
)

func TestRefreshMaterializedView(t *testing.T) {
	setupIsolatedSqliteForTesting(t)
	db := getDb()

	setupMaterializedView()

	// for sqlite a regular view is used, the refresh must be a safe no-op
	refreshMaterializedView()

	// the view must reflect the current data
	channel := getOverlapTestChannel(t)
	now := time.Now()
	entry := createOverlapTestEntry(t, &channel, now, time.Hour)
	refreshMaterializedView()

	var statusInfo StatusInfoModel
	if err := db.Table(statusInfo.TableName()).First(&statusInfo).Error; err != nil {
		t.Fatalf("the status view must be queryable after the refresh: %v", err)
	}
	if statusInfo.ProgramEntryCount != 1 {
		t.Fatalf("expected the view to reflect 1 program entry, got %d", statusInfo.ProgramEntryCount)
	}
	if entry.ID == 0 {
		t.Fatalf("test setup error, program entry was not created")
	}
}
