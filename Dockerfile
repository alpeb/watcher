FROM golang:1.18-alpine as go-lang
WORKDIR /
COPY * /
RUN CGO_ENABLED=0 go build .
ENTRYPOINT ["/watcher"]

