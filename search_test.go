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
