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

// this file contains all HTTP JSON API methods for the gin server implementation of this project

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

// getProgramOf generating a program entry list response for each channel (or if channel = nil for ALL channels) in given time range. returns ProgramResponse
func getProgramOf(start *time.Time, end *time.Time, channel *Channel) *ProgramResponse {
	db := getDb()
	var entries []ProgramEntry
	order := db.Model(&ProgramEntry{}).Where("start_date_time between ? and ?", start, end).
		Preload("ImageLinks").
		Preload("CollisionEntries").
		Order("channel_id")
	if channel != nil {
		order.Where("channel_id", channel.ID)
	}
	order.Find(&entries)
	response := ProgramResponse{
		From:             start,
		To:               end,
		Size:             len(entries),
		ProgramEntryList: &entries,
	}
	return &response
}

func getChannels() *[]Channel {
	db := getDb()
	var channels []Channel
	db.Model(&Channel{}).Find(&channels)
	return &channels
}

func getChannelFamilies() *[]ChannelFamily {
	db := getDb()
	families := &[]ChannelFamily{}
	db.Model(&ChannelFamily{}).Find(families)
	return families
}

func getCount(model interface{}) uint64 {
	db := getDb()
	var count int64
	db.Model(model).Count(&count)
	return uint64(count)
}

func getImageLinksCount() uint64 {
	return getCount(&ImageLink{})
}

func getProgramEntryCount() uint64 {
	return getCount(&ProgramEntry{})
}

func getTvShowCount() uint64 {
	return getCount(&TvShow{})
}

func getProgramYesterdayHandler(c *gin.Context) {
	y := time.Now().Add(-24 * time.Hour)
	yStart := time.Date(y.Year(), y.Month(), y.Day(), 0, 0, 0, 0, y.Location())
	yEnd := time.Date(y.Year(), y.Month(), y.Day(), 23, 59, 59, 0, y.Location())

	c.JSON(http.StatusOK, getProgramOf(&yStart, &yEnd, nil))
}

func getYesterdayProgramWithChannelHandler(c *gin.Context) {
	cid := c.Param("channel_id")
	channel, isValid := isChannelValid(c, cid, false)
	if !isValid {
		return
	}

	y := time.Now().Add(-24 * time.Hour)
	yStart := time.Date(y.Year(), y.Month(), y.Day(), 0, 0, 0, 0, y.Location())
	yEnd := time.Date(y.Year(), y.Month(), y.Day(), 23, 59, 59, 0, y.Location())

	c.JSON(http.StatusOK, getProgramOf(&yStart, &yEnd, channel))
}

func getProgramTodayHandler(c *gin.Context) {
	now := time.Now()
	tStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	tEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

	c.JSON(http.StatusOK, getProgramOf(&tStart, &tEnd, nil))
}

func getProgramTomorrowHandler(c *gin.Context) {
	tom := time.Now().Add(+24 * time.Hour)
	toStart := time.Date(tom.Year(), tom.Month(), tom.Day(), 0, 0, 0, 0, tom.Location())
	toEnd := time.Date(tom.Year(), tom.Month(), tom.Day(), 23, 59, 59, 0, tom.Location())

	c.JSON(http.StatusOK, getProgramOf(&toStart, &toEnd, nil))
}

func getTomorrowProgramWithChannelHandler(c *gin.Context) {
	cid := c.Param("channel_id")

	channel, isValid := isChannelValid(c, cid, false)
	if !isValid {
		return
	}

	tom := time.Now().Add(+24 * time.Hour)
	toStart := time.Date(tom.Year(), tom.Month(), tom.Day(), 0, 0, 0, 0, tom.Location())
	toEnd := time.Date(tom.Year(), tom.Month(), tom.Day(), 23, 59, 59, 0, tom.Location())

	c.JSON(http.StatusOK, getProgramOf(&toStart, &toEnd, channel))
}

