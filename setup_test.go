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
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func setupInMemoryDbForTesting() {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	dBReference = db
	setupPersistence()
}

func TestSetupPersistence(t *testing.T) {
	verboseGlobal = true
	setupInMemoryDbForTesting()

	// test channels are set up
	channels := getChannels()
	if len(*channels) != 27 {
		t.Fatalf("Invalid number of channels after db setup: %d.\n", len(*channels))
	}
	// test channel families are set up
	channelFamilies := getChannelFamilies()
	if len(*channelFamilies) != 4 {
		t.Fatalf("Invalid number of channel families after db setup: %d", len(*channelFamilies))
	}
}
