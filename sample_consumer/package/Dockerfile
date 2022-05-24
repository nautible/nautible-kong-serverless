FROM golang:1.18.0 as builder

WORKDIR /go/src

COPY ./cmd/  ./cmd/
COPY ./go.mod  ./go.mod

ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64
RUN go build \
    -ldflags '-s -w' \
    -o /go/bin/main cmd/main.go

FROM scratch as runner

COPY --from=builder /go/bin/main /app/main

ENTRYPOINT ["/app/main"]