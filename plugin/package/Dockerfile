FROM kong/go-plugin-tool:latest-alpine-latest as builder
ENV GO111MODULE=on

RUN mkdir /go-plugins
COPY ./cmd/ /go-plugins/cmd/
COPY ./pkg/ /go-plugins/pkg/
COPY ./go.mod /go-plugins/
COPY ./go.sum /go-plugins/
RUN cd /go-plugins && \
    go build -o /go-plugins/bin/serverless cmd/main.go

FROM kong:2.8.1

COPY --from=builder /go-plugins/bin/serverless /usr/local/bin/serverless
