//
// oerc, alias oer-collector
// Copyright (C) 2021 emschu[aet]mailbox.org
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
				Title:       "ARD – Das Erste",
				URL:         "",
				Hash:        "28106",
				TechnicalID: "28106",
				Homepage:    "https://ard.de",
			},
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "BR Fernsehen",
				URL:         "",
				Hash:        "28107",
				TechnicalID: "28107",
				Homepage:    "https://www.br.de/fernsehen/index.html",
			},
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "HR Fernsehen",
				URL:         "",
				Hash:        "28108",
				TechnicalID: "28108",
				Homepage:    "https://www.hr-fernsehen.de/",
			},
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "MDR Fernsehen",
				URL:         "",
				Hash:        "28229",
				TechnicalID: "28229",
				Homepage:    "https://www.mdr.de/tv/",
			},
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "NDR Fernsehen",
				URL:         "",
				Hash:        "28226",
				TechnicalID: "28226",
				Homepage:    "https://www.ndr.de",
			},
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "RBB Fernsehen",
				URL:         "",
				Hash:        "28205",
				TechnicalID: "28205",
				Homepage:    "https://www.rbb-online.de/fernsehen/",
			},
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "Radio Bremen TV",
				URL:         "",
				Hash:        "28385",
				TechnicalID: "28385",
				Homepage:    "https://www.radiobremen.de/fernsehen/",
			},
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "SR Fernsehen",
				URL:         "",
				Hash:        "28486",
				TechnicalID: "28486",
				Homepage:    "https://www.sr.de/sr/home/fernsehen/",
			},
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "SWR BW Fernsehen",
				URL:         "",
				Hash:        "28113",
				TechnicalID: "28113",
				Homepage:    "https://www.swrfernsehen.de/tv-programm/",
			},
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "SWR RP Fernsehen",
				URL:         "",
				Hash:        "28231",
				TechnicalID: "28231",
				Homepage:    "https://www.swrfernsehen.de/tv-programm/",
			},
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "WDR Fernsehen",
				URL:         "",
				Hash:        "28111",
				TechnicalID: "28111",
				Homepage:    "http://www.wdr.de/tv/",
			},
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "ARD ALPHA",
				URL:         "",
				Hash:        "28487",
				TechnicalID: "28487",
				Homepage:    "http://www.br.de/fernsehen/ard-alpha/",
			},
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "Tagesschau24",
				URL:         "",
				Hash:        "28721",
				TechnicalID: "28721",
				Homepage:    "http://programm.tagesschau24.de/",
			},
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "ARD One",
				URL:         "",
				Hash:        "28722",
				TechnicalID: "28722",
				Homepage:    "http://www.one.ard.de/",
			},
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "KIKA",
				URL:         "",
				Hash:        "28008",
				TechnicalID: "28008",
				Homepage:    "http://www.kika.de/",
			},
		},
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
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "ZDFinfo",
				URL:         "",
				Hash:        "zdfinfo",
				TechnicalID: "zdfinfo",
				Homepage:    "https://www.zdf.de/dokumentation/zdfinfo-doku",
			},
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "ZDFneo",
				URL:         "",
				Hash:        "zdfneo",
				TechnicalID: "zdfneo",
				Homepage:    "https://www.zdf.de/sender/zdfneo",
			},
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "Phoenix",
				URL:         "",
				Hash:        "phoenix",
				TechnicalID: "phoenix",
				Homepage:    "http://www.phoenix.de/",
			},
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "3Sat",
				URL:         "",
				Hash:        "3sat",
				TechnicalID: "3sat",
				Homepage:    "https://www.3sat.de",
			},
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "ARTE",
				URL:         "",
				Hash:        "arte",
				TechnicalID: "arte",
				Homepage:    "https://www.arte.tv/de",
			},
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
				TechnicalID: "srf-1",
				Homepage:    "https://www.srf.ch/tv",
			},
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "SRF zwei",
				URL:         "",
				Hash:        "srf-2",
				TechnicalID: "srf-2",
				Homepage:    "https://www.srf.ch/tv/srf-2",
			},
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "SRF info",
				URL:         "",
				Hash:        "srf-info",
				TechnicalID: "srf-info",
				Homepage:    "http://www.srf.ch/tv",
			},
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
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "ORF 2",
				URL:         "",
				Hash:        "orf2",
				TechnicalID: "orf2",
				Homepage:    "http://tv.orf.at/program/orf2",
			},
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "ORF III",
				URL:         "",
				Hash:        "orf3",
				TechnicalID: "orf3",
				Homepage:    "https://tv.orf.at/orf3",
			},
		},
		{
			ManagedRecord: ManagedRecord{
				Title:       "ORF Sport +",
				URL:         "",
				Hash:        "orfs",
				TechnicalID: "orfs",
				Homepage:    "http://sport.orf.at/",
			},
		},
	}
}

// method to get a channel family record by it's id, e.g. "ARD", "ZDF" ...
func getChannelFamily(db *gorm.DB, channelFamilyID string) *ChannelFamily {
	var channelFamily ChannelFamily
	db.Where("title = ?", channelFamilyID).First(&channelFamily)
	return &channelFamily
}

// method to get all channels of a given ChannelFamily
func getChannelsOfFamily(db *gorm.DB, channelFamilyRecord *ChannelFamily) []Channel {
	var channels []Channel
	db.Where("channel_family_id = ?", &channelFamilyRecord.ID).Find(&channels)
	return channels
}
