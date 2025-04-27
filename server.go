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
	rice "github.com/GeertJohan/go.rice"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

// StartServer method to start the built-in gin web server to serve the JSON Api
func StartServer() {
	if isDebug() {
		gin.SetMode("debug")
	} else {
		gin.SetMode("release")
	}

	go setupMaterializedView()

	r := initRouter()
	log.Printf("Starting API server...\n")
	ip := net.ParseIP(GetAppConf().ServerHost)
	var err error
	if ip.To4() == nil {
		// ipv6
		err = r.Run(fmt.Sprintf("[%s]:%d", GetAppConf().ServerHost, GetAppConf().ServerPort))
	} else {
		err = r.Run(fmt.Sprintf("%s:%d", GetAppConf().ServerHost, GetAppConf().ServerPort))
	}
	if err != nil {
		log.Printf("Problem starting server %v\n", err)
		return
	}
}

func setupMaterializedView() {
	db := getDb()
	db.Exec(`drop materialized view status_info`)
	db.Exec(fmt.Sprintf(`create materialized view status_info as %s`, materializedStatusView))
	if isDebug() {
		log.Printf("Materialized status view")
	}
}

// StatusInfoModel a gorm model for the materialized view
type StatusInfoModel struct {
	ChannelFamilyCount  uint64
	ChannelCount        uint64
	ProgramEntryCount   uint64
	TvShowCount         uint64
	ImageLinkCount      uint64
	LogCount            uint64
	RecommendationCount uint64
	DataStartTime       string
	DataEndTime         string
}

// TableName of the materialized view prefixed by db schema
func (s *StatusInfoModel) TableName() string {
	return fmt.Sprintf("%s.status_info", appConf.DbSchema)
}

// StatusResponse api object
type StatusResponse struct {
	ChannelFamilyCount  uint64           `json:"channel_family_count"`
	ChannelCount        uint64           `json:"channel_count"`
	ProgramEntryCount   uint64           `json:"program_entry_count"`
	TvShowCount         uint64           `json:"tv_show_count"`
	ImageLinksCount     uint64           `json:"image_links_count"`
	LogCount            uint64           `json:"log_count"`
	RecommendationCount uint64           `json:"recommendation_count"`
	Version             string           `json:"version"`
	ServerDateTime      string           `json:"server_date_time"`
	DataStartTime       string           `json:"data_start_time"`
	DataEndTime         string           `json:"data_end_time"`
	TvChannels          *[]Channel       `json:"tv_channels"`
	TvChannelFamilies   *[]ChannelFamily `json:"tv_channel_families"`
}

// ProgramResponse program response object
type ProgramResponse struct {
	From             *time.Time      `json:"from"`
	To               *time.Time      `json:"to"`
	ChannelID        int64           `json:"channel_id"`
	Size             int             `json:"size"`
	ProgramEntryList *[]ProgramEntry `json:"program_list"`
}

// ChannelResponse channel response object
type ChannelResponse struct {
	Data *[]Channel `json:"data"`
	Size int        `json:"size"`
}

// LogEntriesResponse response object
type LogEntriesResponse struct {
	Elements   *[]LogEntry `json:"elements"`
	Size       int64       `json:"size"`
	Page       int64       `json:"page"`
	PageCount  int64       `json:"page_count"`
	EntryCount int64       `json:"entry_count"`
}

// Error api return object
type Error struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// initRouter initialize routing information
func initRouter() *gin.Engine {
	r := gin.New()
	r.RedirectTrailingSlash = true

	r.Use(gin.Logger())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	if appConf.ProxyURL != "" {
		err := r.SetTrustedProxies([]string{appConf.ProxyURL})
		if err != nil {
			log.Fatal(fmt.Sprintf("Problem to trust proxy url '%s'", appConf.ProxyURL))
			return nil
		}
	} else {
		// trust no proxies - except the user specified one
		err := r.SetTrustedProxies(nil)
		if err != nil {
			log.Fatal(fmt.Sprint("Problem to trust nil proxy url"))
			return nil
		}
	}
	r.Use(gin.Recovery())

	box := rice.MustFindBox("spec")
	if box == nil {
		log.Fatal("Error retrieving specs rice box")
	}
	r.StaticFS("/spec", box.HTTPBox())

	// define group
	apiPrefix := "/api/v2"
	apiV2 := r.Group(apiPrefix)
	apiV2.Use(func(context *gin.Context) {
		if len(GetAppConf().AccessControlAllowOrigin) > 0 {
			context.Header("Access-Control-Allow-Origin", fmt.Sprintf("%s", GetAppConf().AccessControlAllowOrigin))
		} else {
			log.Println("Warning! Using insecure default value '*' for 'Access-Control-Allow-Origin' (CORS) header. Please specify a value 'AccessControlAllowOrigin' in .oerc.yaml.")
			context.Header("Access-Control-Allow-Origin", "*")
		}
		// enable 10 h caching for /program, /log and /channel endpoints for browsers by default
		if strings.HasPrefix(context.FullPath(), apiPrefix+"/program") ||
			context.FullPath() == apiPrefix+"/log" ||
			strings.HasPrefix(context.FullPath(), apiPrefix+"/channel") {
			context.Header("Cache-Control", "public, max-age=36000")
		}
	})

	// default redirect to client app and ping redirect
	r.GET("/", func(context *gin.Context) {
		context.Redirect(301, "/client")
	})
	r.GET("/ping", func(context *gin.Context) {
		context.Redirect(301, "/api/v2/ping")
	})

	// ping
	apiV2.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "Pong", "date": time.Now()})
	})
	apiV2.GET("/status", getStatusHandler)
	// channel data
	apiV2.GET("/channels", getChannelsHandler)
	apiV2.GET("/channel/:channel_id", getSingleChannelHandler)
	// program data
	apiV2.GET("/program/yesterday", getProgramYesterdayHandler)
	apiV2.GET("/program/daily", getProgramTodayHandler)
	apiV2.GET("/program/tomorrow", getProgramTomorrowHandler)
	apiV2.GET("/program/daily/:channel_id", getDailyProgramWithChannelHandler)
	apiV2.GET("/program/tomorrow/:channel_id", getTomorrowProgramWithChannelHandler)
	apiV2.GET("/program/yesterday/:channel_id", getYesterdayProgramWithChannelHandler)
	apiV2.GET("/program", getProgramHandler)
	apiV2.GET("/program/entry/:id", getSingleProgramEntryHandler)
	// log data
	apiV2.GET("/log", getLogEntriesHandler)
	apiV2.GET("/log/entry/:id", getSingleLogEntriesHandler)
	apiV2.DELETE("/log/entry/:id", deleteSingleLogEntriesHandler)
	apiV2.DELETE("/log/clear", clearAllLogEntriesHandler)
	// recommendations
	apiV2.GET("/recommendations", getRecommendationsHandler)
	// search
	apiV2.GET("/search", getSearchHandler)
	apiV2.GET("/xmltv", getXmlTvHandler)

	if GetAppConf().ClientEnabled {
		clientBox := rice.MustFindBox("client/dist/client").HTTPBox()
		r.StaticFS("/client", clientBox)
	}

	return r
}
