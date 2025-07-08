
FROM golang:1.24.4-alpine AS builder

WORKDIR /usr/src/lib

COPY lib/ ./

WORKDIR /usr/src/app

COPY dns/go.mod dns/go.sum ./
RUN go mod download && go mod verify

COPY dns/ ./

ENV CGO_ENABLED=0 

RUN go build -ldflags="-s -w" -trimpath -v -o /usr/local/bin/nameserver ./cmd/nameserver

FROM tbd/go-deploy:latest AS final

COPY --from=builder /usr/local/bin/nameserver /

CMD [ "./nameserver" ]
