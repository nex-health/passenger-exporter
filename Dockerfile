ARG GO_VERSION=1.21

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} as builder

WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY cmd/ cmd/
COPY collector/ collector/

ENV CGO_ENABLED=0
RUN go build -trimpath -a -o passenger_exporter cmd/passenger_exporter/main.go

FROM scratch
LABEL maintainer="Aurel Canciu <aurel.canciu@nexhealth.com>"

COPY --from=builder /workspace/passenger_exporter /bin/passenger_exporter

USER nobody
EXPOSE 9149
ENTRYPOINT [ "/bin/passenger_exporter" ]