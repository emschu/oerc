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
	"strings"
	"testing"
	"time"
)

func TestXMLTvDocumentHasTvRootElement(t *testing.T) {
	start := time.Date(2026, 3, 2, 20, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)

	channelFamily := ChannelFamily{ID: 1, Title: "ARD"}
	channel := Channel{
		BaseModel: BaseModel{ID: 1},
		ManagedRecord: ManagedRecord{
			Title:         "Das Erste",
			ChannelFamily: channelFamily,
		},
	}
	entry := ProgramEntry{
		BaseModel: BaseModel{ID: 42},
		ManagedRecord: ManagedRecord{
			Title:           "Test show",
			URL:             "https://programm-api.ard.de/program/api/teaser?teaserId=1",
			ChannelFamilyID: channelFamily.ID,
		},
		StartDateTime: &start,
		EndDateTime:   &end,
		Channel:       channel,
	}

	doc, err := newXMLTvDocument(&[]Channel{channel}, &[]ProgramEntry{entry})
	if err != nil {
		t.Fatalf("unexpected error creating the xmltv document: %v", err)
	}

	output, err := writeXMLTVToString(doc)
	if err != nil {
		t.Fatalf("unexpected error writing the xmltv document: %v", err)
	}

	trimmed := strings.TrimSpace(output)
	if !strings.HasPrefix(trimmed, "<tv") {
		t.Fatalf("XMLTV output must start with a <tv> root element, was: %.80s", trimmed)
	}
	if !strings.HasSuffix(trimmed, "</tv>") {
		t.Fatalf("XMLTV output must end with </tv>, was: %.80s", trimmed)
	}
	if strings.Contains(output, "Document") {
		t.Fatalf("XMLTV output must not contain the go struct name 'Document'")
	}
	if !strings.Contains(output, `<channel id="DE: Das Erste"`) {
		t.Fatalf("XMLTV output must contain the channel element, was: %s", trimmed)
	}
	if !strings.Contains(output, `<programme channel="DE: Das Erste"`) {
		t.Fatalf("XMLTV output must contain the programme element, was: %s", trimmed)
	}
	if !strings.Contains(output, "Test show") {
		t.Fatalf("XMLTV output must contain the programme title, was: %s", trimmed)
	}
}
