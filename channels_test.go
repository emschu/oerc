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

import "testing"

func TestGetChannels(t *testing.T) {
	ardChannels := getArdChannels()
	if len(*ardChannels) != 14 {
		t.Fatalf("Invalid amount of ard channels: %d", len(*ardChannels))
	}
	zdfChannels := getZdfChannels()
	if len(*zdfChannels) != 6 {
		t.Fatalf("Invalid amount of zdf channels: %d", len(*zdfChannels))
	}
	srfChannels := getSrfChannels()
	if len(*srfChannels) != 3 {
		t.Fatalf("Invalid amount of srf channels: %d", len(*srfChannels))
	}
	orfChannels := getOrfChannels()
	if len(*orfChannels) != 4 {
		t.Fatalf("Invalid amount of orf channels: %d", len(*orfChannels))
	}
}

func TestGetChannelFamily(t *testing.T) {
	setupInMemoryDbForTesting()
	handleChannelsSetup()
	if getDb() == nil {
		t.Fatalf("nil db")
	}
	families := []string{
		"ARD", "ZDF", "ORF", "SRF",
	}

	for _, family := range families {
		dbFamily := getChannelFamily(getDb(), family)
		if dbFamily == nil {
			t.Fatalf("channel family cannot be found")
		}
	}
}

func TestGetChannelsOfFamily(t *testing.T) {
	setupInMemoryDbForTesting()
	handleChannelsSetup()
	if getDb() == nil {
		t.Fatalf("nil db")
	}
	ardFamily := getChannelFamily(getDb(), "ARD")
	zdfFamily := getChannelFamily(getDb(), "ZDF")
	orfFamily := getChannelFamily(getDb(), "ORF")
	srfFamily := getChannelFamily(getDb(), "SRF")

	ardChannels := getChannelsOfFamily(getDb(), ardFamily)
	zdfChannels := getChannelsOfFamily(getDb(), zdfFamily)
	orfChannels := getChannelsOfFamily(getDb(), orfFamily)
	srfChannels := getChannelsOfFamily(getDb(), srfFamily)

	if len(ardChannels) != len(*getArdChannels()) {
		t.Fatalf("invalid number of ard channels")
	}
	if len(zdfChannels) != len(*getZdfChannels()) {
		t.Fatalf("invalid number of zdf channels")
	}
	if len(orfChannels) != len(*getOrfChannels()) {
		t.Fatalf("invalid number of orf channels")
	}
	if len(srfChannels) != len(*getSrfChannels()) {
		t.Fatalf("invalid number of srf channels")
	}
}

func TestChannelsOrdering(t *testing.T) {
	setupInMemoryDbForTesting()
	handleChannelsSetup()

	channels := *getChannels()
	if len(channels) < 2 {
		t.Skip("Not enough channels to test ordering")
	}

	// Reverse the channels
	reversedChannels := make([]Channel, len(channels))
	for i, j := 0, len(channels)-1; i < len(channels); i, j = i+1, j-1 {
		reversedChannels[i] = channels[j]
	}

	// Update priority based on reversed order
	db := getDb()
	for i, channel := range reversedChannels {
		db.Model(&Channel{}).Where("id = ?", channel.ID).Update("priority", i)
	}

	// Fetch channels again and check order
	orderedChannels := *getChannels()
	for i, channel := range orderedChannels {
		if channel.ID != reversedChannels[i].ID {
			t.Errorf("Channel at index %d has wrong ID. Expected %d, got %d", i, reversedChannels[i].ID, channel.ID)
		}
		if channel.Priority != i {
			t.Errorf("Channel at index %d has wrong priority. Expected %d, got %d", i, i, channel.Priority)
		}
	}
}
