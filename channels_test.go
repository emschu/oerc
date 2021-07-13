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

import "testing"

func TestGetChannels(t *testing.T) {
	ardChannels := getArdChannels()
	if len(*ardChannels) != 15 {
		t.Fatalf("Invalid amount of ard channels: %d", len(*ardChannels))
	}
	zdfChannels := getZdfChannels()
	if len(*zdfChannels) != 6 {
		t.Fatalf("Invalid amount of zdf channels: %d", len(*zdfChannels))
	}
	srfChannels := getSrfChannels()
	if len(*srfChannels) != 3 {
		t.Fatalf("Invalid amount of srf channels: %d", len(*srfChannels))
	}
	orfChannels := getOrfChannels()
	if len(*orfChannels) != 5 {
		t.Fatalf("Invalid amount of orf channels: %d", len(*orfChannels))
	}
}
