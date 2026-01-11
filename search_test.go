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
)

func TestUniqueStringSlice(t *testing.T) {
	var testSlice []string
	if len(uniqueStrSlice(&testSlice)) != 0 {
		t.Fatalf("Expected empty slice as return value")
	}
	testSlice2 := []string{"1", "2"}
	if len(uniqueStrSlice(&testSlice2)) != 2 {
		t.Fatalf("Expected slice with 2 elements")
	}
	testSlice3 := []string{"1", "1"}
	if len(uniqueStrSlice(&testSlice3)) != 1 {
		t.Fatalf("Expected slice with 1 element")
	}
}
