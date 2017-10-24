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
