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
	if len(*channels) != 28 {
		t.Fatalf("Invalid number of channels after db setup: %d.\n", len(*channels))
	}
	// test channel families are set up
	channelFamilies := getChannelFamilies()
	if len(*channelFamilies) != 4 {
		t.Fatalf("Invalid number of channel families after db setup: %d", len(*channelFamilies))
	}
}
