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
	"compress/gzip"
	"crypto/md5"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/microcosm-cc/bluemonday"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

var (
	newLineMatcher  = *regexp.MustCompile(`\r?\n+`)
	tabMatcher      = *regexp.MustCompile(`\t+`)
	wsMatcher       = *regexp.MustCompile(`\s+`)
	sanitizerPolicy = bluemonday.UGCPolicy()
	httpClient      = &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			Proxy:               getHTTPProxy(),
		},
		Timeout: time.Second * 45,
	}
)

const (
	settingKeyLastFetch             = "last_fetch_time"
	settingKeyLastSearch            = "last_search_time"
	settingKeyRequestsTotal         = "general_request_counter"
	settingKeyRequestsLastExecution = "last_request_counter"
)

var dBReference *gorm.DB

// connection closing logic should be handled by another part of the application, it's not implicit
func getDb() (*gorm.DB, error) {
	if dBReference == nil {
		conf := GetAppConf()

		gormLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: 3 * time.Second,
				LogLevel:      logger.Silent,
				Colorful:      true,
			},
		)

		if conf.DbType == "postgres" {
			var sslModeStr string
			if !conf.DbSSLEnabled {
				sslModeStr = "disable"
			} else {
				sslModeStr = "enable"
			}
			connectionString := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s TimeZone=%s",
				conf.DbHost, conf.DbPort, conf.DbUser, conf.DbName, conf.DbPassword, sslModeStr, appConf.TimeZone)
			db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
				SkipDefaultTransaction: true,
				DisableAutomaticPing:   false,
				Logger:                 gormLogger,
				PrepareStmt:            true,
			})
			if err != nil {
				log.Printf("Error connecting to the database. Is it running and configured correctly?\n")
				log.Fatal(err)
				return nil, err
			}
			s, err := db.DB()
			if err != nil {
				log.Printf("Error connecting to the database. Is it running and configured correctly?\n")
				log.Fatal(err)
				return nil, err
			}
			s.SetMaxOpenConns(50)

			dBReference = db
		} else {
			log.Fatalf("DbType '%s' is not implemented.", appConf.DbType)
		}
	}
	return dBReference, nil
}

// getHTTPProxy: method to get the http proxy of the app or nil of none configured
func getHTTPProxy() func(r *http.Request) (*url.URL, error) {
	// do proxy configuration
	proxyURLStr := GetAppConf().ProxyURL
	if len(strings.TrimSpace(proxyURLStr)) > 0 {
		proxyURL, err := url.ParseRequestURI(proxyURLStr)
		if err != nil {
			log.Printf("Error in proxy configuration url: %v\n", err)
			return http.ProxyFromEnvironment
		}
		if proxyURL.Scheme == "" || proxyURL.Host == "" || !strings.Contains(proxyURL.Host, ":") {
			log.Printf("Warning: Invalid proxy scheme/host or missing port detected!\n")
			return http.ProxyFromEnvironment
		}
		return http.ProxyURL(proxyURL)
	}
	return http.ProxyFromEnvironment
}

