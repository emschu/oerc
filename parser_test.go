// oerc, alias oer-collector
// Copyright (C) 2021-2025 emschu[aet]mailbox.org
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

func TestARDParserInterfaceCompliance(t *testing.T) {
	var parserInst interface{} = &ARDParser{}
	_, ok := parserInst.(ParserInterface)
	if !ok {
		t.Fatalf("Incompliant ARDParser instance detected!")
	}
}

func TestZDFParserInterfaceCompliance(t *testing.T) {
	var parserInst interface{} = &ZDFParser{}
	_, ok := parserInst.(ParserInterface)
	if !ok {
		t.Fatalf("Incompliant ZDFParser instance detected!")
	}
}

func TestORFParserInterfaceCompliance(t *testing.T) {
	var parserInst interface{} = &ORFParser{}
	_, ok := parserInst.(ParserInterface)
	if !ok {
		t.Fatalf("Incompliant ORFParser instance detected!")
	}
}

func TestSRFParserInterfaceCompliance(t *testing.T) {
	var parserInst interface{} = &SRFParser{}
	_, ok := parserInst.(ParserInterface)
	if !ok {
		t.Fatalf("Incompliant SRFParser instance detected!")
	}
}

func TestDaysInFuture(t *testing.T) {
	p := Parser{}
	now := time.Now()
	tomorrow := now.Add(25 * time.Hour)
	yesterday := now.Add(-25 * time.Hour)

	if !p.isMoreThanXDaysInFuture(&now, 0) {
		t.Fatalf("Invalid calculation!")
	}
	if !p.isMoreThanXDaysInPast(&now, 0) {
		t.Fatalf("Invalid calculation!")
	}
	if p.isMoreThanXDaysInFuture(&now, 1) {
		t.Fatalf("Invalid calculation!")
	}
	if p.isMoreThanXDaysInPast(&now, 1) {
		t.Fatalf("Invalid calculation!")
	}

	if !p.isMoreThanXDaysInFuture(&tomorrow, 1) {
		t.Fatalf("Invalid calculation!")
	}
	if !p.isMoreThanXDaysInPast(&yesterday, 1) {
		t.Fatalf("Invalid calculation!")
	}

	if p.isMoreThanXDaysInFuture(&tomorrow, 2) {
		t.Fatalf("Invalid calculation!")
	}
	if p.isMoreThanXDaysInPast(&yesterday, 2) {
		t.Fatalf("Invalid calculation!")
	}
}