func getSingleChannelHandler(c *gin.Context) {
	var input = c.Param("channel_id")
	channel, isValid := isChannelValid(c, input, false)
	if !isValid {
		return
	}
	c.JSON(http.StatusOK, channel)
}

func getDailyProgramWithChannelHandler(c *gin.Context) {
	cid := c.Param("channel_id")
	channel, isValid := isChannelValid(c, cid, false)
	if !isValid {
		return
	}

	now := time.Now()
	tStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	tEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

	c.JSON(http.StatusOK, getProgramOf(&tStart, &tEnd, channel))
}

func getProgramHandler(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")

	cid := c.Query("channel_id")
	from := c.Query("from")
	to := c.Query("to")

	if len(from) == 0 || len(to) == 0 {
		c.JSON(http.StatusBadRequest, Error{
			Status:  "400",
			Message: "invalid empty query parameters",
		})
		return
	}

	var err error
	channel, isValid := isChannelValid(c, cid, true)
	if !isValid {
		return
	}
	var start, end time.Time
	start, err = time.Parse(time.RFC3339, from)
	if err != nil || start.IsZero() {
		c.JSON(http.StatusBadRequest, Error{
			Status:  "400",
			Message: "invalid 'from' date time parameter",
		})
		return
	}
	location, _ := time.LoadLocation(GetAppConf().TimeZone)
	start = start.In(location)

	end, err = time.Parse(time.RFC3339, to)
	if err != nil || end.IsZero() {
		c.JSON(http.StatusBadRequest, Error{
			Status:  "400",
			Message: "invalid 'to' date time parameter",
		})
		return
	}

	end.In(location)
	if end.Before(start) || end.Equal(start) {
		c.JSON(http.StatusBadRequest, Error{
			Status:  "400",
			Message: "invalid date range, end is before or equal to start",
		})
		return
	}
	if end.Sub(start).Hours() > 21*24 {
		c.JSON(http.StatusBadRequest, Error{
			Status:  "400",
			Message: "invalid date range, exceeds limit of 21 days",
		})
		return
	}
	c.JSON(http.StatusOK, getProgramOf(&start, &end, channel))
}

func getSingleProgramEntryHandler(c *gin.Context) {
	pID := c.Param("id")
	pEID, err := strconv.ParseInt(pID, 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, Error{Status: "404", Message: "Invalid program entry id"})
		return
	}

	db := getDb()
	var programEntry ProgramEntry
	db.Model(ProgramEntry{}).Preload("ImageLinks").Preload("CollisionEntries").First(&programEntry, pEID)
	if programEntry.ID == 0 {
		c.JSON(http.StatusNotFound, Error{Status: "404", Message: "Invalid program entry id"})
		return
	}
	c.JSON(http.StatusOK, programEntry)
}

// isChannelValid: helper method to check if a given channel id is valid and exists, if existent a pointer to a Channel object is returned
func isChannelValid(c *gin.Context, cid string, acceptZero bool) (*Channel, bool) {
	var channel Channel
	// parse channel id
	if len(cid) > 0 {
		cid, err := strconv.ParseInt(cid, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, Error{
				Status:  "400",
				Message: "invalid channel id integer64 parameter",
			})
			return nil, false
		}
		if acceptZero && cid == 0 {
			// 0 = "all"
			return nil, true
		}
		if !acceptZero && cid == 0 {
			c.JSON(http.StatusNotFound, Error{
				Status:  "404",
				Message: "invalid channel id",
			})
			return nil, false
		}

		// check if channel exists
		channelExists := false
		channels := getChannels()
		for _, c := range *channels {
			if c.ID == uint(cid) {
				channelExists = true
				channel = c
				break
			}
		}
		if !channelExists {
			c.JSON(http.StatusNotFound, Error{
				Status:  "404",
				Message: "invalid channel id",
			})
			return nil, false
		}
		return &channel, true
	}
	if acceptZero {
		return nil, true
	}
	return nil, false
}

