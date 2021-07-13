# OERC

... is a short name for **OER-Collector**, which is a software project to store, view and search the program data of public-law ("öffentlich-rechtliche") TV stations in Germany, Austria and Switzerland.

This tool needs an external PostgreSQL database and some configuration parameters, then you can:
- **Collecting TV program data** of 28 channels (`oerc fetch`)
- **Search for interesting program** items by looking for your own keywords (`oerc search`)
- Running an **HTTP backend server** to access program data in JSON format (`oerc server`)
- Running an HTTP server to serve a **client web application** to view the program data and 
  your personal recommendations (`oerc server`)

With the help of `oerc` you can build and use your own, private, TV program recommendation tool while ALL
information is processed locally.

This project is written in Go and it is *AGPL v3* licensed. You are encouraged to participate and improve functionality. 
[Just create an issue!](https://github.com/emschu/oerc/issues)

The focus of this project lies in providing program data for individuals - ready to be enriched or analysed as long as there is no Open Data policy of the public-law sector.

At the moment its not intended to create links between program data and media(thek) data.
If you are looking for this have a look at [similar projects](#similar-projects).

*Note 1:* This server and client software is not ready to be used directly in the internet without further changes.
It is recommended to use it locally only or in protected environments and don't expose it to the internet.

*Note 2*: Unfortunately the public-law web pages this software needs to access are restricted to certain geographic IP regions in general.

**Current software quality: Beta**

![oer-collector logo](./docs/logo.png)

# Install

1. Get this application
```shell
go get -u github.com/emschu/oerc
```
**OR** download the latest release for your platform from the [GitHub release page](https://github.com/emschu/oerc/releases).

2. Set up a PostgreSQL database and run it
3. Run `oerc init`.
This will copy a sample configuration file to the path `~/.oerc.yaml` - if the file does not exist already.
   You have to change some of the values, at least you have to replace `<db_name>`, `<db_user>`, `<db_password>` in the 
   configuration file to reach the database.
4. Run `oerc fetch` for the first time and wait until the first program data is collected for you.

# Description

```shell
NAME:
   oerc - Command line tool to manage the oerc application

USAGE:
   oerc [global options] command [command options] [arguments...]

VERSION:
   1.0.0, License: AGPLv3, https://github.com/emschu/oer-collector

DESCRIPTION:
   Fetch, view and search TV program data of public-law stations in Germany, Switzerland and Austria

COMMANDS:
   clear       Clear the database. Be careful!
   fetch, f    Get latest data
   init, i     Initial database and configuration setup check
   search, sc  Search program data and create recommendations
   server, sv  Start API HTTP backend server
   status, s   show status
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value, -c value  Path to the yaml configuration file (default: ~/.oerc.yaml)
   --verbose                 Verbose log output (default: false)
   --help, -h                show help (default: false)
   --version, -v             print the version (default: false)
```

## Configuration options

The following preferences are important to understand your possibilities to control this software.
You can find this file [here](./config/.oerc_default.dist.yaml) and if you run `oerc init` this file will be created at
`~/.oerc.yaml` for you. You *must* provide valid postgres database connection details.

If you don't want to put your configuration at the user's home directory, you can also use the 
`-c <path-to-your-oerc>.yaml` 
argument for all `oerc` commands.


```yaml
# general fetch settings
ForceUpdate: false
# don't update entries which were processed already in the last 6 hours
TimeToRefreshInMinutes: 360
# if you run "oerc fetch" the last 2 and the next 7 days of program data will be fetched
DaysInPast: 2
DaysInFuture: 7
# enable/disable channel families
EnableARD: true
EnableZDF: true
EnableORF: true
EnableSRF: true
EnableProgramEntryCollection: true
EnableTVShowCollection: true
# provide a full URL to a HTTP(S) or SOCKS proxy or it will fail
ProxyUrl:
TimeZone: Europe/Berlin
# backend server settings
ServerHost: 127.0.0.1
ServerPort: 8080
# client server settings
ClientEnabled: true
# db settings
DbType: postgres
DbHost: localhost
DbPort: 5432
DbName: <db_name>
DbUser: <db_user>
DbPassword: <db_password>
DbSSLEnabled: false
# search settings
# only search for recommendations in the next 4 days
SearchDaysInFuture: 4
# these are example values. Feel free to create you own list of keywords :)
SearchKeywords:
  - Stromberg
  - Die Anstalt
  - Max Uthoff
  - Claus von Wagner
  - Loriot
  - Zapp
  - Kroymann
  - James Bond
  - Satire
# these channels won't be recognized during the "search" for recommendations based on your keywords
SearchSkipChannels:
  - KIKA
  - ORF Sport +
```

# Run and use

After installing `oerc` and setting it up, you should run at least one time the `oerc fetch` command.

It is recommended to update the program data regularly, e.g. daily, by using a **cron job** which runs `oerc fetch` 
and `oerc search`.

While it is possible to run `oerc server` in a user session, you should consider to create a systemd
service to run and control the web server (backend + frontend) in the background persistently. 

The following two systemd service files
are simple examples for you to integrate `oerc` with systemd:

**oerc.service:**
```
[Unit]
Description=oer-collector service
After=network.target

[Service]
Type=simple
ExecStart=<path_to_oerc_bin> server
StandardOutput=journal
KillMode=process

[Install]
WantedBy=multi-user.target
```

*Note*: You *must* replace `<path_to_oerc_bin>` in with a correct path to the `oerc` binary. If you don't know how to do this
type `which oerc`.

Copy the modified system service templates to your systemd services directory, e.g. `/etc/systemd/system` and reload the systemd daemon by executing

`$ sudo systemctl daemon-reload`. 

After the last command you can use 

`$ sudo systemctl [start|stop] oerc` 

to start (or stop) the services. 
If you want to get the servers up after system (re-)boot, you need to execute 

`$ sudo systemctl [enable|disable] oerc`.

If you do so (enabling both services by default), please keep in mind that the PostgreSQL database needs 
to be available, too, so use systemctl to enable the postgres service as well.

## Containerized run
Will be added in future releases.


## Channel list

Note: The first column does not necessarily have to correspond to the channel id in the database.

| No. | Channel    | Version |
| --- | --------------| ---- |
|1| ARD               |   v2 |
|2| ZDF               |   v2 |
|3| 3Sat              |   v2 |
|4| ARTE              |   v2 |
|5| ZDFInfo           |   v2 |
|6| ZDFNeo            |   v2 |
|7| Phoenix           |   v2 |
|8| KiKa              |   v2 |
|9| ARD One           |   v2 |
|10| Tagesschau24     |   v2 |
|11| ARD Alpha        |   v2 |
|12| SWR RP Fernsehen |   v2 |
|13| WDR Fernsehen	  |   v2 |
|14| SWR BW Fernsehen |   v2 |
|15| SR Fernsehen     |   v2 |
|16| Radio Bremen TV  |   v2 |
|17| RBB Fernsehen    |   v2 |
|18| NDR Fernsehen 	  |   v2 |
|19| MDR Fernsehen	  |   v2 |
|20| HR Fernsehen     |   v2 |
|21| BR Fernsehen     |   v2 |
|22| ORF eins         |   v2 |
|23| ORF 2            |   v2 |
|24| ORF III          |   v2 |
|25| ORF Sport +      |   v2 |
|26| SRF-1            |   v2 |
|27| SRF-zwei         |   v2 |
|28| SRF-info         |   v2 |

## Data import limits

| Channel family | Earliest date       | Latest date     |
| ---------------| ------------------- | --------------- |
| ARD/ZDF        | ~ 2011              | Today + 6 weeks |
| ORF            | Today - 14 days     | Today + 22 days |
| SRF            | Today - 14 days     | Today + 29 days |

# Project guidelines
- This project is non-commercial.
- Private/commercial sector TV or radio stations will *never* be part of this project.
- This project shall be an instrument mainly to analyze the program and constructively
  improve public-law TV and radio stations. Or just use it privately.
- This project would be superfluous, if there was a public API for public data, OpenData...
- Minimise traffic and external load to the least needed.
- Avoid security problems on the client side and maintain privacy of the users.
- All parts of the software should work on "low-resource" platforms, e.g. a Raspberry Pi 3b+

<a name="similar-projects"></a>
## Similar projects:
- [cemrich/zapp-backend](https://github.com/mediathekview/zapp-backend)
- [MediathekView(Web)-Project](https://github.com/mediathekview)
- [MediathekDirekt](https://mediathekdirekt.de/) + [Sources](https://gitlab.com/mediathekdirekt/mediathekdirekt)
- [EPG Scraper for ARD TV Stations to Use With tvheadend External XMLTV Grabber](https://projects.webvoss.de/2019/04/14/legal-epg-scraper-for-ard-tv-stations-to-use-with-tvheadend-external-xmltv-grabber/)

## OpenApi 3 specification
If you run the (backend) server (just run `oerc server`) an OpenApi 3 specification is shipped at
`/spec/openapi3.json`, respectively `/spec/openapi3.yaml`. Or - alternatively - have a look at the spec 
files in [this directory](./docs).

Please notify the maintainer of this project if you build something around the JSON HTTP API oerc
offers (for the mail address see below in `License` section).

# License

This project is licensed under [GNU Affero General Public License](./LICENSE).

```text
oerc, alias oer-collector
Copyright (C) 2021 emschu[aet]mailbox.org

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
```

## Used libraries

- [Urfave/cli v2](https://github.com/urfave/cli/) – CLI is a simple, fast, and fun package for building command line apps in Go.
- [Gorm](https://gorm.io/) – The fantastic ORM library for Golang
- [Colly](http://go-colly.org/) – Fast and Elegant Scraping Framework for Gophers
- [Gin](https://github.com/gin-gonic/gin) – Gin is a HTTP web framework
- [Bluemonday](https://github.com/microcosm-cc/bluemonday) – A fast golang HTML sanitizer

# Development
This project is shipped with a `Makefile` to ease the development and testing process. Be sure to run `make build` before
filing a pull request.

## Database
### Development postgres container
```console
# docker run --name oer-postgres -p 5432:5432 -e POSTGRES_PASSWORD=root -e POSTGRES_DB=oer_server_dev -d postgres:12.3

// test connection with
$ psql -U postgres -h 127.0.0.1 -W 
```

## Contributing
- File issues on GitHub to request and discuss new features or bugs there.
- You need a new feature/improvement? -> File an issue.
- Contribute code through pull requests or submit patch files.

# More information
- [German] [Die Vermessung des TV-Programms auf datenjournalist.de](https://www.datenjournalist.de/die-vermessung-des-tv-programms/)
