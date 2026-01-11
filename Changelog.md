# Changelog oerc

## 0.21.0 - 2026/XX
- Support for SQlite database backend
- Change default database to sqlite
- Improve the accuracy and performance of overlap check
- Add channels to the web UI, allow reordering of channels
- Dependency and library updates

## 0.20.0 - 2025/10/26
- Golang 1.24
- Dependency updates 
- Update OpenApi spec and remove conversion logic from project

## 0.19.0 - 2025/04/27
- Add new `xmltv` subcommand and create `xmltv` http api endpoint

## 0.18.0 - 2025/04/13
- Fix recent problems with zdf api
- update dependencies

## 0.17.1 - 2024/12/27
- Fix problem validating ard image urls

## 0.17.0 - 2024/12/26
- Golang 1.23, toolchain and dependency updates
- Changed fetching ARD program information using new ARD program and tv show API 
- Completely new ard channels (the former ones get deprecated during first startup)

## 0.16.0 - 2024/07/04
- Golang 1.21
- Dependency update

## 0.15.0 - 2024/02-04
- Fix ORF fetching
- update go and npm dependencies

## 0.14.0 - 2023/12/27
- go 1.20
- dependency update
- Add dynamic channel integration into UI, static sorting
- Channels can be deprecated now: Mechanism to add and deprecate channels in channel list
- Improve update handling of materialized view for fast status info

## 0.13.0 - 2023/08/12
- update dependencies
- Use new SRF JSON API instead of markup parsing, fixing fetch process for SRF

## 0.12.0 - 2023/04/22
- update dependencies
- client: ui/ux improvements

## 0.11.0 - 2023/01/03
- update toolchains and dependencies

## 0.10.2 - 2022/09/03
- Dependency updates of the application and the client
- client: add modal if there is no program data
- small ard parser improvements (mostly for historical data)
- refine parser log output for easier debugging
- Golang `1.18` is the current minimal version supported 

## 0.10.1 - 2022/06/25
- Add `DbSchema` configuration property tu support arbitrary postgres schema
- Web-Client: Minor improvements

## 0.9.18 - 2022/05/07
- API field change in StatusResponse: `problem_count` -> `log_count`
- docker-compose improvements and documentation added
- Dependency updates for web client
- Clean up configuration file

## 0.9.17 - 2022/03/13
- introduce `fetch-range` subcommand to import all available information
- More code improvements and more tests

## 0.9.16 - 2022/03/06
- Improving code structure in general which also brings more performance
- now `postgres` and `postgresql` as values for `DbType` will be accepted

## 0.9.15 - 2022/02/13
- new configuration option `AccessControlAllowOrigin` to control the CORS-Header of server directly
- minor improvements of collection process
- minor changes to program endpoint of HTTP api

## 0.9.14 - 2021/12/06
- add screenshots to repository
- add new action (`oerc overlap-check`) to run the overlap check separately from other commands
- improve documentation

## 0.9.13 - 2021/11/12
- bugfixes

## 0.9.12 - 2021/11/12
- New Feature: Auto-detection of overlapping program items is added to the `fetch` command and the web client.
- New Command: `oerc full-overlap-check`. It will calculate the overlaps of ALL the program entries 
  in the database. It could take some while.
- Improve recommendation sql search
- Fixing date time offset problem with ARD records.
- Web-Client: Fixing timezone/localization issues in browser client, but there are still browser-specific bugs
- Web-Client: Datepicker added to gui and overlapping program items can be displayed
- HTTP API: Extend LogResponse object of api

## 0.9.11 - 2021/07/31
- Integration of GitHub CI
- Fix of server status request - if there are no program entries
- Multiple small fixes
- Polished `Makefile` and integration of [license-eye check tool](https://github.com/apache/skywalking-eyes)
- More information added to README
- `docker-compose` and `Dockerfile` integration to support the containerized usage of oerc.

## 0.9.10 - 2021/05
- Complete rewrite in Golang, mostly functionally equivalent, but no mass mode fetch.
- Frontend is included

## 0.0.x-0.9.5 – 2019 – Old Java Spring Boot implementation
- initial release, very slow and resource hungry
- Got out-of-date
- GitHub repository was named `emschu/oer-collector`
- If you want to access the old Java source code, please contact the project maintainer.