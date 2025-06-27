
ARG BASE_IMAGE=golang:1.24.4-alpine3.22

FROM ${BASE_IMAGE} AS busybox

RUN apk update

ENV HEALTH_PROBE_VERSION=v0.4.39

RUN go install github.com/grpc-ecosystem/grpc-health-probe@${HEALTH_PROBE_VERSION}
