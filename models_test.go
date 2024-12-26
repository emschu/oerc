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
	"testing"
	"time"
)

func TestProgramEntryTagConsidering(t *testing.T) {
	verboseGlobal = true
	setupInMemoryDbForTesting()

	now := time.Now()
	dateInFuture := now.Add(1 * time.Hour)
	p := ProgramEntry{
		StartDateTime: &now,
		EndDateTime:   &dateInFuture,
	}
	if len(p.Tags) != 0 {
		t.Fatalf("Tags should be empty")
	}
	testTag := "testTag"
	p.considerTagExists(&testTag)
	p.saveProgramEntryRecord(getDb())
	if len(p.Tags) != len(testTag) {
		t.Fatalf("Tags should contain one tag '%s'", testTag)
	}

	p.considerTagExists(&testTag)
	if len(p.Tags) != len(testTag) {
		t.Fatalf("Tags should contain one tag after second add '%s'", testTag)
	}
}

func TestProgramEntryImageLinkConsidering(t *testing.T) {
	verboseGlobal = true
	setupInMemoryDbForTesting()

	now := time.Now()
	dateInFuture := now.Add(1 * time.Hour)
	p := ProgramEntry{
		StartDateTime: &now,
		EndDateTime:   &dateInFuture,
	}
	if len(p.ImageLinks) != 0 {
		t.Fatalf("Tags should be empty")
	}
	testImageLink := "testImageLink"
	p.considerImageLinkExists(testImageLink)
	p.saveProgramEntryRecord(getDb())
	if len(p.ImageLinks) != 1 {
		t.Fatalf("Tags should contain one tag '%s'", testImageLink)
	}

	p.considerImageLinkExists(testImageLink)
	if len(p.ImageLinks) != 1 {
		t.Fatalf("Tags should contain one tag after second add '%s'", testImageLink)
	}
}
