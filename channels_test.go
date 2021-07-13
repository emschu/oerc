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
