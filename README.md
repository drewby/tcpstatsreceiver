# Building an OpenTelemetry Collector

This guide provides a hands-on example of creating a custom [OpenTelemetry Collector](https://github.com/open-telemetry/opentelemetry-collector) (OTel Collector) distribution, with a unique [receiver](tcpstatsreceiver/README.md) to scrape TCP stats from Linux-based systems.

Utilizing the Go-based [OpenTelemetry Collector Builder](https://github.com/open-telemetry/opentelemetry-collector/tree/main/cmd/builder), this custom distribution can be smoothly packaged into a Docker image.

Components in the build are:

* `prometheus` exporter: Provides a /metrics endpoint for Prometheus scraping.
* `attributes` processor: Enhances telemetry data with supplementary attributes like host name.
* `tcpstats` receiver: Reads /proc/net/tcp to generate metrics on queue size and length.

## Repository Breakdown

1. [`builder-config.yaml`](./builder-config.yaml): Configures the OTel Collector Builder with details such as the OTel Collector version, the build's output path, and the exporters, processors, and receivers to include in the distribution.

1. [`otelcol.yaml`](./otelcol.yaml): Specifies the configuration for the OTel Collector, denoting the receivers, processors, and exporters the collector should utilize. This sets up a metrics pipeline including:
    * `tcpstats` receiver: Generates metrics from /proc/net/tcp,
    * `attributes` processor: Adds a host.name attribute to the metrics,
    * `prometheus` exporter: Creates a /metrics endpoint on port 8889.

1. [`Dockerfile`](./Dockerfile): Directs the building of a Docker image for the custom OTel Collector distribution. It first constructs the OTel Collector using the builder, and then copies the built collector and its configuration into a minimal Docker image.

1. [`Makefile`](./Makefile): Contains targets for installing dependencies, compiling the binary, and crafting a Docker image.

1. [`tcpstatsreceiver`](./tcpstatsreceiver/README.md): This is a custom receiver implemented to generate metrics for TCP queue size and length.

## Steps to Follow

### Install Dependencies

Ensure you have `go` version 1.19 or later installed and available in your path. 

Run:

```shell
make setup
```

This installs `builder@latest` and `mdatagen@latest`. The `builder` is used to process `builder-config.yaml`, fetch component sources, and build the binary. `mdatagen` generates metric functions from the `metadata.yaml` data and stores them in `tcpstatsreceiver/internal`.

### Build and Run the Binary

Run `make all` or `make` to build the binary.

```shell
make
```

If `metadata.yaml` has been updated, `mdatagen` will regenerate the metric functions.

The binaries are by default output to `otelcol-dev`, but this can be adjusted in the `builder-config.yaml` file.

To run the newly built binary, use:

```shell
./otelcol-dev/otelcol --config otelcol.yaml
```

### Build and Launch the Docker Image

Build the Docker image with:

```shell
make docker
```

The Docker image, defined in `Dockerfile`, will be tagged as `otelcol:latest`.

Run the Docker image with:

```shell
docker run -p 8889:8889 -v $(pwd)/otelcol.yaml:/etc/otelcol/config.yaml -v $(pwd)/tcpstatsreceiver/testdata:/testdata otelcol:latest
```

This command starts the OTel Collector, collecting files matching the pattern (`/testdata/*.log_bucket`). The `/testdata` directory is mapped to a local folder, while otelcol.yaml is mapped to /etc/otelcol/config.yaml.

# Test the TCP Stats Receiver

A test TCP server is implemented in the [testapp/main.go](testapp/main.go). It's designed to delay responses, causing the TCP queue and length values to increase. 

Start the server with:

```shell
go run testapp/main.go --listenAddr 127.0.0.1:8005 --sleepSeconds 5
```

Both arguments are optional with default values provided.

With the server running, send TCP requests to the listening address from a separate shell:

```shell
for count in $(seq 1 10); do
  echo -n "Some data to cause the queue size to grow" | nc localhost 8005 &
done
```

This sends 10 requests to the server, answered sequentially every five seconds.

Check the metric values at http://localhost:8889.

# Using the Dev Container

The Dev Container is configured with several services including all the dependencies required to build and run the custom Open Telemetry collector. 

It also contains Prometheus and Grafana, allowing you to experiment with querying and visualizing the new TCP metrics.

Start the Dev Container in Visual Studio Code or as a GitHub Codespace. After starting the dev container you can launch either tool:

Prometheus: http://localhost:9090

Grafana: http://localhost:3000.

## Disclaimer

This project is an example and may not be fully optimized for a production environment. We strongly recommend thorough testing of any configuration before deploying in a live production setting.
