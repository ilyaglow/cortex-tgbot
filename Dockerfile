FROM golang:alpine as build
LABEL maintainer="contact@ilyaglotov.com" \
      repository="https://github.com/ilyaglow/cortex-tgbot"

ENV GO111MODULE=on

COPY . /go/src/github.com/ilyaglow/cortex-tgbot/

RUN apk --update --no-cache add ca-certificates \
                                git \
                                sqlite \
                                build-base \
  && cd /go/src/github.com/ilyaglow/cortex-tgbot \
  && go mod download \
  && cd ./cmd/cortexbot \
  && go build -ldflags="-s -w" \
              -a \
              -installsuffix static \
              -o /cortexbot

FROM alpine:latest
COPY --from=build /cortexbot /app/cortexbot

RUN apk --update --no-cache add ca-certificates \
  && adduser -h /app -D app \
  && mkdir -p /app/data \
  && chown -R app /app

USER app

VOLUME /app/data

WORKDIR /app/data

ENTRYPOINT ["../cortexbot"]