// central method to do http get requests in this application, retry count relates to (idle) connection count!
//goland:noinspection GoDeferInLoop
func doGetRequest(target string, requestHeaders map[string]string, retries uint) (*string, error) {
	// global http client
	var lastErr error
	for attempt := -1; attempt <= int(retries); attempt++ {
		req, err := http.NewRequest("GET", target, nil)
		// increment global request counter
		status.TotalRequests++
		if err != nil {
			lastErr = err
			incrErr()
			continue
		}
		// default headers
		req.Header.Add("Accept-Encoding", "gzip, deflate")
		req.Header.Add("Accept-Language", "de,en-US")
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; rv:60.0) Gecko/20100101 Firefox/60.0")

		// apply given headers
		if requestHeaders != nil {
			for key, value := range requestHeaders {
				if strings.ToLower(key) == "host" {
					req.Host = value
				} else {
					req.Header.Set(key, value)
				}
			}
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			lastErr = err
			incrErr()
			continue
		}
		if resp == nil {
			incrErr()
			continue
		}
		// close resp
		defer func() {
			err := resp.Body.Close()
			if err != nil {
				log.Printf("Error during closing http body: %v.\n", err)
			}
		}()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			log.Printf("Unexpected HTTP status code '%d' for url '%s'.\n", resp.StatusCode, target)
			lastErr = fmt.Errorf("invalid HTTP status code '%d'", resp.StatusCode)
			incrErr()
			continue
		}

		if attempt >= int(retries) {
			// leave
			incrErr()
			continue
		}
		// handle encoding-specific stuff
		var reader io.ReadCloser
		switch resp.Header.Get("Content-Encoding") {
		//goland:noinspection GoDeferInLoop
		case "gzip":
			if resp.Body != nil {
				reader, err = gzip.NewReader(resp.Body)
				if err != nil {
					log.Printf("problem reading gzip input. error: %v\n", err)
					lastErr = err
					incrErr()
					continue
				}
				defer func() {
					err := reader.Close()
					if err != nil {
						appLog(fmt.Sprintf("Problem during close of gzip stream: %v\n", err))
						incrErr()
					}
				}()
			}
		default:
			// no compression
			reader = resp.Body
		}
		body, _ := ioutil.ReadAll(reader)
		text := string(body)
		if len(text) > 0 {
			resetErr()
			return &text, nil
		}
	}
	appLog(fmt.Sprintf("Maximum number of retries reached for http request for url '%s'. Error: %v", target, lastErr))
	// fallback
	if lastErr != nil {
		return nil, lastErr
	}
	return nil, nil
}

// generateDateRange date range generator for today +- given days
func generateDateRange(daysInPast, daysInFuture uint) *[]time.Time {
	var dates []time.Time
	year, month, day := time.Now().Date()
	today, _ := time.Parse("2006-01-02T15:04:05", fmt.Sprintf("%04d-%02d-%02dT00:00:00", year, int8(month), day))

	// set time zone
	location, _ := time.LoadLocation(GetAppConf().TimeZone)
	today = today.In(location)

	dates = append(dates, today)
	for i := 1; i <= int(daysInFuture); i++ {
		dates = append(dates, today.AddDate(0, 0, i))
	}
	for i := 1; i <= int(daysInPast); i++ {
		dates = append(dates, today.AddDate(0, 0, -i))
	}
	// sort slice by date, from past to future
	sort.Slice(dates, func(i, j int) bool {
		return dates[i].Before(dates[j])
	})
	return &dates
}

// trimAndSanitizeString this function should be used for all user input (strings)
func trimAndSanitizeString(rawString string) string {
	res := newLineMatcher.ReplaceAllString(rawString, "")
	res = tabMatcher.ReplaceAllString(res, "")
	res = wsMatcher.ReplaceAllString(res, " ")
	res = sanitizeContent(&res)
	return strings.TrimSpace(res)
}

// appLog this function should be used to write log entries to the db log
func appLog(msg string) {
	db, _ := getDb()
	parsingError := &LogEntry{}
	parsingError.Message = trimAndSanitizeString(msg)
	db.Save(parsingError)

	log.Printf("error in parse process: '%s'", msg)
}

// saveProgramEntryRecord: method to store or create a program entry gorm record
func saveProgramEntryRecord(db *gorm.DB, programEntry *ProgramEntry) {
	if programEntry.ID != 0 {
		// limit description field
		if len(programEntry.Description) > 30000 {
			if isDebug() {
				panic(fmt.Errorf("Program Entry description too long %s", programEntry.URL))
			}
			programEntry.Description = programEntry.Description[0:30000]
		}

		// TODO handle different technical id case
		if verboseGlobal {
			log.Printf("Updating program entry #%d.\n", programEntry.ID)
		}
		db.Model(&programEntry).Updates(map[string]interface{}{
			"Description":     programEntry.Description,
			"StartDateTime":   programEntry.StartDateTime,
			"EndDateTime":     programEntry.EndDateTime,
			"DurationMinutes": programEntry.DurationMinutes,
			"Tags":            programEntry.Tags,
			"Title":           programEntry.Title,
			"Homepage":        programEntry.Homepage,
			"LastCheck":       time.Now(),
		}).Association("ImageLinks")
		status.TotalUpdatedPE++
	} else {
		now := time.Now()
		programEntry.LastCheck = &now

		if verboseGlobal {
			log.Printf("Create new program entry '%s' starting at '%s'.\n",
				programEntry.Title,
				programEntry.StartDateTime.Format(time.RFC3339))
		}
		db.Create(&programEntry)
		status.TotalCreatedPE++
		if db.Error != nil {
			log.Println(db.Error)
		}
	}
}

