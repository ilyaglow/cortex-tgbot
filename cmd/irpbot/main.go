package main

import (
	"irpbot"
	"log"

	"gopkg.in/telegram-bot-api.v4"
)

func main() {
	c := irpbot.NewClient()
	c.Bot.Debug = true
	log.Printf("Authorized on account %s", c.Bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := c.Bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {

		if update.CallbackQuery != nil {
			go c.ProcessCallback(update.CallbackQuery)
			continue
		}

		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			go c.ProcessCommand(update.Message)
			continue
		}
	}
}
