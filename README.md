# Introduction
pgpool-consul is a daemon help to register pgpool instance into consul as service.

# Configuration
Sample configuration. At the moment TLS not supported yet.
```yaml
consul:
  addr: "http://localhost:8500" # Consul client address
  services:
    name: "pgpool" # Service name will be combine with hostname
    addr: "localhost" # Service healthcheck address; Should be using pgpool listen address
    port: 9999 # Service healthcheck port
    interval: "5s" # Service healthcheck interval
    timeout: "2s" # Service healtcheck timeout

pgpool:
  listen: "localhost" # PGPool listen address
  pcppassfile: "~/home/postgres/.pcppass" # This will be use to execute pcp_watchdog_info
  pcpport: 9898 # PCP Port listen on the instance
  pcpuser: "pgpool" # # PCP User to execute pcp_watchdog_info
  id: "0" # Refer to pgpool_node_id file

global:
  interval: 3 # Interval for pgpool-consul service to check pgpool state

```

# Build
Follow steps as below to build pgpool-consul. Build directory will be created and binary will be included in build directory
```bash
make build

```

# Testing
Make sure consul client is up and running in dev environment. To mock `pcp_watchdog_info` please refer to `script` folder for a sample `pcp_watchdog_info` command.
```bash
# Clean go test cache 
make clean-test

# Run go test
make test

```
