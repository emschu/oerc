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

// this file contains all HTTP JSON API methods for the gin server implementation of this project

import (
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// use this method to get program entries for the web frontend, includes deprecated entries
func getProgramOfWeb(start *time.Time, end *time.Time, channel *Channel) *ProgramResponse {
	db := getDb()
	var entries []ProgramEntry
	// 14 day = max range
	var endDateTime = *end
	if end.Sub(*start).Hours()/24 > 14 {
		endLimit := start.Add(14 * 24 * time.Hour)
		endDateTime = endLimit
	}
	order := db.Model(&ProgramEntry{}).Where("start_date_time between ? and ?", start, endDateTime).
		Preload("ImageLinks").
		Preload("CollisionEntries").
		Order("channel_id")
	if channel != nil {
		order.Where("channel_id", channel.ID)
	}
	result := order.Find(&entries)
	if result.Error != nil {
		log.Fatalf("error fetching program items: %v", result.Error)
		return nil
	}
	response := ProgramResponse{
		From:             start,
		To:               end,
		ChannelID:        0,
		Size:             len(entries),
		ProgramEntryList: &entries,
	}
	if channel != nil {
		response.ChannelID = int64(channel.ID)
	}
	return &response
}

func getChannels() *[]Channel {
	db := getDb()
	var channels []Channel
	result := db.Model(&Channel{}).Preload("ChannelFamily").Where("is_deprecated = ?", false).Order("priority asc").Find(&channels)
	if result.Error != nil {
		log.Fatalf("error fetching channels: %v", result.Error)
		return nil
	}
	return &channels
}

func getChannelFamilies() *[]ChannelFamily {
	db := getDb()
	families := &[]ChannelFamily{}
	result := db.Model(&ChannelFamily{}).Find(families)
	if result.Error != nil {
		log.Fatalf("error fetching channel channelFamilyKeys: %v", result.Error)
		return nil
	}
	return families
}

func getProgramYesterdayHandler(c *gin.Context) {
	y := time.Now().Add(-24 * time.Hour)
	yStart := time.Date(y.Year(), y.Month(), y.Day(), 0, 0, 0, 0, y.Location())
	yEnd := time.Date(y.Year(), y.Month(), y.Day(), 23, 59, 59, 0, y.Location())

	c.JSON(http.StatusOK, getProgramOfWeb(&yStart, &yEnd, nil))
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

	c.JSON(http.StatusOK, getProgramOfWeb(&yStart, &yEnd, channel))
}

func getProgramTodayHandler(c *gin.Context) {
	now := time.Now()
	tStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	tEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

	c.JSON(http.StatusOK, getProgramOfWeb(&tStart, &tEnd, nil))
}

func getProgramTomorrowHandler(c *gin.Context) {
	tom := time.Now().Add(+24 * time.Hour)
	toStart := time.Date(tom.Year(), tom.Month(), tom.Day(), 0, 0, 0, 0, tom.Location())
	toEnd := time.Date(tom.Year(), tom.Month(), tom.Day(), 23, 59, 59, 0, tom.Location())

	c.JSON(http.StatusOK, getProgramOfWeb(&toStart, &toEnd, nil))
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

	c.JSON(http.StatusOK, getProgramOfWeb(&toStart, &toEnd, channel))
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

	c.JSON(http.StatusOK, getProgramOfWeb(&tStart, &tEnd, channel))
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

	end = end.In(location)
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
	c.JSON(http.StatusOK, getProgramOfWeb(&start, &end, channel))
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

		db := getDb()
		var channel Channel
		result := db.Model(&Channel{}).Preload("ChannelFamily").Where("is_deprecated is false AND id = ?", uint(cid)).First(&channel)
		if result.Error != nil {
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
	c.JSON(http.StatusBadRequest, Error{
		Status:  "400",
		Message: "missing or invalid channel id",
	})
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

	var statusInfoModel = StatusInfoModel{}

	statusViewQueryErr := db.Model(StatusInfoModel{}).Find(&statusInfoModel)
	if statusViewQueryErr.Error != nil {
		log.Printf("ERROR fetching information from status info table: %s", statusViewQueryErr.Error.Error())
	}

	var response = StatusResponse{
		ChannelFamilyCount:  statusInfoModel.ChannelFamilyCount,
		ChannelCount:        statusInfoModel.ChannelCount,
		ImageLinksCount:     statusInfoModel.ImageLinkCount,
		ProgramEntryCount:   statusInfoModel.ProgramEntryCount,
		TvShowCount:         statusInfoModel.TvShowCount,
		LogCount:            statusInfoModel.LogCount,
		RecommendationCount: statusInfoModel.RecommendationCount,
		Version:             version,
		ServerDateTime:      time.Now().Format(time.RFC3339),
		DataStartTime:       statusInfoModel.DataStartTime,
		DataEndTime:         statusInfoModel.DataEndTime,
		TvChannels:          channels,
		TvChannelFamilies:   channelFamilies,
	}
	return &response
}

func getChannelsHandler(c *gin.Context) {
	channels := getChannels()
	c.JSON(http.StatusOK, ChannelResponse{channels, len(*channels)})
}

func putChannelsHandler(c *gin.Context) {
	var channels []Channel
	if err := c.ShouldBindJSON(&channels); err != nil || channels == nil {
		c.JSON(http.StatusBadRequest, Error{Status: "400", Message: "Invalid channel data"})
		return
	}

	db := getDb()
	for i, channel := range channels {
		db.Model(&Channel{}).Where("id = ?", channel.ID).Update("priority", i)
	}

	updatedChannels := getChannels()
	c.JSON(http.StatusOK, ChannelResponse{updatedChannels, len(*updatedChannels)})
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

func clearAllLogEntriesHandler(context *gin.Context) {
	ClearLogs()
	context.JSON(http.StatusOK, "OK")
}

func getRecommendationsHandler(context *gin.Context) {
	location, _ := time.LoadLocation(GetAppConf().TimeZone)

	fromStr := context.Query("from")
	var from time.Time

	now := time.Now()
	if len(fromStr) == 0 {
		from = now
	} else {
		var err error
		from, err = time.Parse(time.RFC3339, fromStr)
		if err != nil {
			context.JSON(http.StatusBadRequest, Error{Status: "400", Message: "Invalid string in 'from' parameter given"})
			return
		}
	}
	from = from.In(location)
	logEntryList := getRecommendationsAt(from)

	context.JSON(http.StatusOK, &logEntryList)
}

func getRecommendationsAt(at time.Time) []Recommendation {
	location, _ := time.LoadLocation(GetAppConf().TimeZone)
	minuteDiff := at.Sub(time.Now().In(location)).Minutes()

	db := getDb()
	var logEntryList []Recommendation
	dbQuery := db.Model(&Recommendation{}).
		Select("recommendations.*").
		Joins("LEFT JOIN program_entries ON (recommendations.program_entry_id = program_entries.id)").
		Order("recommendations.start_date_time asc").
		Preload("ProgramEntry").
		Preload("ProgramEntry.ImageLinks").
		Preload("ProgramEntry.CollisionEntries")

	if minuteDiff < 15 {
		// include televised items of this moment, but don't do this for requests of future recommendations
		dbQuery = dbQuery.Where("((recommendations.start_date_time <= ? AND program_entries.end_date_time >= ?) OR program_entries.start_date_time >= ?)", at, at, at)
	} else {
		dbQuery = dbQuery.Where("(recommendations.start_date_time >= ?)", at)
	}
	dbQuery.Find(&logEntryList)

	return logEntryList
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
	if db.Dialector.Name() == "sqlite" {
		db.Model(&ProgramEntry{}).Where("start_date_time >= datetime('now', '-1 day') AND (UPPER(title) LIKE ? OR UPPER(description) LIKE ?)", strings.ToUpper(queryStr), strings.ToUpper(queryStr)).
			Offset(int(offset)).Limit(int(limit)).
			Order("start_date_time ASC").
			Preload("ImageLinks").
			Preload("CollisionEntries").
			Find(&programEntryList)
	} else {
		db.Model(&ProgramEntry{}).Where("start_date_time >= (NOW() - interval '1 day') AND (title ILIKE ? OR description ILIKE ?)", queryStr, queryStr).
			Offset(int(offset)).Limit(int(limit)).
			Order("start_date_time ASC").
			Preload("ImageLinks").
			Preload("CollisionEntries").
			Find(&programEntryList)
	}

	if len(programEntryList) == 0 {
		context.JSON(http.StatusOK, []ProgramEntry{})
		return
	}
	context.JSON(http.StatusOK, programEntryList)
}

func getXMLTvHandler(context *gin.Context) {
	context.Writer.Header().Set("Content-Type", "application/xml; charset=utf-8")

	from := context.Query("from")
	to := context.Query("to")

	if len(from) == 0 || len(to) == 0 {
		context.XML(http.StatusBadRequest, Error{
			Status:  "400",
			Message: "invalid empty date range parameters: 'from' or 'to' is missing",
		})
		return
	}

	var err error
	var start, end time.Time
	start, err = time.Parse(time.RFC3339, from)
	if err != nil || start.IsZero() {
		context.XML(http.StatusBadRequest, Error{
			Status:  "400",
			Message: "invalid 'from' date time parameter",
		})
		return
	}
	location, _ := time.LoadLocation(GetAppConf().TimeZone)
	start = start.In(location)

	end, err = time.Parse(time.RFC3339, to)
	if err != nil || end.IsZero() {
		context.XML(http.StatusBadRequest, Error{
			Status:  "400",
			Message: "invalid 'to' date time parameter",
		})
		return
	}

	end = end.In(location)
	if end.Before(start) || end.Equal(start) {
		context.XML(http.StatusBadRequest, Error{
			Status:  "400",
			Message: "invalid date range, end is before or equal to start",
		})
		return
	}

	xmltv, err := exportToXMLTV(start, end)
	if err != nil {
		context.XML(http.StatusInternalServerError, Error{
			Status:  "500",
			Message: "Error while generating XMLTV file",
		})
		return
	}
	context.XML(http.StatusOK, &xmltv)
}
