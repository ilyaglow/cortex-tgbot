package cortexbot

import (
	"log"

	"gopkg.in/telegram-bot-api.v4"
)

func (c *Client) Run() {
	defer c.DB.Close()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := c.Bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Authorized on account %s", c.Bot.Self.UserName)

	for update := range updates {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.ReplyToMessageID = update.Message.MessageID

		if update.Message.IsCommand() &&
			update.Message.Command() == "start" &&
			!c.CheckAuth(update.Message.From.UserName) {

			msg.Text = "Enter your password"
			c.Bot.Send(msg)
			continue
		}

		if c.CheckAuth(update.Message.From.UserName) {
			if err := c.ProcessCortex(update.Message); err != nil {
				msg.Text = err.Error()
				c.Bot.Send(msg)
			}
		} else {
			c.Auth(update.Message)
		}
	}
}
