FROM golang:alpine as build
LABEL maintainer="contact@ilyaglotov.com" \
      repository="https://github.com/ilyaglow/cortex-tgbot"

COPY cmd/cortexbot/main.go /go/src/cortexbot/main.go

RUN apk --update --no-cache add ca-certificates \
                                git \
  && cd /go/src/cortexbot/ \
  && go get -t . \
  && CGO_ENABLED=0 go build -ldflags="-s -w" \
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
