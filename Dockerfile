ARG GO_VERSION=1.25

FROM golang:${GO_VERSION} AS builder

ARG TARGETOS
ARG TARGETARCH

RUN apt-get update -yqq && apt-get install -yqq git

ARG PROMU_VERSION=0.17.0
ADD https://github.com/prometheus/promu/releases/download/v${PROMU_VERSION}/promu-${PROMU_VERSION}.${TARGETOS}-${TARGETARCH}.tar.gz /tmp/
RUN tar xf /tmp/promu-${PROMU_VERSION}.${TARGETOS}-${TARGETARCH}.tar.gz -C /tmp && \
    install -m 755 /tmp/promu-${PROMU_VERSION}.${TARGETOS}-${TARGETARCH}/promu /usr/bin

WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY cmd/ cmd/
COPY collector/ collector/

COPY .promu.yml .promu.yml
COPY VERSION VERSION

ENV CGO_ENABLED=0
RUN --mount=target=.git,type=bind,source=.git \
    promu build

FROM quay.io/prometheus/busybox:latest
LABEL maintainer="Aurel Canciu <aurel.canciu@nexhealth.com>"

COPY --from=builder /workspace/passenger_exporter /bin/passenger_exporter

USER nobody
EXPOSE 9149
ENTRYPOINT [ "/bin/passenger_exporter" ]