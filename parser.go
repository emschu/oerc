// oerc, alias oer-collector
// Copyright (C) 2021-2026 emschu[aet]mailbox.org
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
	"github.com/alitto/pond"
	"gorm.io/gorm"
	"log"
	"time"
)

// Parser common data structure of all parsers
type Parser struct {
	ChannelFamilyKey     string
	ChannelFamily        ChannelFamily
	db                   *gorm.DB
	dateRangeHandler     dateRangeHandler
	parallelWorkersCount int
}

// ParserInterface all parsers should implement this interface
type ParserInterface interface {
	handleDay(chn Channel, day time.Time) // process a single day for a single channel

	fetchTVShows() // handle tv shows

	postProcess() // run tasks after the program and tv show entries are fetched

	preProcess() bool // run tasks before the fetch process (like receiving API keys etc.), return success state

	isDateValidToFetch(day *time.Time) (bool, error)
}

type dateRangeHandler interface {
	getDateRange() *[]time.Time
}

type defaultDateRangeHandler struct {
	DaysInPast   uint
	DaysInFuture uint
}

type specificDateRangeHandler struct {
	StartDateTime time.Time
	EndDateTime   time.Time
}

func newDefaultDateRangeHandler() dateRangeHandler {
	return &defaultDateRangeHandler{
		DaysInPast:   GetAppConf().DaysInPast,
		DaysInFuture: GetAppConf().DaysInFuture,
	}
}

func newDefaultDateRangeHandlerPadded(paddingDays uint) dateRangeHandler {
	return &defaultDateRangeHandler{
		DaysInPast:   GetAppConf().DaysInPast - paddingDays,
		DaysInFuture: GetAppConf().DaysInFuture + paddingDays,
	}
}

func newSpecificDateRangeHandler(startDateTime time.Time, endDateTime time.Time) dateRangeHandler {
	return &specificDateRangeHandler{
		StartDateTime: startDateTime,
		EndDateTime:   endDateTime,
	}
}

func newSpecificDateRangeHandlerPadded(startDateTime time.Time, endDateTime time.Time, paddingDays uint) dateRangeHandler {
	return &specificDateRangeHandler{
		StartDateTime: startDateTime.Add(time.Duration(paddingDays) * -24 * time.Hour),
		EndDateTime:   endDateTime.Add(time.Duration(paddingDays) * 24 * time.Hour),
	}
}

func (d *defaultDateRangeHandler) getDateRange() *[]time.Time {
	return generateDateRangeInPastAndFuture(d.DaysInPast, d.DaysInFuture)
}

func (s *specificDateRangeHandler) getDateRange() *[]time.Time {
	return generateDateRangeBetweenDates(s.StartDateTime, s.EndDateTime)
}

// Fetch generic fetch function ready to handle all parsers, sets ChannelFamily and DB to the instance
func (p *Parser) Fetch(parserInterface ParserInterface) {
	// setup db
	db := getDb()
	p.db = db

	// get channel family db record and save it to the parser instance
	var channelFamily = getChannelFamily(db, p.ChannelFamilyKey)
	if channelFamily.ID == 0 {
		log.Fatalf("ChannelFamilyKey '%s' was not found!\n", p.ChannelFamilyKey)
		return
	}
	p.ChannelFamily = *channelFamily

	if parserInterface == nil {
		log.Fatalf("Incompliant parser instance received! Key: '%s'\n", p.ChannelFamilyKey)
	}

	isReady := parserInterface.preProcess()
	if !isReady {
		log.Printf("Parser is not ready to start\n")
		return
	}

	// import tv shows
	if GetAppConf().EnableTVShowCollection {
		parserInterface.fetchTVShows()
	}

	timeRange := p.dateRangeHandler.getDateRange()
	var times []time.Time
	if timeRange != nil {
		times = *timeRange
	} else {
		log.Printf("No valid time range received!\n")
		return
	}

	if GetAppConf().EnableProgramEntryCollection {
		// import program entries for the configured date range
		pool := pond.New(p.parallelWorkersCount, 100, getWorkerPoolIdleTimeout(), pond.PanicHandler(func(i interface{}) {
			log.Printf("Problem with goroutine pool: %v\n", i)
		}))
		for _, channel := range getChannelsOfFamily(db, channelFamily) {
			for _, day := range times {
				chn := channel
				dayToFetch := day
				ok, err := parserInterface.isDateValidToFetch(&dayToFetch)
				if ok && err == nil {
					pool.Submit(func() {
						parserInterface.handleDay(chn, dayToFetch)
					})
				} else {
					if verboseGlobal {
						log.Printf("Skipping date '%s' of channel #%s\n", dayToFetch.Format("2006-01-02"), chn.ChannelFamily.Title)
					}
				}
			}
		}
		// wait for finish
		pool.StopAndWait()

		// general tag linking
		parserInterface.postProcess()
	}

	if verboseGlobal {
		log.Printf("%s parsed successfully\n", p.ChannelFamilyKey)
	}
}

func (p *Parser) isMoreThanXDaysInFuture(day *time.Time, days uint) bool {
	if day == nil {
		return false
	}
	if days == 0 {
		return true
	}
	now := time.Now()
	return day.After(now) && day.Sub(now) > time.Duration(days)*24*time.Hour
}

func (p *Parser) isMoreThanXDaysInPast(day *time.Time, days uint) bool {
	if day == nil {
		return false
	}
	if days == 0 {
		return true
	}
	now := time.Now()
	return day.Before(now) && now.Sub(*day) > time.Duration(days)*24*time.Hour
}

func (p *Parser) logRecentFetch(customMessage string) {
	timeOfNextUpdate := getTimeOfNextUpdate()
	log.Printf("%s, due to recent fetch. Next regular fetch will be at: %s. "+
		"Use 'forceUpdate' = true to ignore this.",
		customMessage,
		timeOfNextUpdate.Format(time.RFC822),
	)
}
