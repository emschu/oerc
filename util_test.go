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

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateDateRange(t *testing.T) {
	if len(*generateDateRangeInPastAndFuture(0, 0)) != 1 {
		t.Error("invalid date range generated")
	}
	if len(*generateDateRangeInPastAndFuture(1, 0)) != 2 {
		t.Error("invalid date range generated")
	}
	if len(*generateDateRangeInPastAndFuture(0, 1)) != 2 {
		t.Error("invalid date range generated")
	}
	if len(*generateDateRangeInPastAndFuture(1, 1)) != 3 {
		t.Error("invalid date range generated")
	}
	if len(*generateDateRangeInPastAndFuture(10, 10)) != 21 {
		t.Error("invalid date range generated")
	}
}

func TestTrim(t *testing.T) {
	in := "\n\t\t\t\t\t\t\tZDX-Morgenmagazin\t\t\t\t\t\tn\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\tModeration: X Y, Z D und F A\t\t\t\t\t\t\t\n\t\t\t\t\t\n\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t"
	trimString := trimAndSanitizeString(in)
	if strings.Contains(trimString, "\t") {
		t.Error("Invalid tab found!")
	}
	if strings.Contains(trimString, "\n") {
		t.Error("Invalid newline found!")
	}
}

func TestIcal(t *testing.T) {
	content, err := parseStartAndEndDateTimeFromIcal("https://programm.ard.de/ICalendar/iCal---Sendung?sendung=281063652013560")
	if err != nil {
		t.Errorf("Error during fetch of ical content '%s'", err)
	}
	if content == nil {
		t.Error("Invalid ical content response")
	}
}

func TestTrimAndSanitizeString(t *testing.T) {
	// they should all evaluate to "test"
	var inputStrings = []string{" test", "test ", "test", "   test", "test   ", "  test  "}
	for _, v := range inputStrings {
		if "test" != trimAndSanitizeString(v) {
			t.Error("Invalid trim and sanitizing function")
		}
	}
	var maliciousStrings = []string{"test<script></script>", "test<script>", "<script>alert('hello')</script>test"}
	for _, v := range maliciousStrings {
		if "test" != trimAndSanitizeString(v) {
			t.Errorf("Invalid return value of trim and sanitizing function. input: '%s', output: '%s'", v, trimAndSanitizeString(v))
		}
	}
	if "<b>test</b>" != trimAndSanitizeString("<b>test</b> ") {
		t.Errorf("Invalid return value of trim and sanitizing function.")
	}
}

func TestGetProxy(t *testing.T) {
	// invalid proxy settings -> fallback to http.ProxyFromEnvironment expected
	appConf.ProxyURL = ""
	if getHTTPProxy() == nil {
		t.Errorf("proxy should NOT be nil")
	}
	appConf.ProxyURL = " "
	if getHTTPProxy() == nil {
		t.Errorf("proxy should NOT be nil")
	}
	appConf.ProxyURL = "test"
	if getHTTPProxy() == nil {
		t.Error("proxy should NOT be nil")
	}
	// missing port
	appConf.ProxyURL = "http://localhost"
	if getHTTPProxy() == nil {
		t.Error("proxy should NOT be nil")
	}
	// valid example, happy case
	appConf.ProxyURL = "http://localhost:7676"
	if getHTTPProxy() == nil {
		t.Error("proxy should NOT be nil")
	}
	// invalid port number
	appConf.ProxyURL = "http://localhost:767676"
	if getHTTPProxy() == nil {
		t.Error("proxy should NOT be nil")
	}
}

func TestGetBaseCollector(t *testing.T) {
	host := []string{"example.com"}
	collector := baseCollector(host)
	if collector == nil {
		t.Error("Collector should not be nil")
	}
	if len(collector.AllowedDomains) == 0 || collector.AllowedDomains[0] != host[0] {
		t.Error("Allowed collector does not allow specified domain")
	}
}

func TestConnectivity(t *testing.T) {
	check, err := connectivityCheck()
	if err != nil || !check {
		t.Error("Basic connection test failed")
	}
}

func TestErrorRegistering(t *testing.T) {
	resetErr()
	if errorCounter != 0 {
		t.Error("Error counter is expected to be 0!")
	}
	incrErr()
	if errorCounter != 1 {
		t.Error("Error counter is expected to be 1")
	}
	resetErr()
	checkErr() // should do nothing, because 1 < errorThreshold
	if errorCounter != 0 {
		t.Error("Error counter is expected to be 0!")
	}
}

func TestHash(t *testing.T) {
	got := buildHash([]string{"3", "2"})
	if got == "" {
		t.Errorf("empty hash! %s", got)
	}
	if len(got) != 32 {
		t.Errorf("invalid length != 32 of hash")
	}
}

func TestAppLog(t *testing.T) {
	setupInMemoryDbForTesting()
	appLog("Test example")

	db := getDb()
	var entry LogEntry
	db.Model(&LogEntry{}).Last(&entry)
	if entry.Message != "Test example" {
		t.Fatalf("Simple app log test failed")
	}
}

