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
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
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

func TestGenerateDateRangeBetweenDates(t *testing.T) {
	today := time.Now()
	dates := generateDateRangeBetweenDates(today, today)
	if len(*dates) != 1 {
		t.Errorf("one day expected, was: %d", len(*dates))
	}
	tomorrow := today.Add(24 * time.Hour)
	dates2 := generateDateRangeBetweenDates(today, tomorrow)
	if len(*dates2) != 2 {
		t.Errorf("two days expected, was: %d", len(*dates2))
	}
	dates3 := generateDateRangeBetweenDates(tomorrow, today)
	if len(*dates3) != 2 {
		t.Errorf("two days expected, was: %d", len(*dates3))
	}
}

func TestGenerateDateRangeBetweenDatesDstAndBounds(t *testing.T) {
	location, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		t.Fatalf("could not load location: %v", err)
	}

	// the day of the daylight saving change (2026-03-29) has only 23 hours,
	// calendar days must be used instead of 24 hour blocks
	start := time.Date(2026, 3, 28, 0, 0, 0, 0, location)
	end := time.Date(2026, 3, 30, 0, 0, 0, 0, location)
	dates := generateDateRangeBetweenDates(start, end)
	if len(*dates) != 3 {
		t.Fatalf("expected 3 calendar days, got %d: %v", len(*dates), *dates)
	}

	// a single day range must return the requested day, not time.Now()
	sameDate := generateDateRangeBetweenDates(start, start)
	if len(*sameDate) != 1 {
		t.Fatalf("expected one element for a single day range, got %d", len(*sameDate))
	}
	if !(*sameDate)[0].Equal(start) {
		t.Fatalf("expected the requested start date, got %v", (*sameDate)[0])
	}

	// reversed ranges must work as well
	reversed := generateDateRangeBetweenDates(end, start)
	if len(*reversed) != 3 {
		t.Fatalf("expected 3 calendar days for the reversed range, got %d", len(*reversed))
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
	// restore the configuration, an invalid proxy would leak into subsequent tests
	origConf := appConf
	defer func() { appConf = origConf }()

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
	errorCount := atomic.LoadUint64(&globalErrorCounter)
	if errorCount != 0 {
		t.Error("Error counter is expected to be 0!")
	}
	incrErr()
	errorCount = atomic.LoadUint64(&globalErrorCounter)
	if errorCount != 1 {
		t.Error("Error counter is expected to be 1")
	}
	resetErr()
	checkErr() // should do nothing, because 1 < errorThreshold
	errorCount = atomic.LoadUint64(&globalErrorCounter)
	if errorCount != 0 {
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
	db.Create(&Recommendation{ProgramEntryID: 456, ChannelID: 4, StartDateTime: &newRec})

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

func TestParseDate(t *testing.T) {
	setupInMemoryDbForTesting()

	location, _ := time.LoadLocation(defaultAppConfig().TimeZone)

	// rfc3339, for example used in orf page attributes
	date, fail := parseDate("2022-02-22T10:10:11Z", location)
	if fail || date.IsZero() {
		t.Fatalf("invalid date")
	}

	// database formats as returned by the status info view, they must be
	// parseable for the global overlap check to work
	dbDates := []string{
		"2026-08-26 21:00:00+02:00", // sqlite
		"2026-08-26 21:00:00+02",    // postgres
		"2026-08-26 21:00:00",       // without time zone information
		"2026-08-26 21:00:00.123456+02:00",
	}
	for _, dbDate := range dbDates {
		date, fail = parseDate(dbDate, location)
		if fail || date.IsZero() {
			t.Fatalf("database date '%s' must be parseable", dbDate)
		}
		if date.Month() != time.August || date.Day() != 26 || date.Hour() != 21 {
			t.Fatalf("unexpected parsed date '%s': %v", dbDate, date)
		}
	}

	invalidDates := []string{
		"",
		"1",
		"2022-02-22T10:10:11",
		"not-a-date",
	}
	for _, invalidDate := range invalidDates {
		date, fail2 := parseDate(invalidDate, location)
		if !fail2 || !date.IsZero() {
			t.Fatalf("invalid date parsing result for '%s'", invalidDate)
		}
	}
}

func TestResolvePath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get user home dir: %v", err)
	}

	tests := []struct {
		input    string
		expected string
	}{
		{"~/test.db", filepath.Join(home, "test.db")},
		{"/tmp/test.db", "/tmp/test.db"},
		{"test.db", "test.db"},
		{"~", home},
	}

	for _, tc := range tests {
		actual := resolvePath(tc.input)
		if actual != tc.expected {
			t.Errorf("resolvePath(%q) = %q; want %q", tc.input, actual, tc.expected)
		}
	}
}

