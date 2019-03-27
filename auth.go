package cortexbot

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// Auth handles simple password authentication of a user
func (c *Cortexbot) Auth(input *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(input.Chat.ID, "")
	msg.ReplyToMessageID = input.MessageID
	if input.Text == c.Password {
		u := &User{
			ID:    input.From.ID,
			Admin: 1,
			About: input.From.String(),
		}
		err := c.addUser(u)
		if err != nil {
			return err
		}

		msg.Text = "Successfully authenticated"
	} else {
		msg.Text = "Wrong password"
	}
	if _, err := c.Bot.Send(msg); err != nil {
		return err
	}

	return nil
}

// CheckAuth checks if user is allowed to interact with a bot.
func (c *Cortexbot) CheckAuth(u *tgbotapi.User) bool {
	return c.userExists(u.ID)
}

// CheckAdmin checks if user is an admin.
func (c *Cortexbot) CheckAdmin(u *tgbotapi.User) bool {
	user, err := c.getUser(u.ID)
	if err != nil {
		return false
	}
	if user.Admin == 1 {
		return true
	}
	return false
}
