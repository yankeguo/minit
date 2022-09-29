FROM golang:1.19 AS builder
ENV CGO_ENABLED 0
ARG VERSION
WORKDIR /go/src/app
ADD . .
RUN go build -mod vendor -ldflags="-X main.GitHash=$(git rev-parse --short HEAD) -X main.Version=${VERSION}" -o /minit

FROM busybox
COPY --from=builder /minit /minit
ENTRYPOINT ["/minit"]
