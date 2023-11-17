// oerc, alias oer-collector
// Copyright (C) 2021-2023 emschu[aet]mailbox.org
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
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	_ "net/http/pprof"

	"github.com/pkg/profile"
)

var (
	version       = "0.14.0"
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

					FetchAll(newDefaultDateRangeHandler())

					FindOverlaps(newDefaultDateRangeHandlerPadded(1))

					// update counters if something has happened
					showFetchResults(startTime, nil, nil)

					return nil
				},
			},
			{
				Name:    "fetch-range",
				Aliases: []string{"fr"},
				Usage:   "Fetch a specific date range",
				Flags: []cli.Flag{
					&cli.TimestampFlag{
						Name:        "from",
						Usage:       "Data range start date, YYYY-MM-DD",
						Layout:      "2006-01-02",
						DefaultText: "empty",
					},
					&cli.TimestampFlag{
						Name:        "to",
						Usage:       "Data range end date, YYYY-MM-DD",
						Layout:      "2006-01-02",
						DefaultText: "empty",
					},
				},
				Action: func(c *cli.Context) error {
					log.Println("Starting fetch-range process")
					startTime := time.Now()
					Startup(c)
					defer Shutdown()

					isNetworkAvailable, err := connectivityCheck()
					if !isNetworkAvailable {
						return err
					}

					rangeStart := c.Timestamp("from")
					rangeEnd := c.Timestamp("to")

					if rangeStart == nil || rangeEnd == nil {
						return fmt.Errorf("invalid date range start or end given")
					}

					FetchAll(newSpecificDateRangeHandler(*rangeStart, *rangeEnd))

					FindOverlaps(newSpecificDateRangeHandlerPadded(*rangeStart, *rangeEnd, 1))

					// update counters if something has happened
					showFetchResults(startTime, rangeStart, rangeEnd)

					return nil
				},
			},
			{
				Name:    "status",
				Aliases: []string{"s"},
				Usage:   "Show some of the app's status information",
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
				Usage:   "Start webserver with oerc API and an embedded browser client",
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
					log.Printf("Starting search process...")
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
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "force",
								Aliases: []string{"f"},
								Value:   false,
							},
						},
						Action: func(c *cli.Context) error {
							log.Println("Clearing all log entries...")
							Startup(c)
							defer Shutdown()

							if c.Bool("force") {
								ClearLogs()
							} else {
								log.Printf("Please set the '--force=true' flag to clear the logs from database.\n")
							}

							return nil
						},
					},
					{
						Name:  "recommendations",
						Usage: "Clears all recommendations",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "force",
								Aliases: []string{"f"},
								Value:   false,
							},
						},
						Action: func(c *cli.Context) error {
							log.Println("Clear recommendations")
							Startup(c)
							defer Shutdown()

							if c.Bool("force") {
								ClearRecommendations()
							} else {
								log.Printf("Please set the '--force=true' flag to clear ALL recommendations from database.\n")
							}

							return nil
						},
					},
					{
						Name:  "recommendations-old",
						Usage: "Clears old recommendations",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "force",
								Aliases: []string{"f"},
								Value:   false,
							},
						},
						Action: func(c *cli.Context) error {
							log.Println("Clear old recommendations")
							Startup(c)
							defer Shutdown()

							if c.Bool("force") {
								ClearOldRecommendations()
							} else {
								log.Printf("Please set the '--force=true' flag to clear old recommendations from database.\n")
							}

							return nil
						},
					},
					{
						Name:  "overlaps",
						Usage: "Clearing overlap status of all program items",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "force",
								Aliases: []string{"f"},
								Value:   false,
							},
						},
						Action: func(c *cli.Context) error {
							log.Println("clear overlaps")
							Startup(c)
							defer Shutdown()

							if c.Bool("force") {
								ClearDeprecations()
							} else {
								log.Printf("Please set the '--force=true' flag to clear program entries' overlap status from database.\n")
							}

							return nil
						},
					},
				},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "force",
						Aliases: []string{"f"},
						Value:   false,
					},
				},
				Action: func(c *cli.Context) error {
					log.Printf("Clearing all data...\n")
					Startup(c)
					defer Shutdown()

					if c.Bool("force") {
						log.Printf("This could take a while...\n")
						ClearAll()
					} else {
						log.Printf("Please set the '--force=true' flag to confirm cleaning the WHOLE database.\n")
					}

					return nil
				},
			},
			{
				Name:  "full-overlap-check",
				Usage: "Run overlap check on all program entries. Could take very long.",
				Action: func(context *cli.Context) error {
					startTime := time.Now()
					Startup(context)
					defer Shutdown()

					FindOverlapsGlobal()

					duration := time.Now().Sub(startTime)
					log.Printf("Duration: %.2f Seconds, %.2f Minutes\n", duration.Seconds(), duration.Minutes())
					return nil
				},
			},
			{
				Name:  "overlap-check",
				Usage: "Run overlap check on currently fetched time range",
				Action: func(context *cli.Context) error {
					startTime := time.Now()
					Startup(context)
					defer Shutdown()

					FindOverlaps(newDefaultDateRangeHandlerPadded(1))

					duration := time.Now().Sub(startTime)
					log.Printf("Duration: %.2f Seconds, %.2f Minutes\n", duration.Seconds(), duration.Minutes())

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

func showFetchResults(startTime time.Time, start *time.Time, end *time.Time) {
	if status.TotalCreatedPE > 0 || status.TotalCreatedTVS > 0 ||
		status.TotalUpdatedPE > 0 || status.TotalUpdatedTVS > 0 {
		setSetting(settingKeyLastFetch, time.Now().Format(time.RFC3339))
	}
	setSetting(settingKeyRequestsLastExecution, strconv.Itoa(int(status.TotalRequests)))
	setting := getSetting(settingKeyRequestsTotal)
	var err error
	var currentRequestsTotal uint64
	if setting != nil {
		currentRequestsTotal, err = strconv.ParseUint(setting.Value, 10, 64)
		if err != nil {
			currentRequestsTotal = 0
		}
	}
	setSetting(settingKeyRequestsTotal, strconv.Itoa(int(currentRequestsTotal+status.TotalRequests)))
	if start != nil && end != nil {
		// if they are nil, the default fetch time range has been used
		dateLayout := "2006-01-02"
		log.Printf("Fetched date range: %s to %s", start.Format(dateLayout), end.Format(dateLayout))
	}
	log.Printf("HTTP request counter of this fetch process: %d\n", status.TotalRequests)
	sub := time.Now().Sub(startTime)
	log.Printf("Duration: %.2f Seconds, %.2f Minutes\n", sub.Seconds(), sub.Minutes())
	log.Printf("Created program entries: %d. Updated: %d. Skipped: %d ", status.TotalCreatedPE, status.TotalUpdatedPE, status.TotalSkippedPE)
	log.Printf("Created tv shows: %d. Updated: %d.", status.TotalCreatedTVS, status.TotalUpdatedTVS)
}

// FetchAll central method to fetch the data based on a dateRangeHandler
func FetchAll(handler dateRangeHandler) {
	// parse the single channels
	if appConf.EnableARD {
		log.Printf("Parsing ARD Start\n")
		var parser = &ARDParser{Parser: Parser{
			ChannelFamilyKey: "ARD",
			dateRangeHandler: handler,
		}}
		parser.Fetch(parser)
		log.Printf("Parsing ARD End\n")
	}
	if appConf.EnableZDF {
		log.Printf("Parsing ZDF Start\n")
		var parser = &ZDFParser{Parser: Parser{
			ChannelFamilyKey: "ZDF",
			dateRangeHandler: handler,
		}}
		parser.Fetch(parser)
		log.Printf("Parsing ZDF End\n")
	}
	if appConf.EnableORF {
		log.Printf("Parsing ORF Start\n")
		var parser = &ORFParser{Parser: Parser{
			ChannelFamilyKey: "ORF",
			dateRangeHandler: handler,
		}}
		parser.Fetch(parser)
		log.Printf("Parsing ORF End\n")
	}
	if appConf.EnableSRF {
		log.Printf("Parsing SRF Start\n")
		var parser = &SRFParser{Parser: Parser{
			ChannelFamilyKey: "SRF",
			dateRangeHandler: handler,
		}}
		parser.Fetch(parser)
		log.Printf("Parsing SRF End\n")
	}
}

// Startup should be called on startup
func Startup(c *cli.Context) {
	appConf = AppConfig{}
	if c != nil {
		appConf.loadConfiguration(c.Path("config"), true)
	} else {
		log.Fatal("Problem with context")
	}

	isValid := appConf.verifyConfiguration()
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
		go func() {
			err := http.ListenAndServe("127.0.0.1:9999", nil)
			if err != nil {
				log.Fatalf("Could not start profiling endpoint: '%v'\n", err)
			}
			log.Printf("Profiling endpoint started at 127.0.0.1:9999")
		}()

		shutdownCb = func() {
			defer profile.Start(profile.MemProfile).Stop()
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
		configurationFile := appConf.loadConfiguration(c.Path("config"), false)
		if len(c.Path("config")) == 0 && configurationFile != nil {
			confBox := rice.MustFindBox("config")
			log.Printf("Trying to create default configuration at '%s'.\n", *configurationFile)
			err := ioutil.WriteFile(*configurationFile, confBox.MustBytes(".oerc_default.dist.yaml"), 0600)
			if err != nil {
				log.Fatalf("Error creating default configuration at '%s': %v.\n", *configurationFile, err)
			}
		}
	}
	isValid := appConf.verifyConfiguration()
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
