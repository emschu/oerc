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

import "gorm.io/gorm"

// this file contains all available channels we handle in this piece of software
// NOTE: the channelFamily <---> Channel association is in some cases just a technical aspect

func getArdChannels() *[]Channel {
	return &[]Channel{
		{
			ManagedRecord: ManagedRecord{
				Title:       "Das Erste",
				URL:         "",
				Hash:        "daserste",
				TechnicalID: "daserste",
				Homepage:    "https://ard.de",
			},
			IsDeprecated: false,
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "BR Fernsehen",
				URL:         "",
				Hash:        "br",
				TechnicalID: "br",
				Homepage:    "https://www.br.de",
			},
			IsDeprecated: false,
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "HR Fernsehen",
				URL:         "",
				Hash:        "hr",
				TechnicalID: "hr",
				Homepage:    "https://www.hr-fernsehen.de/",
			},
			IsDeprecated: false,
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "MDR Fernsehen",
				URL:         "",
				Hash:        "mdr",
				TechnicalID: "mdr",
				Homepage:    "https://www.mdr.de/tv/",
			},
			IsDeprecated: false,
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "NDR Fernsehen",
				URL:         "",
				Hash:        "ndr",
				TechnicalID: "ndr",
				Homepage:    "https://www.ndr.de",
			},
			IsDeprecated: false,
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "RBB Fernsehen",
				URL:         "",
				Hash:        "rbb",
				TechnicalID: "rbb",
				Homepage:    "https://www.rbb-online.de/fernsehen/",
			},
			IsDeprecated: false,
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "Radio Bremen TV",
				URL:         "",
				Hash:        "radiobremen",
				TechnicalID: "radiobremen",
				Homepage:    "https://www.radiobremen.de/fernsehen/",
			},
			IsDeprecated: false,
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "SR Fernsehen",
				URL:         "",
				Hash:        "sr",
				TechnicalID: "sr",
				Homepage:    "https://www.sr.de/sr/home/fernsehen/",
			},
			IsDeprecated: false,
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "SWR Fernsehen",
				URL:         "",
				Hash:        "swr",
				TechnicalID: "swr",
				Homepage:    "https://www.swr.de",
			},
			IsDeprecated: false,
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "WDR Fernsehen",
				URL:         "",
				Hash:        "wdr",
				TechnicalID: "wdr",
				Homepage:    "http://www.wdr.de/tv/",
			},
			IsDeprecated: false,
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "ARD ALPHA",
				URL:         "",
				Hash:        "alpha",
				TechnicalID: "alpha",
				Homepage:    "http://www.br.de/fernsehen/ard-alpha/",
			},
			IsDeprecated: false,
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "Tagesschau24",
				URL:         "",
				Hash:        "tagesschau24",
				TechnicalID: "tagesschau24",
				Homepage:    "http://programm.tagesschau24.de/",
			},
			IsDeprecated: false,
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "ARD One",
				URL:         "",
				Hash:        "one",
				TechnicalID: "one",
				Homepage:    "http://www.one.ard.de/",
			},
			IsDeprecated: false,
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "KIKA",
				URL:         "",
				Hash:        "kika",
				TechnicalID: "kika",
				Homepage:    "http://www.kika.de/",
			},
			IsDeprecated: false,
		},
		// also available: [
		//  "3sat 3sat",
		//  "arte ARTE",
		//  "kika KiKA",
		//  "phoenix phoenix",
		// ]

	}
}

func getZdfChannels() *[]Channel {
	return &[]Channel{
		{
			ManagedRecord: ManagedRecord{
				Title:       "ZDF",
				URL:         "",
				Hash:        "zdf",
				TechnicalID: "zdf",
				Homepage:    "http://www.zdf.de/",
			},
			IsDeprecated: false,
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "ZDFinfo",
				URL:         "",
				Hash:        "zdfinfo",
				TechnicalID: "zdfinfo",
				Homepage:    "https://www.zdf.de/dokumentation/zdfinfo-doku",
			},
			IsDeprecated: false,
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "ZDFneo",
				URL:         "",
				Hash:        "zdfneo",
				TechnicalID: "zdfneo",
				Homepage:    "https://www.zdf.de/sender/zdfneo",
			},
			IsDeprecated: false,
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "Phoenix",
				URL:         "",
				Hash:        "phoenix",
				TechnicalID: "phoenix",
				Homepage:    "http://www.phoenix.de/",
			},
			IsDeprecated: false,
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "3Sat",
				URL:         "",
				Hash:        "3sat",
				TechnicalID: "3sat",
				Homepage:    "https://www.3sat.de",
			},
			IsDeprecated: false,
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "ARTE",
				URL:         "",
				Hash:        "arte",
				TechnicalID: "arte",
				Homepage:    "https://www.arte.tv/de",
			},
			IsDeprecated: false,
		},
	}
}

func getSrfChannels() *[]Channel {
	return &[]Channel{
		{
			ManagedRecord: ManagedRecord{
				Title:       "SRF 1",
				URL:         "",
				Hash:        "srf-1",
				TechnicalID: "23FFBE1B-65CE-4188-ADD2-C724186C2C9F",
				Homepage:    "https://www.srf.ch/tv",
			},
			IsDeprecated: false,
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "SRF zwei",
				URL:         "",
				Hash:        "srf-zwei",
				TechnicalID: "E4D5AD08-C1E8-46A3-BB58-4875051D60D2",
				Homepage:    "https://www.srf.ch/tv/srf-2",
			},
			IsDeprecated: false,
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "SRF info",
				URL:         "",
				Hash:        "srf-info",
				TechnicalID: "34c2819e-e715-43d7-9026-40a443152a97",
				Homepage:    "https://www.srf.ch/tv",
			},
			IsDeprecated: false,
		},
	}
}

func getOrfChannels() *[]Channel {
	return &[]Channel{
		{
			ManagedRecord: ManagedRecord{
				Title:       "ORF eins",
				URL:         "",
				Hash:        "orf1",
				TechnicalID: "orf1",
				Homepage:    "http://tv.orf.at/program/orf1",
			},
			IsDeprecated: false,
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "ORF 2",
				URL:         "",
				Hash:        "orf2",
				TechnicalID: "orf2",
				Homepage:    "http://tv.orf.at/program/orf2",
			},
			IsDeprecated: false,
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "ORF III",
				URL:         "",
				Hash:        "orf3",
				TechnicalID: "orf3",
				Homepage:    "https://tv.orf.at/orf3",
			},
			IsDeprecated: false,
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "ORF Sport +",
				URL:         "",
				Hash:        "orfs",
				TechnicalID: "orfs",
				Homepage:    "http://sport.orf.at/",
			},
			IsDeprecated: false,
		},
	}
}

// method to get a channel family record by id, e.g. "ARD", "ZDF" ...
func getChannelFamily(db *gorm.DB, channelFamilyID string) *ChannelFamily {
	var channelFamily ChannelFamily
	db.Where("title = ?", channelFamilyID).First(&channelFamily)
	return &channelFamily
}

// method to get all channels of a given ChannelFamily
func getChannelsOfFamily(db *gorm.DB, channelFamilyRecord *ChannelFamily) []Channel {
	var channels []Channel
	db.Where("channel_family_id = ? and is_deprecated = false", &channelFamilyRecord.ID).Find(&channels)
	return channels
}
