package cortexbot

import (
	"log"

	"gopkg.in/telegram-bot-api.v4"
)

// Auth handles simple password authentication of a user
func (c *Client) Auth(input *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(input.Chat.ID, "")
	msg.ReplyToMessageID = input.MessageID
	if input.Text == c.Password {
		c.AllowedUsernames[input.From.UserName] = true
		log.Printf("Allowed users: %v", c.AllowedUsernames)
		msg.Text = "Successfully authenticated"
	} else {
		msg.Text = "Wrong password"
	}
	c.Bot.Send(msg)
}

// CheckAuth checks if user is allowed to interact with a bot
func (c *Client) CheckAuth(u string) bool {
	if _, ok := c.AllowedUsernames[u]; ok {
		return true
	}
	return false
}
