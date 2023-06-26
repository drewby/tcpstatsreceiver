FROM golang:1-bullseye AS build-env

RUN go install go.opentelemetry.io/collector/cmd/builder@latest
RUN go install github.com/open-telemetry/opentelemetry-collector-contrib/cmd/mdatagen@latest

WORKDIR /otelcol
COPY ./tcpstatsreceiver ./tcpstatsreceiver
COPY builder-config.yaml builder-config.yaml
RUN builder --config builder-config.yaml

FROM golang:1-bullseye

ARG USER_UID=10001
USER ${USER_UID}

COPY --chmod=755 --from=build-env otelcol/otelcol /otelcol
COPY ./otelcol.yaml /etc/otelcol/config.yaml
ENTRYPOINT ["/otelcol"]
CMD ["--config", "/etc/otelcol/config.yaml"]
EXPOSE 4317 55678 55679