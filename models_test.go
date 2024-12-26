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
