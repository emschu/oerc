# OER Collector Docker Guide

## Prepare (one time)

Build project for one time (needs Apache Maven):

```bash
# in project directory root:
$ sh package.sh
```
OR: If Maven is not available, copy released `oer-collector.jar` + `start_collector.sh` + `stop_collector.sh` from [here](https://github.com/emschu/oer-collector/releases) to `tools/docker/oer_collector` and the same for server component docker directory by hand.

Then:

## Setup and Run
```bash
$ cd tools/docker
$ docker-compose build
$ docker-compose up
```


### Squid proxy (recommended in mass_mode)
```bash
$ cd squid
$ sudo docker build -t squid .
# execute as single container
$ sudo docker run -d -p 3128:3128 --restart=always --name cache-proxy squid
```

#### Monitoring squid proxy (in container)
```bash
# squidclient -h localhost cache_object://localhost/ mgr:info
# squidclient -h localhost cache_object://localhost/counters
```