FROM alpine:latest
LABEL maintainer "contact@ilyaglotov.com"
LABEL repository "https://github.com/ilyaglow/cortex-tgbot"

ENV CORTEXBOT_VERSION "0.9.3"
ADD https://github.com/ilyaglow/cortex-tgbot/releases/download/v${CORTEXTBOT_VERSION}/cortexbot_linux_amd64.tar.gz /app/

RUN chmod +x /app/cortexbot

RUN adduser -D app \
  && chown -R app /app

USER app

VOLUME /app

WORKDIR /app

ENTRYPOINT ["./cortexbot"]
