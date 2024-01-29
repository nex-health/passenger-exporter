# Passenger Exporter

![GitHub Actions](https://github.com/nex-health/passenger-exporter/actions/workflows/ci.yml/badge.svg)

Export Passenger metrics to Prometheus.

To run it:

```bash
make
./passenger_exporter [flags]
```

## Exported Metrics

| Metric                             | Meaning                                                  | Type    |
| ---------------------------------- | -------------------------------------------------------- | ------- |
| passenger_up                       | Passenger state.                                         | Gauge   |
| passenger_version                  | Phusion Passenger version.                               | Gauge   |
| passenger_top_level_queue          | Number of requests in the top-level queue.               | Gauge   |
| passenger_max_processes            | Configured maximum number of processes.                  | Gauge   |
| passenger_current_processes        | Current number of processes.                             | Gauge   |
| passenger_app_count                | Number of apps.                                          | Gauge   |
| passenger_app_queue                | Number of requests in app process queues.                | Gauge   |
| passenger_app_group_queue          | Number of requests in app group process queues.          | Gauge   |
| passenger_app_procs_spawning       | Number of processes spawning.                            | Gauge   |
| passenger_requests_processed_total | Number of processes served by a process.                 | Counter |
| passenger_current_sessions         | Number of sessions currently being handled by a process. | Gauge   |
| passenger_proc_start_time_seconds  | Number of seconds since processor started.               | Gauge   |
| passenger_proc_memory              | Memory consumed by a process.                            | Gauge   |

### Flags

```bash
./passenger_exporter --help
```

* __`passenger.instance-registry`:__ Path to the instance registry directory.
  (default: /tmp)
* __`passenger.pid-file`:__ Optional path to a file containing the
  passenger/nginx PID for additional metrics.
* __`log.format`:__ Output format of log messages. One of: [logfmt, json]
  (default: `logfmt`).
* __`log.level`:__ Only log messages with the given severity or above. One of:
  [debug, info, warn, error] (default: `info`).
* __`web.listen-address`:__ Addresses on which to expose metrics and web
  interface. Repeatable for multiple addresses (default: `:9144`).
* __`web.telemetry-path`:__ Path under which to expose metrics (default: `/metrics`).
* __`version`:__ Show application version.

## Using Containers

You can run this exporter using the [ghcr.io/nex-health/passenger-exporter](https://github.com/nex-health/passenger-exporter/pkgs/container/passenger-exporter) container image.

```bash
docker run -d \
           -p 9149:9149 \
           -v /path/to/passenger-inst-reg-dir:/var/run/passenger-instreg:ro \
           -v /path/to/passenger.pid:/tmp/passenger.pid:ro \
           ghcr.io/nex-health/passenger-exporter \
             --passenger.pid-file /tmp/passenger.pid \
             --passenger.instance-registry /var/run/passenger-instreg
```
