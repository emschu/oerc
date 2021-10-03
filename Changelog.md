# Changelog oerc

## 0.9.12 - 2021/10/03
- Fixing date time offset problem with ARD records. Update existing database entries *before* 
  the first run of `oerc fetch` by executing the following query ONCE:
  - `UPDATE program_entries SET start_date_time = start_date_time - interval '2 hours' where channel_family_id = 1;`
- Package updates
- Improve recommendation quality

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