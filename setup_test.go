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
	"os"
	"path/filepath"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupInMemoryDbForTesting() {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	dBReference = db
	setupPersistence()
}

// setupIsolatedSqliteForTesting creates a temporary isolated sqlite database for
// testing, the database and the configuration are restored on test cleanup
func setupIsolatedSqliteForTesting(t *testing.T) {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "oerc-isolated-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})

	origConf := appConf
	t.Cleanup(func() {
		if sqlDB, err := getDb().DB(); err == nil {
			sqlDB.Close()
		}
		dBReference = nil
		appConf = origConf
	})

	appConf = *defaultAppConfig()
	appConf.DbType = "sqlite"
	appConf.DbHost = filepath.Join(tmpDir, "isolated.sqlite")
	appConf.TimeZone = "UTC"
	dBReference = nil
	setupPersistence()
}

// mustGetChannels fetches the channels and fails the test on database errors
func mustGetChannels(t *testing.T) *[]Channel {
	t.Helper()
	channels, err := getChannels()
	if err != nil {
		t.Fatalf("error fetching channels: %v", err)
	}
	return channels
}

func TestSetupPersistence(t *testing.T) {
	verboseGlobal = true
	setupInMemoryDbForTesting()

	// test channels are set up
	channels, err := getChannels()
	if err != nil {
		t.Fatalf("error fetching channels: %v", err)
	}
	if len(*channels) != 27 {
		t.Fatalf("Invalid number of channels after db setup: %d.\n", len(*channels))
	}
	// test channel families are set up
	channelFamilies, err := getChannelFamilies()
	if err != nil {
		t.Fatalf("error fetching channel families: %v", err)
	}
	if len(*channelFamilies) != 4 {
		t.Fatalf("Invalid number of channel families after db setup: %d", len(*channelFamilies))
	}
}
