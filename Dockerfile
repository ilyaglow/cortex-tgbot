FROM alpine:latest
LABEL maintainer="contact@ilyaglotov.com" \
      repository="https://github.com/ilyaglow/cortex-tgbot"

ENV CORTEXBOT_VERSION "0.9.5"
RUN apk --update --no-cache add ca-certificates \
  && mkdir app \
  && wget -O /app/cortexbot.tar.gz https://github.com/ilyaglow/cortex-tgbot/releases/download/v${CORTEXBOT_VERSION}/cortexbot_${CORTEXBOT_VERSION}_linux_amd64.tar.gz \
  && cd /app \
  && tar xzf cortexbot.tar.gz \
  && chmod +x /app/cortexbot \
  && adduser -D app \
  && chown -R app /app

USER app

VOLUME /app

WORKDIR /app/data

ENTRYPOINT ["../cortexbot"]
