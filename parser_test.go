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
