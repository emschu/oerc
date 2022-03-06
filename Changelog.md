# Changelog oerc

## 0.9.16 - 2022/03/XX
- Improving code structure in general which brings more performance
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
- New Command: `oerc full-overlap-check`. It will calculate the overlaps on ALL of the program entries 
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