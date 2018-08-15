FROM alpine:latest
LABEL maintainer "contact@ilyaglotov.com"
LABEL repository "https://github.com/ilyaglow/cortex-tgbot"

ENV CORTEXBOT_VERSION "0.9.2"

ADD https://github.com/ilyaglow/cortex-tgbot/releases/download/v${CORTEXBOT_VERSION}/cortexbot-amd64 /app/cortexbot

RUN chmod +x /app/cortexbot

RUN adduser -D app \
  && chown -R app /app

USER app

VOLUME /app

WORKDIR /app

ENTRYPOINT ["./cortexbot"]
