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
	"os"
	"path/filepath"
	"testing"
)

func TestSqliteIntegration(t *testing.T) {
	// Setup a temporary directory for the sqlite database
	tmpDir, err := os.MkdirTemp("", "oerc-sqlite-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test_oerc.sqlite")

	// Configure app to use sqlite
	appConf.DbType = "sqlite"
	appConf.DbHost = dbPath
	appConf.TimeZone = "UTC"
	// Reset dBReference to force re-initialization in getDb()
	dBReference = nil

	// Initialize database
	setupPersistence()

	// Verify that the database file was created
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Fatalf("SQLite database file was not created at %s", dbPath)
	}

	db := getDb()
	if db == nil {
		t.Fatal("Failed to get database connection")
	}

	// Verify that tables were created and populated
	var familyCount int64
	db.Model(&ChannelFamily{}).Count(&familyCount)
	if familyCount != 4 {
		t.Errorf("Expected 4 channel families, got %d", familyCount)
	}

	var channelCount int64
	db.Model(&Channel{}).Count(&channelCount)
	if channelCount == 0 {
		t.Error("Expected channels to be populated, but found 0")
	}

	// Test basic CRUD operation
	testFamily := ChannelFamily{Title: "TestFamily"}
	if err := db.Create(&testFamily).Error; err != nil {
		t.Fatalf("Failed to create record in SQLite: %v", err)
	}

	var foundFamily ChannelFamily
	if err := db.First(&foundFamily, "title = ?", "TestFamily").Error; err != nil {
		t.Fatalf("Failed to retrieve record from SQLite: %v", err)
	}

	if foundFamily.Title != "TestFamily" {
		t.Errorf("Expected Title 'TestFamily', got '%s'", foundFamily.Title)
	}

	// Clean up for other tests if they run in the same process
	dBReference = nil
}

func TestSqliteReconnect(t *testing.T) {
	// Setup a temporary directory for the sqlite database
	tmpDir, err := os.MkdirTemp("", "oerc-sqlite-reconnect-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "reconnect_oerc.sqlite")

	// 1. Initial setup
	appConf.DbType = "sqlite"
	appConf.DbHost = dbPath
	appConf.TimeZone = "UTC"
	dBReference = nil

	setupPersistence()

	db := getDb()
	testFamily := ChannelFamily{Title: "PersistentFamily"}
	db.Create(&testFamily)

	// 2. "Close" and reconnect
	dBReference = nil

	// Re-initialize (mocking app restart)
	// We want to avoid the "Unknown channel family" log by not having unexpected families in the DB
	// but for this test it's fine as long as it works.
	setupPersistence()
	db = getDb()

	var foundFamily ChannelFamily
	if err := db.First(&foundFamily, "title = ?", "PersistentFamily").Error; err != nil {
		t.Fatalf("Failed to retrieve record after reconnect: %v", err)
	}

	if foundFamily.Title != "PersistentFamily" {
		t.Errorf("Expected Title 'PersistentFamily', got '%s'", foundFamily.Title)
	}

	dBReference = nil
}
