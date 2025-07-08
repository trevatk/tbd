
FROM golang:1.24.4-alpine AS builder

WORKDIR /usr/src/lib

COPY lib/ ./

WORKDIR /usr/src/app

COPY wellknown/go.mod wellknown/go.sum ./
RUN go mod download && go mod verify

COPY wellknown/ ./

ENV CGO_ENABLED=0 

RUN go build -ldflags="-s -w" -trimpath -v -o /usr/local/bin/wellknown ./cmd/server

FROM tbd/go-deploy:latest AS final

COPY --from=builder /usr/local/bin/wellknown /

HEALTHCHECK --interval=30s \
    --timeout=30s \
    --start-period=5s \
    --retries=3 \
    CMD [ "grpc_health_probe", "-addr=localhost:50051" ]

VOLUME [ "/config" ]

CMD [ "./wellknown" ]
