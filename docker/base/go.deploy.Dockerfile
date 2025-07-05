
ARG BASE_IMAGE=alpine:3.22

FROM ${BASE_IMAGE} AS builder

RUN apk update && apk add curl

ENV PROBE_VERSION=0.4.39
ENV OS=linux
ENV ARCH=amd64

WORKDIR /downloads

RUN curl -XGET https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v${PROBE_VERSION}/grpc_health_probe-${OS}-${ARCH} \
    -o grpc_health_probe

FROM scratch AS final

WORKDIR /
COPY --from=builder /downloads/grpc_health_probe ./

HEALTHCHECK --interval=30s \
    --timeout=30s \
    --start-period=5s \
    --retries=3 \
    CMD [ "grpc_health_probe" "-addr=localhost:50051" ]
