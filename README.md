# OER-Collector

[![Build](https://travis-ci.org/emschu/oer-collector.svg?branch=master)](https://travis-ci.org/emschu/oer-collector)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=org.emschu.oer%3Aoer-collector-parent&metric=alert_status)](https://sonarcloud.io/dashboard?id=org.emschu.oer%3Aoer-collector-parent)
[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=org.emschu.oer%3Aoer-collector-parent&metric=ncloc)](https://sonarcloud.io/dashboard?id=org.emschu.oer%3Aoer-collector-parent)

A java software bundle to store and expose program data of public-law tv channels in Germany via a REST JSON interface.

This project is licensed under *AGPL v3* and you are encouraged to participate and improve functionality.

The focus of this project lies in providing reliable program data - ready to be enriched or analysed as long as there is no Open Data policy of the public-law sector in Germany.

At the moment its not intended to create links between program data and media(thek) data.
If you are looking for this have a look at [similar projects](#similar-projects).

**NOTE:** This server software is not ready to be used in any kind of *production* server environment.
Only use it locally and/or in protected environments.

**Current version:** v0.9.4

# Description

This software contains two end-user relevant components:

- **Collector:** Collects data and is running in a background process.
- **Server:** JSON REST interface. Have a look at the [OpenApi 2/Swagger specification of this server](./docs/openapi2-oer-server-web.json).

*Collector* and *Server* component run as independent processes so *Collector* problems won't affect the *Server* and vice versa.
You can run the *Collector* by triggering it manually, by cron in background or by programmatic (java) cron.

In addition there is a very simple **configuration sample** to integrate the data model into [Apache Solr](http://lucene.apache.org/solr/).

To get an impression of the data there is an **Angular example frontend**.

Note: The public-law web pages this software needs to access are restricted to certain geographic IP regions. 

## Requirements
- JRE 8/10/11
- Collector: ~ 164 MB of RAM
- One of the following database backends: MariaDB/MySQL/Postgres 

# Setup 
- Download release package [here](https://github.com/emschu/oer-collector/releases).
- Extract package
- Create an empty database and setup `config.json` JDBC connection parameters in `oer.core.jdbc_*`. MariaDb is recommended. See [database options](#database_support) for different driver options.
- Database is initially set up by running the **Collector** *or* **Server** component once. Use the `--foreground` flag to see if everything is running fine.


## Configuration options

The following preferences are important to understand your possibilities to control this software.
Configuration of Collector and Server is done in the central `config.json` file.

- **oer.core.jdbc_**\*: Database and Hibernate configuration and connection parameters.
- **oer.collector.update_mode**: Flag to skip updating of \"old\" (as defined in hours by *oer.collector.invalidate_update_hours*) entries in general.
- **oer.collector.enable_tv_show_collect**: Flag to control collection of tv shows.
- **oer.collector.enable_program_entry_collect**: Flag to control collection of program entries.
- **oer.collector.proxy_host**: Optional HTTP proxy host.
- **oer.collector.proxy_port**: Optional HTTP proxy port.
- **oer.collector.skip_ard**: Flag to skip collection of ARD data.
- **oer.collector.skip_zdf**: Flag to skip collection of ZDF data.
- **oer.collector.cron_definition**: Collector is running in endless cron mode. Not recommended at the moment. Should never be used in combination with mass mode. Example: `*/15 * * * * *`. Default: `null`
- **oer.collector.invalidate_update_hours**: Number of hours to consider data records in db as old and web data is used to refresh/update record. Has no effect if not in update mode. Default: `72` 

### Default collection mode
- **oer.collector.collect_past_program_days_max**: Number of days in past to collect data for.
- **oer.collector.collect_future_program_days_max**: Number of days in future to collect data for.

### Mass mode
- **oer.collector.mass_mode**: Flag to enable collecting date ranges. Should be used with caution and should be monitored.
- **oer.collector.start_date**: Start of date range to collect data for. Only active if `mass_mode: true`.
- **oer.collector.end_date**: End of date range to collet data for. Only active if `mass_mode: true`.


# Run

By default the Collector and Server component are running as background processes, until they stop due to configuration (e.g. if no Java *cron* definition is used) or you stop them.

Stopping a component should never be problematic, except during database migrations at (first time) startup.

To debug your (database) configuration, use the `--foreground` flag of the following start scripts.

```bash
$ ./start.sh
```
This script basically calls both `start_collector.sh` and `start_server.sh`.

## Manage Collector (only)
```bash
$ ./start_collector.sh [--foreground]
# Follow log output
$ tail -f collector.log
$ ./stop_collector.sh
```
## Manage Server (only)
```bash
$ ./start_server.sh [--foreground]
$ ./stop_server.sh
```

## Stop services
```bash
$ ./stop.sh
```

## Angular example frontend
Located in `oer-example-client`. To use it, use the following commands:
```bash
$ npm install
$ npm run start
```
After the simple example client ui is reachable at [http://localhost:4200](http://localhost:4200).  
Feel free to improve and enhance!

## Apache Solr-Integration

*Work in progress.*  
Have a look at `tools/solr/` dir of this project.

# Channel list
| Nr. | Sender     | Version |
| --- | --------------| ---- |
|1| ARD               |   v1 |
|2| ZDF               |   v1 |
|3| 3Sat              |   v1 |
|4| ARTE              |   v1 |
|5| ZDFInfo           |   v1 |
|6| ZDFNeo            |   v1 |
|7| Phoenix           |   v1 |
|8| KiKa              |   v1 |
|9| ARD One           |   v1 |
|10| Tagesschau24     |   v1 |
|11| ARD Alpha        |   v1 |
|12| SWR RP Fernsehen |   v1 |
|13| WDR Fernsehen	  |   v1 |
|14| SWR BW Fernsehen |   v1 |
|15| SR Fernsehen     |   v1 |
|16| Radio Bremen TV  |   v1 |
|17| RBB Fernsehen    |   v1 |
|18| NDR Fernsehen 	  |   v1 |
|19| MDR Fernsehen	  |   v1 |
|20| HR Fernsehen     |   v1 |
|21| BR Fernsehen     |   v1 |

# Project guidelines
- This project is non-commercial.
- Private/commercial sector tv or radio stations will *never* be part of this project.
- This project shall be an instrument mainly to analyze the program and constructively 
improve public-law tv and radio stations.
- This project would be superfluous, if there was a public API for public data, OpenData...
- Minimize traffic and external load to the least needed.

<a name="similar-projects"></a>
## Similar projects:
- [cemrich/zapp-backend](https://github.com/cemrich/zapp-backend)
- [MediathekView(Web)-Project](https://github.com/mediathekview)
- [MediathekDirekt](https://mediathekdirekt.de/) + [Sources](https://gitlab.com/mediathekdirekt/mediathekdirekt)


# REST API
Run the server and visit the following site to see an OpenApi v2 REST interface specification *oer-collector* provides:

[OpenApi v2 API description of this project](http://127.0.0.1:8081/openapi2)

# Data Model
### Entities
- ProgramEntry
- Image
- Tag
- TvShow
- Artist
- Channel

### ER-Diagram

![ER-Diagram](https://raw.githubusercontent.com/emschu/oer-collector/master/docs/er-model-v1.png)

Java source classes can be found in *org.emschu.oer.core.model*.

# <a name="database_support"></a>
# Database support
Example configuration snippets for different database backends.


### MariaDB 5.5+

```json
{
  "oer.core.jdbc_datasource": "jdbc:mariadb://localhost:3306/oer_server?useUnicode=yes&characterEncoding=UTF-8",
  "oer.core.jdbc_datasource_user": "oer_server",
  "oer.core.jdbc_datasource_password": "oer_server_pw",
  "oer.core.jdbc_hibernate_dialect": "org.hibernate.dialect.MariaDBDialect",
  "oer.core.jdbc_hibernate_tmp_metadata_defaults": false
}
```

### Oracle MySQL 5.5+
```json
{
  "oer.core.jdbc_datasource": "jdbc:mysql://localhost:3306/oer_server?useUnicode=yes&characterEncoding=UTF-8",
  "oer.core.jdbc_datasource_user": "oer_server",
  "oer.core.jdbc_datasource_password": "oer_server_pw",
  "oer.core.jdbc_hibernate_dialect": "org.hibernate.dialect.MySQLDialect",
  "oer.core.jdbc_hibernate_tmp_metadata_defaults": true
}
```

### PostgreSQL 9.6/10.7
```json
{
  "oer.core.jdbc_datasource": "jdbc:postgresql://127.0.0.1:5432/oer_server",
  "oer.core.jdbc_datasource_user": "postgres",
  "oer.core.jdbc_datasource_password": "postgrespw",
  "oer.core.jdbc_hibernate_dialect": "org.hibernate.dialect.PostgreSQL9Dialect",
  "oer.core.jdbc_hibernate_tmp_metadata_defaults": false
}
```

# License

This project is licensed under [GNU Affero General Public License](./LICENSE).

```text
    oer-collector
    Copyright (C) 2019 emschu@mailbox.org

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
- [Spring (Boot) Framework](https://spring.io/)
- [Spring data + Hibernate](http://projects.spring.io/spring-data/)
- [Jsoup](https://jsoup.org/)
- [Gson](https://github.com/google/gson)
- [Springfox](http://springfox.github.io/springfox/)
- [iCal4J](https://github.com/ical4j/ical4j)

# Contributing
- File issues on github to request and discuss new features/bugs there.
- You need a new feature/improvement? -> File an issue.
- Contribute code through pull requests or submit patch files.

**Planned - and not yet implemented - support of tv/radio channels and SRF + ORF support:** 

*Radio:*
 - DLR/DLF
 - DW
 - ...

*TV:*
 - SRF-Channel-Family
 - ORF-Channel-Family

Support for container environments is under construction.


# More information
- [German] [Die Vermessung des TV-Programms auf datenjournalist.de](https://www.datenjournalist.de/die-vermessung-des-tv-programms/)