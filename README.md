# OER-Collector

[![Build](https://travis-ci.org/emschu/oer-collector.svg?branch=master)](https://travis-ci.org/emschu/oer-collector)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=org.emschu.oer%3Aoer-collector-parent&metric=alert_status)](https://sonarcloud.io/dashboard?id=org.emschu.oer%3Aoer-collector-parent)
[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=org.emschu.oer%3Aoer-collector-parent&metric=ncloc)](https://sonarcloud.io/dashboard?id=org.emschu.oer%3Aoer-collector-parent)

A java software bundle to store and expose program data of public-law tv channels in Germany, Austria and Switzerland via a REST JSON interface.

This project is licensed under *AGPL v3* and you are encouraged to participate and improve functionality.

The focus of this project lies in providing reliable program data - ready to be enriched or analysed as long as there is no Open Data policy of the public-law sector.

At the moment its not intended to create links between program data and media(thek) data.
If you are looking for this have a look at [similar projects](#similar-projects).

**NOTE:** This server software is not ready to be used in any kind of *production* server environment.
Only use it locally and/or in protected environments.

**Current version:** v0.9.5

# Description

This software contains two components which are relevant to the end-user:

- **Collector:** Collects data and is running in a background process.
- **Server:** JSON REST interface. Have a look at the [OpenApi 2/Swagger specification of this server](./docs/openapi2-oer-server-web.json).

*Collector* and *Server* component run as independent processes so *Collector* problems won't affect the *Server* and vice versa.

You can run the *Collector* by triggering it manually, by crontab in background or by programmatic (java) cron.

In addition there is a very simple **configuration sample** to integrate the data model into [Apache Solr](http://lucene.apache.org/solr/).

To get an impression of the data there is an **Angular example frontend**.

Note: The public-law web pages this software needs to access are restricted to certain geographic IP regions. 

## Requirements
- JRE 8/10/11
- Collector: ~ 164 MB of RAM
- One of the following database backends: MariaDB/MySQL/Postgres 

# Setup 
- Download release package [here](https://github.com/emschu/oer-collector/releases)
- Extract package
- Create an empty database and setup `config.json` JDBC connection parameters in `oer.core.jdbc_*`. MariaDB is recommended. See [database options](#database_support) for different driver options.
- Database is initially set up by running the **Collector** *or* **Server** component once. Use the `--foreground` flag to see if everything is running fine.


## Configuration options

The following preferences are important to understand your possibilities to control this software.
Configuration of Collector and Server is done in the central `config.json` file.

- **oer.core.jdbc_**\*: Database and Hibernate configuration and connection parameters.
- **oer.collector.update_mode**: Flag to skip updating of \"old\" (as defined in hours by *oer.collector.invalidate_update_hours*) entries in general.
- **oer.collector.enable_tv_show_collect**: Flag to control collecting tv show data.
- **oer.collector.enable_program_entry_collect**: Flag to control collecting program entry data.
- **oer.collector.proxy_host**: (Optional) HTTP proxy host.
- **oer.collector.proxy_port**: (Optional) HTTP proxy port.
- **oer.collector.skip_ard**: Flag to skip collecting *ARD* data.
- **oer.collector.skip_zdf**: Flag to skip collecting *ZDF* data.
- **oer.collector.skip_orf**: Flag to skip collecting *ORF* data.
- **oer.collector.skip_srf**: Flag to skip collecting *SRF* data.
- **oer.collector.cron_definition**: Collector is running in endless cron mode. Not recommended at the moment. Should never be used in combination with mass mode. Example for twice executions a day at 4 AM and 4 PM: `0 0 4,16 * * *`. Default: `null`
- **oer.collector.cron_mode_run_at_startup**: Runs collection process at program startup. Effective if *cron_definition* is not empty. Default: `true`
- **oer.collector.invalidate_update_hours**: Number of hours to consider data records in db as old and web data is used to refresh/update record. Has no effect if not in update mode. Default: `72` 
- **server.address**: Host the webserver will listen to. Default: `127.0.0.1`
- **server.port**: Port the webserver will listen to. Default: `8081`


### Default collection mode
- **oer.collector.collect_past_program_days_max**: Number of days in past to collect data for.
- **oer.collector.collect_future_program_days_max**: Number of days in future to collect data for.

### Mass mode
- **oer.collector.mass_mode**: Flag to enable collecting date ranges. Should be used with caution and should be monitored.
- **oer.collector.start_date**: Start of date range to collect data for. Only active if `mass_mode: true`.
- **oer.collector.end_date**: End of date range to collet data for. Only active if `mass_mode: true`.


# Run and use

By default the Collector and Server component are running as background processes, until they stop due to configuration (e.g. if no Java *cron* definition is used) or you stop them.

Stopping a component should never be problematic, except during database migrations at (first time) startup.

To debug your (database) configuration, use the `--foreground` flag of the following start scripts.

```bash
$ ./start.sh
```
This script basically calls both `start_collector.sh` and `start_server.sh`.

NOTE: For complete ORF program data, you should run the collector at least twice.

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
Located in `oer-example-client`. Basically it contains a timeline with all channels and program entry data on it. If you want to have a look, use the following commands:
```bash
$ npm install
$ npm run start
```
After the second command the simple example client ui should be reachable at [http://localhost:4200](http://localhost:4200).  
Feel free to improve and enhance!

## Apache Solr-Integration

*This is a work in progress.*

There is a basic sample how to setup and use Apache Solr with MariaDB and the collected program entry data.

Have a look at `tools/solr` directory and customise. `tools/solr/setup_with_mariadb.sh` and `data-config.xml` for your needs. Please share improvements. 

## Container environment

OER Collector + Server can be used with *docker-compose* by following this [guide](https://github.com/emschu/oer-collector/tree/master/tools/docker). A `docker-compose.yml` is contained in release packages and can be found in `tools/docker` directory (of this repository).  

Images: MariaDB, Squid Proxy, OER Collector, OER Server.

# Channel list

Note: The first column does not necessarily have to correspond to the channel id in the database.

| No. | Channel    | Version |
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
|22| ORF eins         |   v1 |
|23| ORF 2            |   v1 |
|24| ORF III          |   v1 |
|25| ORF Sport +      |   v1 |
|26| SRF-1            |   v1 |
|27| SRF-zwei         |   v1 |
|28| SRF-info         |   v1 |


## Data import limits

| Channel family | Earliest date       | Latest date     |
| ---------------| ------------------- | --------------- |
| ARD/ZDF        | ~ 2011              | Today + 6 weeks |
| ORF            | Today - 14 days     | Today + 22 days |
| SRF            | Today - 14 days     | Today + 29 days |

# Project guidelines
- This project is non-commercial.
- Private/commercial sector tv or radio stations will *never* be part of this project.
- This project shall be an instrument mainly to analyse the program and constructively 
improve public-law tv and radio stations.
- This project would be superfluous, if there was a public API for public data, OpenData...
- Minimise traffic and external load to the least needed.

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

**Planned - and not yet implemented - support of tv/radio channels and SRF support:** 

*Radio:*
 - DLR/DLF
 - DW
 - ...

# More information
- [German] [Die Vermessung des TV-Programms auf datenjournalist.de](https://www.datenjournalist.de/die-vermessung-des-tv-programms/)