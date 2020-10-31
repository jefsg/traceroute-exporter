# Traceroute Exporter

  A Prometheus Exporter for the `traceroute` command line utility. Currently calls the CLI traceroute utility via `exec`.

  Exposes the following metrics on the /trace endpoint

  * route_hop_count - Number of hops taken along route
  * route_success - Indicates whether the trace was successful (1 = success, 0 = failure, -1 = exporter error). 
  * route_hop_latency - Latency to the indicated hop (in seconds)

## Prometheus Configuration

  Similar to the official snmp exporter, the traceroute exporter expects to be passed an address or hostname as a parameter. This can be achieved through relabeling> Example:

  ```
  scrape_configs:
  - job_name: 'traceroute'
    scrape_interval: 60s
    scrape_timeout: 60s
    static_configs:
      - targets:
        - example.com  # host or address to be passed to traceroute.
    metrics_path: /trace
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: target
      - target_label: __address__
        replacement: 127.0.0.1:9805  # The traceroute exporter's real hostname:port.
  ```

## Build Manually

  ```
  git clone https://github.com/jefsg/traceroute-exporter.git
  go get .
  go build .
  ```

## Prebuilt Docker Image

  [DockerHub](https://hub.docker.com/repository/docker/jefsg/traceroute-exporter)

## TODO

  1. Refactor to use a Go based traceroute library.