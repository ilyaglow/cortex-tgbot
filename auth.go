package cortexbot

import (
	"log"
	"strings"

	"github.com/ilyaglow/telegram-bot-api"
)

// Auth handles simple password authentication of a user
func (c *Client) Auth(input *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(input.Chat.ID, "")
	msg.ReplyToMessageID = input.MessageID
	if input.Text == c.Password {
		c.registerUser(input.From.UserName, "password")
		log.Printf("Allowed users: %s", strings.Join(c.listUsers(), ","))
		msg.Text = "Successfully authenticated"
	} else {
		msg.Text = "Wrong password"
	}
	c.Bot.Send(msg)
}

// CheckAuth checks if user is allowed to interact with a bot
func (c *Client) CheckAuth(u string) bool {
	if c.userExists(u) {
		return true
	}
	return false
}
