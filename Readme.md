# Traceroute Exporter

  A Prometheus Exporter for *nix `traceroute` command. Currently calls the CLI traceroute utility via `exec`.

  Exposes the following metrics

  * route_hop_count - Number of hops taken along route
  * route_success - Indicates whether the trace was successful (1 = success, 0 = failure, -1 = exporter error). 
  * route_hop_latency - Latency to the indicated hop (in seconds)

## Build Manually

  ```Go
  git clone https://github.com/jefsg/traceroute-exporter.git
  go get .
  go build .
  ```

## Prebuilt Docker Image

  [DockerHub](https://hub.docker.com/repository/docker/jefsg/traceroute-exporter)

## TODO

  1. Refactor to use a Go based traceroute library.