// getDocument: method to get a goquery object of a URL
func getDocument(apiURL string) (*goquery.Document, error) {
	res, err := doGetRequest(apiURL, map[string]string{}, 3)
	if res == nil || err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(*res))
	if err != nil {
		return nil, fmt.Errorf("problem fetching document URL '%s'. %v", apiURL, err)
	}
	return doc, nil
}

// saveTvShowRecord: general function to save a tv show record
func saveTvShowRecord(db *gorm.DB, tvShow *TvShow) {
	if tvShow.ID == 0 {
		// create new tv show
		db.Create(&tvShow)
		status.TotalCreatedTVS++
		if verboseGlobal {
			log.Printf("create tv show #%d \n", tvShow.ID)
		}
	} else {
		// tv show exists: update
		db.Model(&tvShow).Updates(TvShow{
			ManagedRecord: ManagedRecord{
				Title:    tvShow.Title,
				URL:      tvShow.URL,
				Homepage: tvShow.Homepage,
			},
		})
		status.TotalUpdatedTVS++
		if verboseGlobal {
			log.Printf("update tv show #%d \n", tvShow.ID)
		}
	}
}

// isRecentlyUpdated: method to check if program entry record was updated recently
func isRecentlyUpdated(entry *ProgramEntry) bool {
	if entry.LastCheck == nil || entry.LastCheck.IsZero() || GetAppConf().ForceUpdate {
		return false
	}
	secDiff := time.Since(*entry.LastCheck).Seconds()
	if secDiff < float64(GetAppConf().TimeToRefreshInMinutes*60) {
		return true
	}
	return false
}

// isRecentlyFetched method to check if there was a recent fetch job
func isRecentlyFetched() bool {
	if GetAppConf().ForceUpdate {
		return false
	}
	set := getSetting(settingKeyLastFetch)
	if set != nil && set.ID != 0 {
		lastUpdateTime, err := time.Parse(time.RFC3339, set.Value)
		if len(set.Value) == 0 {
			return false
		}
		if err != nil {
			log.Printf("Could not parse '%s' as date", set.Value)
			return false
		}

		location, _ := time.LoadLocation(GetAppConf().TimeZone)
		lastUpdateTime = lastUpdateTime.In(location)
		minDiff := time.Since(lastUpdateTime).Minutes()
		if minDiff < float64(GetAppConf().TimeToRefreshInMinutes) {
			return true
		}
	} else {
		log.Printf("Could not find setting value of last fetch\n")
	}
	return false
}

// buildHash: Method to build a hash of the given string parts
func buildHash(in []string) string {
	var val = strings.Join(in, ":")
	return fmt.Sprintf("%x", md5.Sum([]byte(val)))
}

// considerTagExists: adds a tag to the program entry.
func considerTagExists(programEntry *ProgramEntry, mainTagName *string) {
	var existingTags []string
	if len(programEntry.Tags) > 0 {
		existingTags = strings.Split(programEntry.Tags, ";")
	} else {
		existingTags = []string{}
	}
	for _, tag := range existingTags {
		if tag == *mainTagName {
			return
		}
	}
	existingTags = append(existingTags, trimAndSanitizeString(*mainTagName))
	if verboseGlobal {
		log.Printf("Save new tag '%s' to program entry #%d\n", *mainTagName, programEntry.ID)
	}

	programEntry.Tags = strings.Join(existingTags, ";")
}

// this definition is important for default values
var settings = map[string]string{
	settingKeyLastFetch:             "",
	settingKeyRequestsTotal:         "0",
	settingKeyRequestsLastExecution: "0",
}

func getSetting(key string) *Settings {
	return getOrCreateSetting(key)
}

func getOrCreateSetting(key string) *Settings {
	db, _ := getDb()

	var setting Settings
	db.Model(&Settings{}).Where("setting_key = ?", key).Find(&setting)
	if setting.ID == 0 {
		// get default defaultValue of definition above
		defaultValue := settings[key]
		db.Model(&Settings{}).Create(&Settings{SettingKey: key, Value: defaultValue})
		db.Model(&Settings{}).Where("setting_key = ?", key).Find(&setting)
	}
	// FIXME!
	if setting.ID == 0 {
		log.Fatal("Cannot find setting object in db")
	}
	return &setting
}

