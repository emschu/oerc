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

import "testing"

func TestArdPageCount(t *testing.T) {
	testCases := []struct {
		totalElements int
		pageSize      int
		want          int
	}{
		{0, 200, 0},
		{1, 200, 1},
		{200, 200, 1},
		{201, 200, 2},
		{400, 200, 2},
		{401, 200, 3},
		{100, 0, 0},
		{-5, 200, 0},
	}
	for _, tc := range testCases {
		if got := ardPageCount(tc.totalElements, tc.pageSize); got != tc.want {
			t.Errorf("ardPageCount(%d, %d) = %d, want %d", tc.totalElements, tc.pageSize, got, tc.want)
		}
	}
}

func TestFilterValidArdImageLinks(t *testing.T) {
	setupInMemoryDbForTesting()

	imageLinks := []ImageLink{
		{URL: "https://api.ardmediathek.de/image-service/image/xyz?width=600"},
		{URL: "https://programm-api.ard.de/images/abc.jpg"},
		{URL: "https://evil.example.com/payload.png"},
		{URL: "javascript:alert(1)"},
	}

	valid := filterValidArdImageLinks(imageLinks)
	if len(valid) != 2 {
		t.Fatalf("expected 2 valid image links, got %d", len(valid))
	}
	for _, link := range valid {
		if link.URL != "https://api.ardmediathek.de/image-service/image/xyz?width=600" &&
			link.URL != "https://programm-api.ard.de/images/abc.jpg" {
			t.Fatalf("unexpected image link survived the filter: %s", link.URL)
		}
	}

	if got := filterValidArdImageLinks(nil); len(got) != 0 {
		t.Fatalf("expected no valid image links for an empty input, got %d", len(got))
	}
}
