FROM golang:1.19 AS builder
ENV CGO_ENABLED 0
ARG VERSION
WORKDIR /go/src/app
ADD . .
RUN go build -ldflags="-X main.GitHash=$(git rev-parse --short HEAD)" -o /minit

FROM busybox
COPY --from=builder /minit /minit
ENTRYPOINT ["/minit"]
