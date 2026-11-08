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
	"testing"
	"time"
)

func TestResolveOrfURL(t *testing.T) {
	testCases := []struct {
		in      string
		want    string
		wantErr bool
	}{
		{in: "/program/orf1/index-20260302~1200~2200-.html", want: "https://tv.orf.at/program/orf1/index-20260302~1200~2200-.html"},
		{in: "program/orf1/index.html", want: "https://tv.orf.at/program/orf1/index.html"},
		{in: "https://tv.orf.at/program/orf1", want: "https://tv.orf.at/program/orf1"},
		{in: "http://tv.orf.at/program/orf1", want: "http://tv.orf.at/program/orf1"},
		{in: "//tv.orf.at/program/orf1", want: "https://tv.orf.at/program/orf1"},
		{in: "", wantErr: true},
		{in: "   ", wantErr: true},
		{in: "://bad", wantErr: true},
	}
	for _, tc := range testCases {
		got, err := resolveOrfURL(tc.in)
		if tc.wantErr {
			if err == nil {
				t.Fatalf("expected an error for input '%s', got '%s'", tc.in, got)
			}
			continue
		}
		if err != nil {
			t.Fatalf("unexpected error for input '%s': %v", tc.in, err)
		}
		if got != tc.want {
			t.Fatalf("resolveOrfURL('%s') = '%s', want '%s'", tc.in, got, tc.want)
		}
	}
}

func TestGetDaysBetween(t *testing.T) {
	origConf := appConf
	defer func() { appConf = origConf }()

	parser := &ORFParser{}

	// simple cases in utc
	appConf.TimeZone = "UTC"
	now := time.Date(2026, 3, 10, 15, 0, 0, 0, time.UTC)
	if got := parser.getDaysBetween(now, now); got != 0 {
		t.Fatalf("expected 0 days for the same day, got %d", got)
	}
	if got := parser.getDaysBetween(now.Add(24*time.Hour), now); got != 1 {
		t.Fatalf("expected 1 day, got %d", got)
	}
	if got := parser.getDaysBetween(now.Add(-48*time.Hour), now); got != -2 {
		t.Fatalf("expected -2 days, got %d", got)
	}

	// the same calendar day in the configured time zone must not be counted
	// as two days, even if the system time zone is on the previous day
	appConf.TimeZone = "Europe/Berlin"
	systemNow := time.Date(2026, 3, 10, 23, 30, 0, 0, time.UTC)    // = 2026-03-11 00:30 Berlin
	configuredDay := time.Date(2026, 3, 11, 0, 30, 0, 0, time.UTC) // same instant, utc date components
	if got := parser.getDaysBetween(configuredDay, systemNow); got != 0 {
		t.Fatalf("expected 0 days for the same calendar day in the configured time zone, got %d", got)
	}

	// daylight saving change on 2026-03-29 in Europe/Berlin, the 29th has only 23 hours
	dayBefore := time.Date(2026, 3, 28, 12, 0, 0, 0, time.UTC)
	dayAfter := time.Date(2026, 3, 30, 12, 0, 0, 0, time.UTC)
	if got := parser.getDaysBetween(dayAfter, dayBefore); got != 2 {
		t.Fatalf("expected 2 calendar days across the daylight saving change, got %d", got)
	}
}
