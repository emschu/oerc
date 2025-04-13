// oerc, alias oer-collector
// Copyright (C) 2021-2024 emschu[aet]mailbox.org
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
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"math"
	"net"
	"os"
	"path"
	"strings"
	"time"
)

// GetAppConf method to get the app's configuration
func GetAppConf() AppConfig {
	return appConf
}

// AppConfig struct to wrap the app's configuration
type AppConfig struct {
	Debug                        bool     `yaml:"Debug,omitempty"`
	ProfilingEnabled             bool     `yaml:"ProfilingEnabled,omitempty"`
	ForceUpdate                  bool     `yaml:"ForceUpdate"`
	TimeToRefreshInMinutes       int32    `yaml:"TimeToRefreshInMinutes"`
	DaysInPast                   uint     `yaml:"DaysInPast"`
	DaysInFuture                 uint     `yaml:"DaysInFuture"`
	EnableARD                    bool     `yaml:"EnableARD"`
	EnableZDF                    bool     `yaml:"EnableZDF"`
	EnableORF                    bool     `yaml:"EnableORF"`
	EnableSRF                    bool     `yaml:"EnableSRF"`
	EnableTVShowCollection       bool     `yaml:"EnableTVShowCollection"`
	EnableProgramEntryCollection bool     `yaml:"EnableProgramEntryCollection"`
	ProxyURL                     string   `yaml:"ProxyUrl"`
	ServerHost                   string   `yaml:"ServerHost"`
	ServerPort                   uint16   `yaml:"ServerPort"`
	ClientEnabled                bool     `yaml:"ClientEnabled"`
	TimeZone                     string   `yaml:"TimeZone"`
	DbType                       string   `yaml:"DbType"`
	DbHost                       string   `yaml:"DbHost"`
	DbPort                       uint32   `yaml:"DbPort"`
	DbName                       string   `yaml:"DbName"`
	DbSchema                     string   `yaml:"DbSchema"`
	DbUser                       string   `yaml:"DbUser"`
	DbPassword                   string   `yaml:"DbPassword"`
	DbSSLEnabled                 bool     `yaml:"DbSSLEnabled"`
	SearchKeywords               []string `yaml:"SearchKeywords"`
	SearchSkipChannels           []string `yaml:"SearchSkipChannels"`
	SearchDaysInFuture           uint     `yaml:"SearchDaysInFuture"`
	AccessControlAllowOrigin     string   `yaml:"AccessControlAllowOrigin,omitempty"`
}

func (a *AppConfig) verifyConfiguration() bool {
	// check time zone is valid
	_, err := time.LoadLocation(a.TimeZone)
	if err != nil {
		log.Printf("Invalid time zone '%s' given!\n", a.TimeZone)
		return false
	}
	// check db type is valid/supported
	dbType := a.DbType
	if len(dbType) == 0 {
		log.Printf("Invalid empty DbType given in configuration!\n")
		return false
	}
	if strings.ToLower(dbType) != "postgres" && strings.ToLower(dbType) != "postgresql" {
		log.Printf("Invalid DbType '%s' given!\n", a.DbType)
		return false
	}
	// check backend server configuration
	ip := net.ParseIP(a.ServerHost)
	if ip == nil {
		log.Printf("Invalid ServerHost provided in configuration!\n")
		return false
	}
	serverPort := a.ServerPort
	if serverPort == 0 || serverPort > uint16(math.Pow(2, 16)-2) {
		log.Printf("Invalid port number for server provided in configuration!\n")
		return false
	}
	if strings.TrimSpace(a.DbSchema) == "" {
		log.Printf("Invalid empty database schema provided in configuration!\n")
		return false
	}
	return true
}

// mechanism to detect the configuration file to use
func (a *AppConfig) loadConfiguration(inputPath string, allowFail bool) *string {
	// at first take the provided path - if possible
	var cleanedPath = path.Clean(inputPath)
	if len(cleanedPath) > 0 {
		providedFilePath, err := os.Stat(cleanedPath)
		if err != nil {
			if allowFail {
				log.Fatal(err)
			}
			return nil
		}
		if providedFilePath.Mode().IsRegular() {
			log.Printf("Loading configuration from file '%s'.\n", cleanedPath)
			loadYaml(cleanedPath)
			return &cleanedPath
		}
	}
	// then look in current directory for config.yaml
	homeDir, err := os.UserHomeDir()
	homeDir = path.Clean(homeDir)
	if err != nil {
		if allowFail {
			log.Fatalf("Home dir cannot be accessed - error: %v", err)
		}
		return nil
	}
	// then look in ~/.oerc.yaml
	homeDirConfigFile := fmt.Sprintf("%s/%s", homeDir, ".oerc.yaml")
	homeDirCfgFileStat, errHomeDir := os.Stat(homeDirConfigFile)
	if errHomeDir != nil {
		if allowFail {
			log.Fatalf("Could not find configuration file at '%s'", homeDirConfigFile)
		}
		return nil
	}
	if verboseGlobal {
		log.Printf("Loading configuration from file '%s'.\n", homeDirConfigFile)
	}
	if homeDirCfgFileStat.Mode().IsRegular() {
		loadYaml(homeDirConfigFile)
		return &homeDirConfigFile
	}
	if allowFail {
		log.Fatalf("Path '%s' is not a valid regular file.", homeDirConfigFile)
	}
	return nil
}

func loadYaml(path string) {
	f, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var appConfLoaded AppConfig
	yamlErr := yaml.UnmarshalStrict(f, &appConfLoaded)

	if yamlErr != nil {
		log.Fatalf("Problem with configuration file: %v", yamlErr)
	}
	appConf = appConfLoaded
}
