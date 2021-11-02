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

import (
	"encoding/json"
	"fmt"
	rice "github.com/GeertJohan/go.rice"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"sort"
	"strconv"
	"time"
)

var (
	version       = "0.9.12"
	appConf       AppConfig
	status        Status
	verboseGlobal = false
	// internal var to store a function which is - if not nil - executed right before the app stops, used by profiling feature only
	shutdownCb func()
)

// main entry point of oerc
func main() {
	app := &cli.App{
		Name:                 "oerc",
		Usage:                "Command line tool to manage the oerc application",
		Description:          "Fetch, view and search TV program data of public-law stations in Germany, Switzerland and Austria",
		EnableBashCompletion: false,
		HideHelp:             false,
		Version:              version + ", License: AGPLv3, https://github.com/emschu/oerc",
		Compiled:             time.Now(),
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:        "config",
				Usage:       "Path to the yaml configuration file",
				DefaultText: "~/.oerc.yaml",
				Aliases:     []string{"c"},
				TakesFile:   true,
			},
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: "Verbose log output",
				Value: false,
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "fetch",
				Aliases: []string{"f"},
				Usage:   "Get latest data",
				Action: func(c *cli.Context) error {
					log.Println("Starting fetch process")
					startTime := time.Now()
					Startup(c)
					defer Shutdown()

					isNetworkAvailable, err := connectivityCheck()
					if !isNetworkAvailable {
						return err
					}

					// parse the single channels
					if appConf.EnableARD {
						log.Printf("Parsing ARD Start\n")
						ParseARD()
						log.Printf("Parsing ARD End\n")
					}
					if appConf.EnableZDF {
						log.Printf("Parsing ZDF Start\n")
						ParseZDF()
						log.Printf("Parsing ZDF End\n")
					}
					if appConf.EnableORF {
						log.Printf("Parsing ORF Start\n")
						ParseORF()
						log.Printf("Parsing ORF End\n")
					}
					if appConf.EnableSRF {
						log.Printf("Parsing SRF Start\n")
						ParseSRF()
						log.Printf("Parsing SRF End\n")
					}

					FindOverlaps()

					// update counters
					if status.TotalCreatedPE > 0 || status.TotalCreatedTVS > 0 ||
						status.TotalUpdatedPE > 0 || status.TotalUpdatedTVS > 0 {
						setSetting(settingKeyLastFetch, time.Now().Format(time.RFC3339))
					}
					setSetting(settingKeyRequestsLastExecution, strconv.Itoa(int(status.TotalRequests)))
					currentRequestsTotal, err := strconv.ParseUint(getSetting(settingKeyRequestsTotal).Value, 10, 64)
					if err != nil {
						currentRequestsTotal = 0
					}
					setSetting(settingKeyRequestsTotal, strconv.Itoa(int(currentRequestsTotal+status.TotalRequests)))
					log.Printf("HTTP request counter of this fetch process: %d\n", status.TotalRequests)
					sub := time.Now().Sub(startTime)
					log.Printf("Duration: %.2f Seconds, %.2f Minutes\n", sub.Seconds(), sub.Minutes())
					log.Printf("Created program entries: %d. Updated: %d. Skipped: %d ", status.TotalCreatedPE, status.TotalUpdatedPE, status.TotalSkippedPE)
					log.Printf("Created tv shows: %d. Updated: %d.", status.TotalCreatedTVS, status.TotalUpdatedTVS)

					return nil
				},
			},
			{
				Name:    "status",
				Aliases: []string{"s"},
				Usage:   "show app's status",
				Action: func(c *cli.Context) error {
					Startup(c)
					defer Shutdown()

					object := getStatusObject()
					object.TvChannels = nil
					object.TvChannelFamilies = nil

					s, _ := json.MarshalIndent(object, "", "\t")
					log.Printf("%s\n", string(s))

					log.Printf("Total request counter: %s\n", getSetting(settingKeyRequestsTotal).Value)

					return nil
				},
			},
			{
				Name:    "server",
				Aliases: []string{"sv"},
				Usage:   "Start API HTTP backend server",
				Action: func(c *cli.Context) error {
					startTime := time.Now()
					Startup(c)
					defer Shutdown()
					log.Printf("Trying to start server at http://%s:%d ...\n", appConf.ServerHost, appConf.ServerPort)

					defer func() {
						log.Println("Server end")
					}()

					StartServer()
					duration := time.Now().Sub(startTime)
					log.Printf("Duration: %.2f Seconds, %.2f Minutes\n", duration.Seconds(), duration.Minutes())
					log.Printf("Created program entries: %d. Updated: %d. Skipped: %d ", status.TotalCreatedPE, status.TotalUpdatedPE, status.TotalSkippedPE)
					log.Printf("Created tv shows: %d. Updated: %d.", status.TotalCreatedTVS, status.TotalUpdatedTVS)

					return nil
				},
			},
			{
				Name:    "init",
				Aliases: []string{"i"},
				Usage:   "Initial database and configuration setup check",
				Action: func(c *cli.Context) error {
					initialStartup(c)
					log.Println("init done")
					return nil
				},
			},
			{
				Name:    "search",
				Aliases: []string{"sc"},
				Usage:   "Search program data and create recommendations",
				Action: func(c *cli.Context) error {
					startTime := time.Now()
					Startup(c)
					defer Shutdown()

					SearchProgram()

					setSetting(settingKeyLastSearch, fmt.Sprintf(time.Now().Format(time.RFC3339)))

					duration := time.Now().Sub(startTime)
					log.Printf("Duration: %.2f Seconds, %.2f Minutes\n", duration.Seconds(), duration.Minutes())
					return nil
				},
			},
			{
				Name:  "clear",
				Usage: "Clear the database. Be careful!",
				Subcommands: []*cli.Command{
					{
						Name:  "log",
						Usage: "Clears all logs",
						Action: func(c *cli.Context) error {
							log.Println("Clearing all log entries...")
							Startup(c)
							defer Shutdown()

							if c.Bool("force") {
								ClearLogs()
							} else {
								log.Printf("Please set the '--force true' flag to clear the logs from database.\n")
							}

							return nil
						},
					},
					{
						Name:  "recommendations",
						Usage: "Clears all recommendations",
						Action: func(c *cli.Context) error {
							log.Println("clear recommendations")
							Startup(c)
							defer Shutdown()

							if c.Bool("force") {
								ClearRecommendations()
							} else {
								log.Printf("Please set the '--force true' flag to clear ALL recommendations from database.\n")
							}

							return nil
						},
					},
					{
						Name:  "recommendations-old",
						Usage: "Clears old recommendations",
						Action: func(c *cli.Context) error {
							log.Println("clear recommendations-old")
							Startup(c)
							defer Shutdown()

							if c.Bool("force") {
								ClearOldRecommendations()
							} else {
								log.Printf("Please set the '--force true' flag to clear old recommendations from database.\n")
							}

							return nil
						},
					},
					{
						Name:  "overlaps",
						Usage: "Clearing overlap status of all program items",
						Action: func(c *cli.Context) error {
							log.Println("clear recommendations-old")
							Startup(c)
							defer Shutdown()

							if c.Bool("force") {
								ClearDeprecations()
							} else {
								log.Printf("Please set the '--force true' flag to clear program entries' overlap status from database.\n")
							}

							return nil
						},
					},
				},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "force",
						Aliases: []string{"f"},
					},
				},
				Action: func(c *cli.Context) error {
					log.Printf("Clear\n")
					Startup(c)
					defer Shutdown()

					if c.Bool("force") {
						ClearAll()
					} else {
						log.Printf("Please set the '--force true' flag to confirm cleaning the WHOLE database.\n")
					}

					return nil
				},
			},
			{
				Name:   "overlap-check-full",
				Hidden: true,
				Usage:  "Run overlap check on all program entries. Could take very long.",
				Action: func(context *cli.Context) error {
					Startup(context)
					defer Shutdown()

					FindOverlapsGlobal()

					return nil
				},
			},
		},
	}
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// Startup should be called on startup
func Startup(c *cli.Context) {
	appConf = AppConfig{}
	if c != nil {
		loadConfiguration(c.Path("config"), true)
	} else {
		log.Fatal("Problem with context")
	}

	isValid := verifyConfiguration()
	if !isValid {
		log.Fatalln("Invalid configuration! Startup cancelled.")
	}
	setupPersistence()
	verbose := c.Bool("verbose")
	if verbose {
		verboseGlobal = true
	}

	// profiling related stuff
	if isProfilingEnabled() {
		dateStr := time.Now().Format(time.RFC3339)

		f, err := ioutil.TempFile("", fmt.Sprintf("oerc-profiling-cpu-%s-*.pprof", dateStr))
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		log.Printf("Profiling (CPU) output is stored in %s\n", f.Name())

		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		shutdownCb = func() {
			defer f.Close()
			defer pprof.StopCPUProfile()

			memF, err := ioutil.TempFile("", fmt.Sprintf("oerc-profiling-mem-%s-*.pprof", dateStr))
			if err != nil {
				log.Fatal("could not create memory profile: ", err)
			}
			log.Printf("Profiling (Memory) output is stored in %s\n", memF.Name())

			runtime.GC() // get up-to-date statistics
			if err := pprof.WriteHeapProfile(memF); err != nil {
				log.Fatal("could not write memory profile: ", err)
			}
			defer memF.Close()
		}
	}
}

