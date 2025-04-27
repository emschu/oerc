// oerc, alias oer-collector
// Copyright (C) 2021-2024 emschu[aet]mailbox.org
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
	"log"
)

// setupPersistence: setup persistence during app's startup process. initializing global db object
func setupPersistence() {
	db := getDb()

	if db == nil {
		log.Fatalf("Error connecting to database!\n")
	}
	// timezone handling
	if GetAppConf().TimeZone != "" {
		db.Exec(fmt.Sprintf("SET TIME ZONE '%s'", GetAppConf().TimeZone))
		if GetAppConf().DbName != "" {
			db.Exec(fmt.Sprintf("ALTER DATABASE %s SET timezone TO '%s'", GetAppConf().DbName, GetAppConf().TimeZone))
		} else {
			log.Printf("Warning: No DbName is configured!\n")
		}
		// set search path to the selected schema
		db.Exec(fmt.Sprintf("SET search_path TO %s;", GetAppConf().DbSchema))
	} else {
		log.Printf("Warning: No TimeZone is configured!\n")
	}

	s, err := db.DB()
	if err != nil {
		log.Fatalf("Error connecting to database!\nError: %v\n", err)
	}
	pingErr := s.Ping()
	if pingErr != nil {
		log.Fatalf("Error connecting and pinging to database!\nError: %v\n", pingErr)
	}

	// check migrations of used models
	migrationErr := db.AutoMigrate(
		&ChannelFamily{},
		&Channel{},
		&TvShow{},
		&ProgramEntry{},
		&LogEntry{},
		&ImageLink{},
		&Settings{},
		&Recommendation{},
	)
	if migrationErr != nil {
		log.Printf("Problem during migration of database.\n")
		log.Fatal(migrationErr)
	}

	// ensure channels exist
	handleChannelFamiliesSetup()
	handleChannelsSetup()
}

// handleChannelFamiliesSetup: setup channel families in db, based on hard-coded channel family information
func handleChannelFamiliesSetup() {
	db := getDb()

	for _, f := range channelFamilyKeys {
		var channelFamily ChannelFamily
		db.Where("title like ?", f).First(&channelFamily)
		if channelFamily.ID <= 0 {
			log.Printf("Creating new channel family %s", f)
			db.Create(&ChannelFamily{
				Title: f,
			})
		}
	}
	if verboseGlobal {
		log.Println("All channel families are present")
	}
}

// handleChannelsSetup: setup channels in db, based on hard-coded channel information
func handleChannelsSetup() {
	db := getDb()

	var listOfUpdatedChannels []uint
	var channelFamilies []ChannelFamily
	db.Find(&channelFamilies)
	for _, channelFam := range channelFamilies {
		var channels *[]Channel
		switch channelFam.Title {
		case "ARD":
			channels = getArdChannels()
			break
		case "ZDF":
			channels = getZdfChannels()
			break
		case "SRF":
			channels = getSrfChannels()
			break
		case "ORF":
			channels = getOrfChannels()
			break
		default:
			log.Printf("Unknown channel family '%s'", channelFam.Title)
			return
		}
		var channelCounter = 0
		for _, c := range *channels {
			var channel Channel
			c.ChannelFamily = channelFam

			db.Where("hash = ? AND is_deprecated is false", c.Hash).First(&channel)
			if channel.ID == 0 {
				db.Create(&c)
				listOfUpdatedChannels = append(listOfUpdatedChannels, c.ID)
			} else {
				// update
				channel.Title = c.Title
				channel.URL = c.URL
				channel.TechnicalID = c.TechnicalID
				channel.Homepage = c.Homepage

				db.Save(&channel)
				listOfUpdatedChannels = append(listOfUpdatedChannels, channel.ID)
			}
			channelCounter++
		}
		if verboseGlobal {
			log.Printf("%s: %d channels present", channelFam.Title, channelCounter)
		}
	}

	// cleanup deprecation status of channels
	var notUpdatedChannels []Channel
	db.Model(&Channel{}).Where("is_deprecated = false AND id not IN(?)", listOfUpdatedChannels).Find(&notUpdatedChannels)

	for _, deprecatedChannel := range notUpdatedChannels {
		deprecatedChannel.IsDeprecated = true
		db.Save(&deprecatedChannel)
		log.Printf("Deprecating channel #%d with name '%s'\n", deprecatedChannel.ID, deprecatedChannel.Title)
	}
}
