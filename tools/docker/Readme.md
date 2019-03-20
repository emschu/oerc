# Docker usage with oer-server


## Squid proxy (recommended in mass_mode)
```bash
$ cd squid
$ sudo docker build -t squid .
$ sudo docker run -d -p 3128:3128 --restart=always --name cache-proxy squid
```

### Monitoring squid proxy (in container)
```bash
# squidclient -h localhost cache_object://localhost/ mgr:info
# squidclient -h localhost cache_object://localhost/counters
```