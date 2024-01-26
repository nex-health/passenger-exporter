ARG GO_VERSION=1.21
ARG XX_VERSION=1.3.0

FROM --platform=$BUILDPLATFORM tonistiigi/xx:${XX_VERSION} AS xx
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} as builder

COPY --from=xx / /

ARG TARGETPLATFORM

WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY cmd/ cmd/
COPY collector/ collector/

ENV CGO_ENABLED=0
RUN xx-go build -trimpath -a -o passenger_exporter cmd/passenger_exporter/main.go

FROM scratch
LABEL maintainer="Aurel Canciu <aurel.canciu@nexhealth.com>"

COPY --from=builder /workspace/passenger_exporter /bin/passenger_exporter

USER nobody
EXPOSE 9149
ENTRYPOINT [ "/bin/passenger_exporter" ]