func getStatusHandler(c *gin.Context) {
	response := getStatusObject()
	if response == nil {
		return
	}
	c.JSON(http.StatusOK, response)
}

func getStatusObject() *StatusResponse {
	channels := getChannels()
	channelFamilies := getChannelFamilies()

	if channels == nil || channelFamilies == nil {
		return nil
	}

	db := getDb()

	var firstEntry time.Time
	var peCount int64
	var lastEntry time.Time

	db.Model(&ProgramEntry{}).Count(&peCount)

	if peCount > 0 {
		errMin := db.Raw("SELECT MIN(start_date_time) from program_entries LIMIT 1").Row().Scan(&firstEntry)
		if errMin != nil {
			log.Fatal("error querying database for MIN(start_date_time)")
		}
		errMax := db.Raw("SELECT MAX(end_date_time) from program_entries LIMIT 1").Row().Scan(&lastEntry)
		if errMax != nil {
			log.Fatal("error querying database for MAX(end_date_time)")
		}
	}

	var firstEntryStr string
	if firstEntry.IsZero() {
		firstEntryStr = ""
	} else {
		firstEntryStr = firstEntry.Format(time.RFC3339)
	}
	var lastEntryStr string
	if lastEntry.IsZero() {
		lastEntryStr = ""
	} else {
		lastEntryStr = lastEntry.Format(time.RFC3339)
	}

	var response = StatusResponse{
		ChannelFamilyCount:  uint64(len(*channelFamilies)),
		ChannelCount:        uint64(len(*channels)),
		ImageLinksCount:     getImageLinksCount(),
		ProgramEntryCount:   getProgramEntryCount(),
		TvShowCount:         getTvShowCount(),
		ProblemCount:        getCount(&LogEntry{}),
		RecommendationCount: getCount(&Recommendation{}),
		Version:             version,
		ServerDateTime:      time.Now().Format(time.RFC3339),
		DataStartTime:       firstEntryStr,
		DataEndTime:         lastEntryStr,
		TvChannels:          channels,
		TvChannelFamilies:   channelFamilies,
	}
	return &response
}

func getChannelsHandler(c *gin.Context) {
	channels := getChannels()
	c.JSON(http.StatusOK, ChannelResponse{channels, len(*channels)})
}

func getLogEntriesHandler(context *gin.Context) {
	db := getDb()
	var logEntryList []LogEntry
	var entryCount, pageCount int64
	db.Model(&LogEntry{}).Count(&entryCount)
	pageCount = int64(math.Ceil(float64(entryCount)))
	db.Model(&LogEntry{}).Limit(500).Order("id desc").Find(&logEntryList)
	context.JSON(http.StatusOK, LogEntriesResponse{&logEntryList, int64(len(logEntryList)), 0, pageCount, entryCount})
}

func getSingleLogEntriesHandler(context *gin.Context) {
	logID := context.Param("id")
	logEntryID, err := strconv.ParseInt(logID, 10, 64)
	if err != nil {
		context.JSON(http.StatusNotFound, Error{Status: "404", Message: "Invalid log entry id"})
		return
	}

	db := getDb()
	var singleLogEntry LogEntry
	db.Model(&LogEntry{}).Where("id", logEntryID).Find(&singleLogEntry)
	if singleLogEntry.ID == 0 {
		context.JSON(http.StatusNotFound, Error{Status: "404", Message: "Invalid log entry id"})
		return
	}
	context.JSON(http.StatusOK, &singleLogEntry)
}

func deleteSingleLogEntriesHandler(context *gin.Context) {
	logID := context.Param("id")
	logEntryID, err := strconv.ParseInt(logID, 10, 64)
	if err != nil {
		context.JSON(http.StatusNotFound, Error{Status: "404", Message: "Invalid log entry id"})
		return
	}

	db := getDb()
	var singleLogEntry LogEntry
	db.Model(&LogEntry{}).Where("id", logEntryID).Find(&singleLogEntry)
	if singleLogEntry.ID == 0 {
		context.JSON(http.StatusNotFound, Error{Status: "404", Message: "Invalid log entry id"})
		return
	}
	db.Delete(&singleLogEntry)
	context.JSON(http.StatusOK, "OK")
}

