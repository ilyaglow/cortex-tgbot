[![Build Status](https://travis-ci.org/ilyaglow/cortex-tgbot.svg?branch=master)](https://travis-ci.org/ilyaglow/cortex-tgbot) [![](https://godoc.org/github.com/ilyaglow/cortex-tgbot?status.svg)](http://godoc.org/github.com/ilyaglow/cortex-tgbot) [![Codacy Badge](https://api.codacy.com/project/badge/Grade/a75cbc20a3524962bb182814048cd186)](https://www.codacy.com/app/ilyaglow/cortex-tgbot?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=ilyaglow/cortex-tgbot&amp;utm_campaign=Badge_Grade) [![Coverage Status](https://coveralls.io/repos/github/ilyaglow/cortex-tgbot/badge.svg?branch=master)](https://coveralls.io/github/ilyaglow/cortex-tgbot?branch=master)

Cortex bot
----------

Simple telegram bot to check indicators' reputation based on [Cortex](https://github.com/CERT-BDF/Cortex) [analyzers](https://github.com/CERT-BDF/Cortex-Analyzers) that can be [easily written](https://github.com/CERT-BDF/CortexDocs/blob/master/api/how-to-create-an-analyzer.md) for any third party feeds or your own API service.

It simply uses a password for authentication, which is probably will be changed in the future prior to a role based model.

## Usage

Start bot from the source code (you can use [compiled version](https://github.com/ilyaglow/cortex-tgbot/releases) too):

```
CORTEX_BOT_PASSWORD=PassphraseForAuth CORTEX_LOCATION=http://127.0.0.1:9000 TGBOT_API_TOKEN=TOKEN go run cmd/cortexbot/main.go
```

Add bot to your contacts, enter the passphrase and here you go - submit data and wait for results.

## Supported data types

By now the following data types are supported for lookup:
* Domain
* Hash
* IP
* URL
* File


PRs are welcome!
