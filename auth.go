package irbot

import (
	"log"

	"gopkg.in/telegram-bot-api.v4"
)

// auth handle simple password authentication of a user
func (c *Client) auth(input *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(input.Chat.ID, "")
	msg.ReplyToMessageID = input.MessageID
	if input.CommandArguments() == c.Password {
		c.AllowedUsernames[input.From.UserName] = true
		log.Printf("Allowed users: %v", c.AllowedUsernames)
		msg.Text = "Successfully authenticated"
	} else {
		msg.Text = "Wrong password"
	}
	c.Bot.Send(msg)
}

// checkAuth checks if user is allowed to interact with a bot
func (c *Client) checkAuth(u string) bool {
	if _, ok := c.AllowedUsernames[u]; ok {
		return true
	}
	return false
}