func clearAllLogEntriesHandler(context *gin.Context) {
	ClearLogs()
	context.JSON(http.StatusOK, "OK")
}

func getRecommendationsHandler(context *gin.Context) {
	location, _ := time.LoadLocation(GetAppConf().TimeZone)

	fromStr := context.Query("from")
	var from time.Time

	if len(fromStr) == 0 {
		from = time.Now()
	} else {
		var err error
		from, err = time.Parse(time.RFC3339, fromStr)
		if err != nil {
			context.JSON(http.StatusBadRequest, Error{Status: "400", Message: "Invalid string in 'from' parameter given"})
			return
		}
	}
	from = from.In(location)
	threeDaysInFuture := from.Add(48 * time.Hour) // maximum
	fourHoursBefore := from.Add(-4 * time.Hour)   // minimum

	db := getDb()
	var logEntryList []Recommendation
	db.Debug().
		Model(&Recommendation{}).
		Select("recommendations.*").
		Joins("LEFT JOIN program_entries ON (recommendations.program_entry_id = program_entries.id)").
		Where("recommendations.start_date_time >= ? AND recommendations.start_date_time <= ? AND program_entries.start_date_time >= ?", fourHoursBefore, threeDaysInFuture, from).
		Order("recommendations.start_date_time asc").
		Preload("ProgramEntry").
		Preload("ProgramEntry.ImageLinks").
		Preload("ProgramEntry.CollisionEntries").
		Find(&logEntryList)

	context.JSON(http.StatusOK, &logEntryList)
}

func getSearchHandler(context *gin.Context) {
	// get query string param
	queryStr := trimAndSanitizeString(context.Query("query"))
	if queryStr == "" {
		context.JSON(http.StatusBadRequest, Error{
			Status:  "400",
			Message: "Empty query string received!",
		})
		return
	}
	if len(queryStr) < 3 || len(queryStr) > 150 {
		context.JSON(http.StatusBadRequest, Error{
			Status:  "400",
			Message: "Query string length has to be between 3 and 150 characters!",
		})
		return
	}
	queryStr = "%" + queryStr + "%"

	// handle limit param
	limitStr := context.Query("limit")
	var limit uint64
	if len(limitStr) > 0 {
		var err error
		limit, err = strconv.ParseUint(limitStr, 10, 64)
		if err != nil {
			context.JSON(http.StatusBadRequest, Error{
				Status:  "400",
				Message: "Invalid value for parameter 'limit'. Should be a positive integer.",
			})
			return
		}
	}

	if limit > 5000 {
		// internal maximum
		limit = 5000
	}
	if limit == 0 || len(limitStr) == 0 {
		limit = 250
	}

	// handle offset param
	offsetStr := context.Query("offset")
	var offset uint64
	if len(offsetStr) > 0 {
		var err error
		offset, err = strconv.ParseUint(offsetStr, 10, 64)
		if err != nil {
			context.JSON(http.StatusBadRequest, Error{
				Status:  "400",
				Message: "Invalid value for parameter 'offset'. Should be a positive integer.",
			})
			return
		}
	} else {
		offset = 0
	}

	db := getDb()
	var programEntryList []ProgramEntry
	db.Model(&ProgramEntry{}).Where("start_date_time >= NOW() AND (title ILIKE ? OR description ILIKE ?)", queryStr, queryStr).
		Offset(int(offset)).Limit(int(limit)).
		Order("start_date_time ASC").
		Preload("ImageLinks").
		Preload("CollisionEntries").
		Find(&programEntryList)

	if len(programEntryList) == 0 {
		context.JSON(http.StatusOK, [0]ProgramEntry{})
		return
	}
	context.JSON(http.StatusOK, &programEntryList)
}
