package main

import "testing"

func TestAppConfig(t *testing.T) {
	config := &AppConfig{}
	if config.verifyConfiguration() {
		t.Fatalf("Empty config not allowed to be valid!")
	}
	// make the config minimally valid
	config.TimeZone = "Europe/Berlin"
	if config.verifyConfiguration() {
		t.Fatalf("Config not yet allowed to be valid!")
	}
	config.DbType = "postgres"
	if config.verifyConfiguration() {
		t.Fatalf("Config not yet allowed to be valid!")
	}
	config.ServerHost = "127.0.0.1"
	if config.verifyConfiguration() {
		t.Fatalf("Config not yet allowed to be valid!")
	}
	config.ServerPort = 8000
	if !config.verifyConfiguration() {
		t.Fatalf("Config should be valid!")
	}
}

func TestLoadConfiguration(t *testing.T) {
	config := AppConfig{}
	demoConfigFile := "config/.oerc.dist.yaml"
	configuration := config.loadConfiguration(demoConfigFile, true)
	if *configuration != demoConfigFile {
		t.Fatalf("Unexpected config file! %v", *configuration)
	}

	failedConfig := config.loadConfiguration("sdfsdf", false)
	if failedConfig != nil {
		t.Fatalf("Wrong path should return nil")
	}
}

func defaultAppConfig() *AppConfig {
	return &AppConfig{
		Debug:      true,
		TimeZone:   "Europe/Berlin",
		DbType:     "postgres",
		ServerHost: "127.0.0.1",
		ServerPort: 8080,
	}
}