func TestIsRecentlyUpdated(t *testing.T) {
	if (&ProgramEntry{LastCheck: nil}).isRecentlyUpdated() || (&ProgramEntry{LastCheck: &time.Time{}}).isRecentlyUpdated() {
		t.Fatalf("cannot be updated if last check is nil")
	}
	const i = 15
	appConf.TimeToRefreshInMinutes = i
	appConf.ForceUpdate = false
	fakeLastCheckTime := time.Now().Add(-i*time.Minute - 1*time.Second)
	fakeLastCheckTime2 := time.Now().Add(-i * time.Minute)
	fakeLastCheckTime3 := time.Now().Add(-i*time.Minute + 1*time.Second)
	if (&ProgramEntry{LastCheck: &fakeLastCheckTime}).isRecentlyUpdated() ||
		(&ProgramEntry{LastCheck: &fakeLastCheckTime2}).isRecentlyUpdated() ||
		!(&ProgramEntry{LastCheck: &fakeLastCheckTime3}).isRecentlyUpdated() {
		t.Fatalf("time range check failed")
	}
	appConf.ForceUpdate = true
	if (&ProgramEntry{LastCheck: &fakeLastCheckTime}).isRecentlyUpdated() ||
		(&ProgramEntry{LastCheck: &fakeLastCheckTime2}).isRecentlyUpdated() ||
		(&ProgramEntry{LastCheck: &fakeLastCheckTime3}).isRecentlyUpdated() {
		t.Fatalf("force-update check failed")
	}
}

func TestIsRecentlyFetched(t *testing.T) {
	setupInMemoryDbForTesting()
	appConf.ForceUpdate = true
	if isRecentlyFetched() {
		t.Fatalf("Expected recent fetch is forced")
	}
	appConf.ForceUpdate = false
	if isRecentlyFetched() {
		t.Fatalf("Expected recent fetch did not take place")
	}
	setSetting(settingKeyLastFetch, time.Now().Format(time.RFC3339))
	if !isRecentlyFetched() {
		t.Fatalf("Expected recent fetch took place")
	}
}

func TestClearOldRecommendations(t *testing.T) {
	setupInMemoryDbForTesting()

	db := getDb()
	oldRec := time.Now().Add(-1 * time.Minute)
	newRec := time.Now().Add(20 * time.Minute)
	db.Create(&Recommendation{ProgramEntryID: 123, ChannelID: 4, StartDateTime: &oldRec})
	db.Create(&Recommendation{ProgramEntryID: 123, ChannelID: 4, StartDateTime: &newRec})

	var counter int64
	db.Model(&Recommendation{}).Count(&counter)
	if counter != 2 {
		t.Fatalf("Test logic fail")
	}
	ClearOldRecommendations()
	db.Model(&Recommendation{}).Count(&counter)
	if counter != 1 {
		t.Fatalf("One entry should be deleted")
	}
	ClearRecommendations()
	db.Model(&Recommendation{}).Count(&counter)
	if counter != 0 {
		t.Fatalf("There should be zero entries after cleanup")
	}
}

func TestConsiderTagExists(t *testing.T) {
	setupInMemoryDbForTesting()
	verboseGlobal = true

	var pe ProgramEntry
	if pe.Tags != "" {
		t.Fatalf("There should be no tags in a new program entry!")
	}

	emptyStr := ""
	testTag := "test"
	pe.considerTagExists(&emptyStr)
	if pe.Tags != "" {
		t.Fatalf("Empty string is not a tag")
	}
	pe.considerTagExists(&testTag)
	if pe.Tags != "test" {
		t.Fatalf("There should be a new tag 'test'")
	}
	pe.considerTagExists(&testTag)
	if pe.Tags != "test" {
		t.Fatalf("There should be a new tag 'test'")
	}
	testTag2 := "test2"
	pe.considerTagExists(&testTag2)
	if pe.Tags != "test;test2" {
		t.Fatalf("There should be a new tag 'test2'")
	}
}

func TestChunkStringSlice(t *testing.T) {
	slice1 := []string{}
	if len(chunkStringSlice(slice1, 0)) != 0 {
		t.Fatalf("Invalid handling of empty slice of chunk size 0")
	}
	if len(chunkStringSlice(slice1, 1)) != 0 {
		t.Fatalf("Invalid handling of empty slice of chunk size 1")
	}
	slice2 := []string{"", ""}
	if len(chunkStringSlice(slice2, 0)) != 0 {
		t.Fatalf("Invalid handling of slice(2) of chunk size 0")
	}
	if len(chunkStringSlice(slice2, 1)) != 2 {
		t.Fatalf("Invalid handling of slice with two elements")
	}
	if len(chunkStringSlice(slice2, 2)) != 1 {
		t.Fatalf("Invalid handling of slice with two elements")
	}
	slice3 := []string{"", "", ""}
	if len(chunkStringSlice(slice3, 0)) != 0 {
		t.Fatalf("Invalid handling of slice(3) of chunk size 0")
	}
	if len(chunkStringSlice(slice3, 1)) != 3 {
		t.Fatalf("Invalid handling of slice with three elements")
	}
	if len(chunkStringSlice(slice3, 2)) != 2 {
		t.Fatalf("Invalid handling of slice with three elements")
	}
}
