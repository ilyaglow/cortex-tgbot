package cortexbot

import (
	"log"
	"strconv"
	"strings"

	"github.com/ilyaglow/telegram-bot-api"
)

// Auth handles simple password authentication of a user
func (c *Client) Auth(chatID, userID int, password string) error {
	var message string
	if password == c.Password {
		err := c.registerUser(userID)
		if err != nil {
			return err
		}
		log.Printf("New user %s was registered!\nAllowed users: %s", userID, strings.Join(c.listUsers(), ","))
		message = "Successfully authenticated"
	} else {
		message = "Wrong password"
	}

	if err := c.Bot.Send(&tb.Chat{chatID}, message); err != nil {
		return err
	}
	return nil
}

// CheckAuth checks if user is allowed to interact with a bot
func (c *Client) CheckAuth(u *tgbotapi.User) bool {
	if c.userExists(strconv.Itoa(u.ID)) {
		return true
	}
	return false
}
