# OERC

... is a software project to locally store, view and search the program data of
public-law ("öffentlich-rechtliche") TV stations in Germany, Austria and Switzerland.

While `oerc` is a command-line tool only, a built-in web application is provided 
at `127.0.0.1:8080` (*default*) if you run `oerc server`.

### Commands

- **Collecting TV program data** of 27 TV channels (`oerc fetch`)
- **Search for interesting program** items by looking for your own keywords (`oerc search`)
- Running an **HTTP server** to access program data in a simple web UI (`oerc server`)
- Find overlaps in data or program changes (`oerc overlap-check`)
- By default, the server contains a **client web application** for your browser to view the program data and
  your personal program recommendations
- Support of XMLTV via CLI and HTTP API

Using `oerc`, you can build and use your own private TV program recommendation tool while ALL
information is processed and kept locally.

This project is written in Golang, and it is *AGPL v3* licensed. You are encouraged to participate and improve
functionality.
[Just create a GitHub issue!](https://github.com/emschu/oerc/issues)

The focus of this project lies in providing program data for individuals - ready to be enriched, searched or analyzed
as long as there is no (real) Open Data policy of the public-law sector.

At the moment it's not intended to create links between program data and Media(thek) information.
If you are looking for this, have a look at [similar projects](#similar-projects).

*Note 1:* This server and client software is not ready to be used directly on the internet without further changes.
It is recommended to use it locally only or in protected environments and not to expose it to the internet.

*Note 2*: Unfortunately, the public-law web pages this software needs to access are restricted to certain geographic IP
regions in general.

**Current software quality: Beta**

![oer-collector logo](./docs/logo.png)

### Web Client

![oerc web client screenshot](./docs/screenshot_web.png)
![oerc web client screenshot with overlaps](./docs/screenshot_web_overlaps.png)
![oerc web client screenshot with recommendations](./docs/screenshot_web_recommendations.png)

# Setup & Run

```shell
go install github.com/emschu/oerc@latest
```

## Download and start pre-built executable

Download the latest `oerc` binary from [here](https://github.com/emschu/oerc/releases) for your platform, make it executable, and you are ready to start. Builds are available for all popular platforms.

## docker/docker-compose
Clone this repository, then:

```bash
$ docker-compose build
$ docker-compose up
```

## Manual Installation

1. Install the application with Go using the following command: `go install github.com/emschu/oerc@latest`

2. Setup a database.
2.1 Option 1: SQLite database (default) - this is the easiest way to get started.
    - The SQLite database file is located at `~/.oerc.db` by default.
2.2 Option 2: Set up a PostgreSQL database (12+), configure a database with a user and start it.
    - The [development section](./DEVELOPMENT.md) contains a simple docker command to run a local database for development or testing. 
3. Run `oerc init`.   
   - This will copy a sample configuration file to the path `~/.oerc.yaml` if the file does not exist already.
   - If you use SQLite (default), you are ready to go.
   - If you use PostgreSQL, you have to change the values for `<db_name>`, `<db_user>`, and `<db_password>` in the configuration file to reach the database you've configured.
4. Run `oerc fetch` for the first time and wait until the first program data is collected.
5. Run `oerc server` to have a browser application at `http://localhost:8080/client` (the endpoint is configurable).
6. Run `oerc search` to look for recommendations based on your self-defined keywords in the configuration.

## Description

```text
NAME:
   oerc - Command line tool to manage the oerc application

USAGE:
   oerc [global options] command [command options]

VERSION:
   0.21.0, License: AGPLv3, https://github.com/emschu/oerc

DESCRIPTION:
   Fetch, view and search TV program data of public-law stations in Germany, Switzerland and Austria

COMMANDS:
   clear               Clear the database. Be careful!
   fetch, f            Get latest data
   fetch-range, fr     Fetch a specific date range
   full-overlap-check  Run overlap check on all program entries. Could take very long.
   init, i             Initial database and configuration setup check
   overlap-check       Run overlap check on currently fetched time range
   search, sc          Search program data and create recommendations
   server, sv          Start webserver with oerc API and an embedded browser client
   status, s           Show some of the app's status information
   xmltv               Export the database to xmltv format
   help, h             Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value, -c value  Path to the yaml configuration file (default: ~/.oerc.yaml)
   --verbose                 Verbose log output (default: false)
   --help, -h                show help
   --version, -v             print the version                                                           
```

### Configuration options

The following preferences are important to understand your possibilities to control this software.
You can find the default configuration file [here](./config/.oerc_default.dist.yaml). Running `oerc init` will create this file at `~/.oerc.yaml`.

If you don't want to put your configuration in the user's home directory, you can also use the
`-c <path-to-your-oerc>.yaml`
argument for all `oerc` commands.

```yaml
# oerc configuration file
# db settings (sqlite is default)
DbType: sqlite
DbHost: ~/.oerc.db
DbPort:
DbName: oerc_db
DbSchema: public
DbUser:
DbPassword:
DbSSLEnabled: false
ForceUpdate: false
TimeToRefreshInMinutes: 360
DaysInPast: 2
DaysInFuture: 7
EnableARD: true
EnableZDF: true
EnableORF: true
EnableSRF: true
EnableProgramEntryCollection: true
EnableTVShowCollection: true
ProxyUrl:
TimeZone: Europe/Berlin
# backend server settings
ServerHost: 127.0.0.1
ServerPort: 8080
# client server settings
ClientEnabled: true
# search settings
SearchDaysInFuture: 5
# feel free to compose the following list of keywords as you want. They are used case-insensitive.
SearchKeywords:
  - Loriot
SearchSkipChannels:
  - KIKA
  - ORF Sport +
# if you don't know what this means, you should probably keep this value
AccessControlAllowOrigin: "http://127.0.0.1:8080"
```

Additional configuration options that may be available:
- `Debug: true/false` - Enable debug logging
- `ProfilingEnabled: true/false` - Enable performance profiling

## Usage

After installing `oerc` and setting it up, you should run the `oerc fetch` command at least once.

It is recommended to update the program data regularly, e.g. daily, by using a (Unix) **cron job** which
runs `oerc fetch`
and `oerc search`.

While it is possible to run `oerc server` in a user session, you should consider creating a systemd
service to run and control the web server (backend + frontend) in the background persistently.

The following systemd service file is a simple example to integrate `oerc` with systemd on a Linux system:

**oerc.service:**

```
[Unit]
Description=oerc service
After=network.target

[Service]
Type=simple
ExecStart=<path_to_oerc_bin> server
StandardOutput=journal
KillMode=process

[Install]
WantedBy=multi-user.target
```

*Note*: You *must* replace `<path_to_oerc_bin>` with an absolute path to the `oerc` binary (`which oerc`).

Copy the modified system service templates to your systemd services directory, e.g. `/etc/systemd/system`, and reload the
systemd daemon by executing

`$ sudo systemctl daemon-reload`.

After the last command you can use

`$ sudo systemctl [start|stop] oerc`

to start (or stop) the services.

To (not) start the service on boot, type

`$ sudo systemctl [enable|disable] oerc`.

If you enable the service by default, please keep in mind that the database needs
to be available too. If you use PostgreSQL, ensure to enable the PostgreSQL database service as well.

### Channel list

**Channels:** ARD, ZDF, 3Sat, ARTE, ZDFInfo, ZDFNeo, Phoenix, KiKa, ARD One, Tagesschau24, ARD Alpha, SWR, WDR, SR, Radio Bremen TV, RBB, NDR, MDR, HR, BR, ORF eins, ORF 2, ORF III, ORF Sport +, SRF-1, SRF-zwei, SRF-info.
### Time range import limits

| Channel family | Earliest date   | Latest date     |
|----------------|-----------------|-----------------|
| ARD            | Today - 8 days  | Today + 8 days  |
| ZDF            | 2015/03         | Today + 6 weeks |
| ORF            | Today - 14 days | Today + 22 days |
| SRF            | Today - 14 days | Today + 29 days |

## Project guidelines

- This project is non-commercial.
- Private/commercial sector TV or radio stations will *never* be part of this project.
- This project shall be an instrument mainly to analyze the program and constructively
  improve public-law TV and radio stations. Or just use it privately.
- This project would be superfluous if there was a public API for public data, OpenData...
- Minimize traffic and external load to the least needed.
- Avoid security problems on the client side and maintain privacy of the users.
- All parts of the software should work on "low-resource" platforms, e.g. a Raspberry Pi 3b+

<a name="similar-projects"></a>

### Similar projects:

- [mediathekview/zapp-backend](https://github.com/mediathekview/zapp-backend)
- [MediathekView(Web)-Project](https://github.com/mediathekview)
- [MediathekDirekt](https://mediathekdirekt.de/) + [Sources](https://gitlab.com/mediathekdirekt/mediathekdirekt)

### OpenAPI 3 specification

If you run the (backend) server (`oerc server`), an OpenAPI 3 specification is shipped at
`http://localhost:8080/spec/oerc-openapi3.json`. Or have a look at the
API specification files in [this directory](./docs).

## License

This project is licensed under [GNU Affero General Public License 3](./LICENSE).

```text
oerc, alias oer-collector
Copyright (C) 2021-2026 emschu[aet]mailbox.org

This program is free software: you can redistribute it and/or modify 
it under the terms of the GNU Affero General Public License as 
published by the Free Software Foundation, either version 3 of the 
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, 
but WITHOUT ANY WARRANTY; without even the implied warranty of 
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the 
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public 
License along with this program. 
If not, see <https://www.gnu.org/licenses/>.
```

### Used libraries

- [Urfave/cli v2](https://github.com/urfave/cli/) – CLI is a simple, fast, and fun package for building command line
  apps in Go.
- [Gorm](https://gorm.io/) – The fantastic ORM library for Golang
- [Colly](http://go-colly.org/) – Fast and Elegant Scraping Framework for Gophers
- [Gin](https://github.com/gin-gonic/gin) – Gin is an HTTP web framework
- [Bluemonday](https://github.com/microcosm-cc/bluemonday) – A fast golang HTML sanitizer

## Development

See [DEVELOPMENT.md](./DEVELOPMENT.md).

### Contributing

- File issues on GitHub to request and discuss new features or bugs there.
- You need a new feature/improvement? -> File an issue.
- Contribute code through pull requests or submit patch files.

## More information

- [German] [Die Vermessung des TV-Programms auf datenjournalist.de](https://www.datenjournalist.de/die-vermessung-des-tv-programms/)
