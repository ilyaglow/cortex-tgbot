FROM scratch
LABEL maintainer "contact@ilyaglotov.com"
LABEL repository "https://github.com/ilyaglow/cortex-tgbot"

ENV CORTEXBOT_VERSION "0.4"

ADD https://github.com/ilyaglow/cortex-tgbot/releases/download/v${CORTEXBOT_VERSION}/cortexbot-amd64 /app/cortexbot

VOLUME /app

WORKDIR /app

ENTRYPOINT ["./cortexbot"]
