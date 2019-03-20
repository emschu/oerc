# Oer-Collector Development guide
 
3 spring environment profiles: 
- test
- dev
- prod

use EnvService to access.

# Maven modules

- **oer-collector:** fills the db
- **oer-collector-server:** exposes data
- **oer-core:** is shared by collector and web server

**JDBC/data access relevant classes:** *org.emschu.oer.core.model*
