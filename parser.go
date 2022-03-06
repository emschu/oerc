package main

import (
	"gorm.io/gorm"
	"time"
)

type ParserInterface interface {
	handleDay(chn Channel, day time.Time)
	fetchTVShows()
}

type DateRangeHandler interface {
	getDateRange() *[]time.Time
}

type DefaultDateRangeHandler struct {
	DaysInPast   uint
	DaysInFuture uint
}
type SpecificDateRangeHandler struct {
	StartDateTime time.Time
	EndDateTime   time.Time
}

func newDefaultDateRangeHandler() DateRangeHandler {
	return &DefaultDateRangeHandler{
		DaysInPast:   GetAppConf().DaysInPast,
		DaysInFuture: GetAppConf().DaysInFuture,
	}
}

func newSpecificDateRangeHandler(startDateTime time.Time, endDateTime time.Time) DateRangeHandler {
	return &SpecificDateRangeHandler{
		StartDateTime: startDateTime,
		EndDateTime:   endDateTime,
	}
}

func (d *DefaultDateRangeHandler) getDateRange() *[]time.Time {
	return generateDateRangeInPastAndFuture(d.DaysInPast, d.DaysInFuture)
}

func (s *SpecificDateRangeHandler) getDateRange() *[]time.Time {
	return generateDateRangeBetweenDates(s.StartDateTime, s.EndDateTime)
}

type Parser struct {
	ChannelFamilyKey string
	ChannelFamily    ChannelFamily
	db               *gorm.DB
	dateRangeHandler DateRangeHandler
}
