ARG VERSION="unknown"

FROM golang:1.23 AS builder

ARG VERSION

ENV CGO_ENABLED="0"

WORKDIR /go/src/app

ADD . .

RUN go build -ldflags="-X main.AppVersion=${VERSION}" -o /minit

FROM busybox

COPY --from=builder /minit /minit

ENTRYPOINT ["/minit"]
