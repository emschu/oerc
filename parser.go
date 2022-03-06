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
