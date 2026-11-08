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
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

// TestDbErrorsAreReturnedInsteadOfFatal ensures the data access functions used by
// request handlers return errors instead of terminating the whole server process
func TestDbErrorsAreReturnedInsteadOfFatal(t *testing.T) {
	setupIsolatedSqliteForTesting(t)
	sqlDB, err := getDb().DB()
	if err != nil {
		t.Fatalf("could not get the sql db: %v", err)
	}
	sqlDB.Close()

	now := time.Now()
	later := now.Add(time.Hour)

	if _, err := getProgramOfWeb(&now, &later, nil); err == nil {
		t.Fatalf("expected an error from getProgramOfWeb for a closed database connection")
	}
	if _, err := getChannels(); err == nil {
		t.Fatalf("expected an error from getChannels for a closed database connection")
	}
	if _, err := getChannelFamilies(); err == nil {
		t.Fatalf("expected an error from getChannelFamilies for a closed database connection")
	}
	if _, err := getProgramOf(&now, &later, nil); err == nil {
		t.Fatalf("expected an error from getProgramOf for a closed database connection")
	}
}

func TestGetLogEntriesHandlerPagination(t *testing.T) {
	setupIsolatedSqliteForTesting(t)
	db := getDb()

	for i := range 251 {
		if err := db.Create(&LogEntry{Message: fmt.Sprintf("entry %d", i)}).Error; err != nil {
			t.Fatalf("could not create log entry: %v", err)
		}
	}

	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v2/log?page=0", nil)
	getLogEntriesHandler(c)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
	var firstPage LogEntriesResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &firstPage); err != nil {
		t.Fatalf("could not decode the response: %v", err)
	}
	if firstPage.PageCount != 2 {
		t.Fatalf("expected page_count 2 for 251 entries, got %d", firstPage.PageCount)
	}
	if firstPage.PageSize != 250 {
		t.Fatalf("expected page_size 250, got %d", firstPage.PageSize)
	}
	if firstPage.Size != 250 {
		t.Fatalf("expected 250 elements on the first page, got %d", firstPage.Size)
	}
	if firstPage.EntryCount != 251 {
		t.Fatalf("expected entry_count 251, got %d", firstPage.EntryCount)
	}

	recorder2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(recorder2)
	c2.Request = httptest.NewRequest(http.MethodGet, "/api/v2/log?page=1", nil)
	getLogEntriesHandler(c2)

	var lastPage LogEntriesResponse
	if err := json.Unmarshal(recorder2.Body.Bytes(), &lastPage); err != nil {
		t.Fatalf("could not decode the response: %v", err)
	}
	if lastPage.Size != 1 {
		t.Fatalf("expected 1 element on the last page, got %d", lastPage.Size)
	}
}

func TestGetXMLTvHandlerRangeLimits(t *testing.T) {
	setupIsolatedSqliteForTesting(t)
	gin.SetMode(gin.TestMode)

	// a date range of more than 14 days must be rejected
	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v2/xmltv?from=2026-01-01T00:00:00Z&to=2026-03-01T00:00:00Z", nil)
	getXMLTvHandler(c)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for a too large date range, got %d", recorder.Code)
	}

	// a small range must be answered with a valid xmltv document
	recorder2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(recorder2)
	c2.Request = httptest.NewRequest(http.MethodGet, "/api/v2/xmltv?from=2026-01-01T00:00:00Z&to=2026-01-02T00:00:00Z", nil)
	getXMLTvHandler(c2)
	if recorder2.Code != http.StatusOK {
		t.Fatalf("expected 200 for a small date range, got %d", recorder2.Code)
	}
	if !strings.Contains(recorder2.Body.String(), "<tv") {
		t.Fatalf("expected a <tv> root element in the response, got: %.80s", recorder2.Body.String())
	}
}

func TestGetRecommendationsAtPagination(t *testing.T) {
	setupIsolatedSqliteForTesting(t)
	db := getDb()

	channel := getOverlapTestChannel(t)
	now := time.Now().Truncate(time.Second)

	for i := range 5 {
		start := now.Add(time.Duration(i) * time.Hour)
		end := start.Add(time.Hour)
		entry := ProgramEntry{
			ManagedRecord: ManagedRecord{
				Title: fmt.Sprintf("recommendation entry %d", i),
				Hash:  buildHash([]string{"recommendation-test", fmt.Sprintf("%d", i)}),
			},
			StartDateTime: &start,
			EndDateTime:   &end,
			ChannelID:     channel.ID,
		}
		if err := db.Create(&entry).Error; err != nil {
			t.Fatalf("could not create program entry: %v", err)
		}
		if err := db.Create(&Recommendation{
			ProgramEntryID: entry.ID,
			ChannelID:      channel.ID,
			StartDateTime:  &start,
			Keywords:       "test",
		}).Error; err != nil {
			t.Fatalf("could not create recommendation: %v", err)
		}
	}

	all, err := getRecommendationsAt(now, 100, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(all) != 5 {
		t.Fatalf("expected 5 recommendations, got %d", len(all))
	}

	firstPage, err := getRecommendationsAt(now, 2, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(firstPage) != 2 {
		t.Fatalf("expected 2 recommendations on the first page, got %d", len(firstPage))
	}
	secondPage, err := getRecommendationsAt(now, 2, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(secondPage) != 2 {
		t.Fatalf("expected 2 recommendations on the second page, got %d", len(secondPage))
	}
	if firstPage[0].ID == secondPage[0].ID {
		t.Fatalf("the offset must shift the window of returned recommendations")
	}
	lastPage, err := getRecommendationsAt(now, 2, 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lastPage) != 1 {
		t.Fatalf("expected 1 recommendation on the last page, got %d", len(lastPage))
	}
}

func TestPutChannelsHandlerLimitsBodySize(t *testing.T) {
	setupIsolatedSqliteForTesting(t)
	gin.SetMode(gin.TestMode)

	// a valid but oversized request body must be rejected
	hugeTitle := strings.Repeat("a", maxChannelRequestBodyBytes+1024)
	payload := fmt.Sprintf(`[{"id":1,"priority":5,"title":"%s"}]`, hugeTitle)

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = httptest.NewRequest(http.MethodPut, "/api/v2/channels", strings.NewReader(payload))
	c.Request.Header.Set("Content-Type", "application/json")
	putChannelsHandler(c)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for an oversized request body, got %d", recorder.Code)
	}

	// a normal sized request must still work
	recorder2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(recorder2)
	c2.Request = httptest.NewRequest(http.MethodPut, "/api/v2/channels", strings.NewReader(`[{"id":1,"priority":5}]`))
	c2.Request.Header.Set("Content-Type", "application/json")
	putChannelsHandler(c2)
	if recorder2.Code != http.StatusOK {
		t.Fatalf("expected 200 for a normal sized request body, got %d", recorder2.Code)
	}
}