// Shutdown should be called on shutdown event
func Shutdown() {
	if httpClient != nil {
		defer httpClient.CloseIdleConnections()
	}
	if shutdownCb != nil {
		shutdownCb()
	}
	// close db
	db := getDb()
	s, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	err2 := s.Close()
	if err2 != nil {
		log.Fatal(err2)
	}
}

func initialStartup(c *cli.Context) {
	appConf = AppConfig{}
	if c != nil {
		configurationFile := loadConfiguration(c.Path("config"), false)
		if len(c.Path("config")) == 0 && configurationFile != nil {
			confBox := rice.MustFindBox("config")
			log.Printf("Trying to create default configuration at '%s'.\n", *configurationFile)
			err := ioutil.WriteFile(*configurationFile, confBox.MustBytes(".oerc_default.dist.yaml"), 0600)
			if err != nil {
				log.Fatalf("Error creating default configuration at '%s': %v.\n", *configurationFile, err)
			}
		}
	}
	isValid := verifyConfiguration()
	if !isValid {
		log.Fatalln("Invalid configuration! Please adjust and fix the configuration. Startup cancelled.")
	}
	setupPersistence()
}

func isDebug() bool {
	return appConf.Debug
}

func isProfilingEnabled() bool {
	return appConf.ProfilingEnabled
}

// Status struct to wrap app's status
type Status struct {
	TotalRequests   uint64
	TotalUpdatedPE  uint64
	TotalUpdatedTVS uint64
	TotalCreatedPE  uint64
	TotalCreatedTVS uint64
	TotalSkippedPE  uint64
}
