
ARG BASE_IMAGE

FROM ${BASE_IMAGE} AS builder

WORKDIR /usr/src/lib

COPY lib ./

WORKDIR /usr/src/app

COPY dns/go.mod dns/go.sum ./
RUN go mod download && go mod verify

COPY dns ./
RUN CGO_ENABLED=0 go build \
    -ldflags="-s -w" -trimpath \
    -o /usr/bin/server \
    ./cmd/resolver

FROM scratch 

COPY --from=builder /usr/bin/server ./

CMD [ "./server" ]
