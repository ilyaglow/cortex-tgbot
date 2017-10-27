[![Build Status](https://travis-ci.org/ilyaglow/cortex-tgbot.svg?branch=master)](https://travis-ci.org/ilyaglow/cortex-tgbot) [![](https://godoc.org/github.com/ilyaglow/cortex-tgbot?status.svg)](http://godoc.org/github.com/ilyaglow/cortex-tgbot) [![Codacy Badge](https://api.codacy.com/project/badge/Grade/a75cbc20a3524962bb182814048cd186)](https://www.codacy.com/app/ilyaglow/cortex-tgbot?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=ilyaglow/cortex-tgbot&amp;utm_campaign=Badge_Grade) [![Coverage Status](https://coveralls.io/repos/github/ilyaglow/cortex-tgbot/badge.svg?branch=master)](https://coveralls.io/github/ilyaglow/cortex-tgbot?branch=master)

Cortex bot
----------

Simple telegram bot that provides capabilities to work with [Cortex](https://github.com/CERT-BDF/Cortex).

It uses a password to validate it's user, that probably will be changed in the future.

## Usage

Start bot from source code (you can use compiled version (linux/amd64) too):

```
CORTEX_BOT_PASSWORD=PassphraseForAuth CORTEX_LOCATION=http://127.0.0.1:9000 TGBOT_API_TOKEN=TOKEN go run cmd/cortexbot/main.go
```

Add bot to your contacts, enter the passphrase and here you go - submit data and wait for results.

## Supported data types

By now the following data types are supported:
* Domain
* Hash
* IP
* URL


PRs are welcome!
