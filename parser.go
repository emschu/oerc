package main

import (
	"gorm.io/gorm"
	"time"
)

// ParserInterface all parsers should implement this interface
type ParserInterface interface {
	Fetch() // is called in main

	handleDay(chn Channel, day time.Time) // process a single day for a single channel

	fetchTVShows() // handle tv shows
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

func newSpecificDateRangeHandler(startDateTime time.Time, endDateTime time.Time) dateRangeHandler {
	return &specificDateRangeHandler{
		StartDateTime: startDateTime,
		EndDateTime:   endDateTime,
	}
}

func (d *defaultDateRangeHandler) getDateRange() *[]time.Time {
	return generateDateRangeInPastAndFuture(d.DaysInPast, d.DaysInFuture)
}

func (s *specificDateRangeHandler) getDateRange() *[]time.Time {
	return generateDateRangeBetweenDates(s.StartDateTime, s.EndDateTime)
}

// Parser common data structure of all parsers
type Parser struct {
	ChannelFamilyKey string
	ChannelFamily    ChannelFamily
	db               *gorm.DB
	dateRangeHandler dateRangeHandler
}