func setSetting(key string, val string) {
	setting := getOrCreateSetting(key)
	db, _ := getDb()
	db.Model(&setting).Update("value", val)
}

// method to sanitize a string
func sanitizeContent(content *string) string {
	return sanitizerPolicy.Sanitize(*content)
}

// error handling vars
var (
	appInError     = false
	errorCounter   = 0
	errorThreshold = 20
)

// incrErr: global error counter increment
func incrErr() {
	errorCounter++
}

// checkErr: global error check method, if counter is too high, exit application with error exit code 1
func checkErr() {
	if !appInError && errorCounter > errorThreshold {
		appInError = true
	}
	if appInError {
		appLog(fmt.Sprintf("Too much HTTP errors (>%d). Cancel program. Please connect to the internet.", errorThreshold))
		os.Exit(1)
	}
}

// resetErr: global error counter reset method
func resetErr() {
	errorCounter = 0
}

// general connectivity check, should be called on startup of the fetch process
func connectivityCheck() (bool, error) {
	hostsToCheck := &[]string{
		ardHostWithPrefix,
		zdfHost,
		orfHostWithPrefix,
		srfHostWithPrefix,
	}
	for _, v := range *hostsToCheck {
		response, err := doGetRequest(v, map[string]string{}, 1)
		if err != nil {
			log.Println("Are you connected to the internet?")
			return false, fmt.Errorf("general network internet connectivity check failed. Message: %v", err)
		}
		if response == nil {
			log.Println("Are you connected to the internet?")
			return false, fmt.Errorf("general network internet connectivity check failed")
		}
	}
	return true, nil
}

// method to get general colly collector instance used by ard, srf and orf parsers
func baseCollector(allowedHost []string) *colly.Collector {
	c := colly.NewCollector(colly.AllowedDomains(allowedHost...), colly.AllowURLRevisit())
	c.MaxDepth = 1
	c.Async = true
	c.CheckHead = false
	c.TraceHTTP = false
	c.UserAgent = "Mozilla/5.0 (Windows NT 6.1; rv:60.0) Gecko/20100101 Firefox/60.0"
	c.SetRequestTimeout(30 * time.Second)

	// configure proxy
	proxyURL := GetAppConf().ProxyURL
	if len(proxyURL) > 0 {
		proxyErr := c.SetProxy(proxyURL)
		if proxyErr != nil {
			panic(proxyErr)
		}
	}

	for _, singleHost := range allowedHost {
		limitErr := c.Limit(&colly.LimitRule{
			DomainGlob:  singleHost + "/**",
			Parallelism: 4,
			Delay:       25 * time.Millisecond,
		})
		if limitErr != nil {
			log.Fatalf("%s\n", limitErr)
		}
	}
	c.OnRequest(func(r *colly.Request) {
		if verboseGlobal {
			log.Printf("Visiting '%s'\n", r.URL)
		}
		status.TotalRequests++
	})
	c.OnResponse(func(response *colly.Response) {
		resetErr()
	})
	c.OnError(func(response *colly.Response, err error) {
		incrErr()
		checkErr()
	})
	return c
}

// ClearLogs method to clear the application's logs
func ClearLogs() {
	db, _ := getDb()

	db.Where("id > 0").Delete(&LogEntry{})
}

// ClearAll method to clear (almost) all the db data - except channels + channels families
func ClearAll() {
	ClearLogs()
	ClearRecommendations()

	db, _ := getDb()
	db.Where("id > 0").Delete(&TvShow{})
	db.Where("id > 0").Delete(&ImageLink{})
	db.Where("id > 0").Delete(&ProgramEntry{})
	db.Where("id > 0").Delete(&Settings{})
}

// ClearRecommendations method to clear ALL the recommendations from the database
func ClearRecommendations() {
	db, _ := getDb()

	db.Where("id > 0").Delete(&Recommendation{})
}

// ClearOldRecommendations method to clear all the OLD(=past) recommendations
func ClearOldRecommendations() {
	db, _ := getDb()

	db.Where("start_date_time < ?", time.Now()).Delete(&Recommendation{})
}