func TestGetDbWithHomeDir(t *testing.T) {
	// Backup original config
	origConf := appConf
	defer func() { appConf = origConf }()

	// Set up config with ~ path
	appConf = *defaultAppConfig()
	appConf.DbType = "sqlite"
	appConf.DbHost = "~/oerc_test_home.sqlite"

	// Reset dBReference to force new connection
	dBReference = nil

	db := getDb()
	if db == nil {
		t.Fatal("Failed to get DB with ~ path")
	}

	// Check if file was created in home dir
	home, _ := os.UserHomeDir()
	expectedPath := filepath.Join(home, "oerc_test_home.sqlite")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("Expected database file to be created at %s, but it doesn't exist", expectedPath)
	}

	// Cleanup
	sqlDB, _ := db.DB()
	sqlDB.Close()
	os.Remove(expectedPath)
}

func TestSaveProgramEntryRecordUpdatesImageLinks(t *testing.T) {
	setupIsolatedSqliteForTesting(t)

	db := getDb()
	channels := *mustGetChannels(t)
	if len(channels) == 0 {
		t.Fatalf("no channels present after setup")
	}
	channel := channels[0]

	now := time.Now().Truncate(time.Second)
	later := now.Add(time.Hour)
	entry := ProgramEntry{
		ManagedRecord: ManagedRecord{
			Title: "image link test",
			Hash:  buildHash([]string{"image-link-test", now.Format(time.RFC3339)}),
		},
		StartDateTime: &now,
		EndDateTime:   &later,
		ChannelID:     channel.ID,
	}
	entry.considerImageLinkExists("https://example.com/first.png")
	entry.saveProgramEntryRecord(db)
	if entry.ID == 0 {
		t.Fatalf("program entry was not created")
	}

	// reload with preloaded image links and add a second one
	var stored ProgramEntry
	db.Preload("ImageLinks").First(&stored, entry.ID)
	stored.considerImageLinkExists("https://example.com/second.png")
	stored.saveProgramEntryRecord(db)

	var linkCount int64
	db.Model(&ImageLink{}).Where("program_entry_id = ?", entry.ID).Count(&linkCount)
	if linkCount != 2 {
		t.Fatalf("expected 2 image links after the update, got %d", linkCount)
	}

	// removing image links must be persisted as well
	var reloaded ProgramEntry
	db.Preload("ImageLinks").First(&reloaded, entry.ID)
	reloaded.ImageLinks = reloaded.ImageLinks[:1]
	reloaded.saveProgramEntryRecord(db)

	db.Model(&ImageLink{}).Where("program_entry_id = ?", entry.ID).Count(&linkCount)
	if linkCount != 1 {
		t.Fatalf("expected 1 image link after the removal, got %d", linkCount)
	}

	// no orphaned image links must remain
	var totalCount int64
	db.Model(&ImageLink{}).Count(&totalCount)
	if totalCount != 1 {
		t.Fatalf("removed image links must be deleted and not orphaned, got %d rows in total", totalCount)
	}
}

func TestHTTPClientProxyIsEvaluatedPerRequest(t *testing.T) {
	transport, ok := httpClient.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("httpClient must use an http.Transport")
	}
	if transport.Proxy == nil {
		t.Fatalf("the transport must have a proxy function")
	}

	req, err := http.NewRequest("GET", "https://example.com", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	origConf := appConf
	defer func() { appConf = origConf }()

	// the proxy has to be resolved from the current configuration per request,
	// not from the empty configuration state at package initialization
	appConf.ProxyURL = "http://localhost:7676"
	proxyURL, err := transport.Proxy(req)
	if err != nil {
		t.Fatalf("unexpected proxy error: %v", err)
	}
	if proxyURL == nil || proxyURL.String() != "http://localhost:7676" {
		t.Fatalf("the configured proxy must be returned, got: %v", proxyURL)
	}

	// without configuration the environment fallback must be used
	appConf.ProxyURL = ""
	if _, err := transport.Proxy(req); err != nil {
		t.Fatalf("unexpected proxy error: %v", err)
	}
}
