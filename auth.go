package cortexbot

import (
	"log"
	"strconv"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

// Auth handles simple password authentication of a user
func (c *Client) Auth(m *tb.Message) error {
	var message string
	if m.Text == c.Password {
		err := c.registerUser(m.Sender)
		if err != nil {
			return err
		}
		log.Printf("New user %s was registered!\nAllowed users: %s", m.Sender.Username, strings.Join(c.listUsers(), ","))
		message = "Successfully authenticated"
	} else {
		message = "Wrong password"
	}

	if _, err := c.Bot.Reply(m, message); err != nil {
		return err
	}
	return nil
}

// CheckAuth checks if user is allowed to interact with a bot
func (c *Client) CheckAuth(uid int) bool {
	if c.userExists(strconv.Itoa(uid)) {
		return true
	}
	return false
}
