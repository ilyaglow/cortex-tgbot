package cortexbot

import (
	"log"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

// Run represents infinite function that waits for a message,
// authenticate user and process task
func (c *Client) Run() {
	defer c.DB.Close()

	c.Bot.Handle(tb.OnCallback, func(m *tb.Callback) {
		if !c.CheckAuth(m.Sender.ID) {
			if err := c.Bot.Respond(m, &CallbackResponse{
				CallbackID: m.ID,
				Text:       "Something fishy going on here :)",
			}); err != nil {
				log.Println(err)
			}
			return
		}

		go func() {
			if err := c.processCallback(m); err != nil {
				log.Println(err)
			}
		}()
	})

	c.Bot.Handle(tb.OnMessage, func(m *tb.Message) {
		if !c.CheckAuth(m.Sender.ID) {
			if err := c.Auth(m.Chat.ID, m.Sender.ID, m.Text); err != nil {
				log.Println(err)
			}
			return
		}
		go func() {
			if err := c.processMessage(m); err != nil {
				log.Println(err)
			}
		}()
	})

	c.Bot.Handle(tb.OnDocument, func(m *tb.Message) {

	})

	log.Printf("Users in database: %s", strings.Join(c.listUsers(), ","))

	for update := range updates {

		if update.CallbackQuery != nil {
			log.Printf("username: %s, id: %d, text: %s", update.CallbackQuery.Message.From.UserName, update.CallbackQuery.Message.From.ID, update.CallbackQuery.Message.Text)
			go func() {
				if err := c.processCallback(update.CallbackQuery); err != nil {
					log.Println(err)
				}
			}()
		} else {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			msg.ReplyToMessageID = update.Message.MessageID

			if update.Message.IsCommand() &&
				update.Message.Command() == "start" &&
				!c.CheckAuth(update.Message.From) {
				msg.Text = "Enter your password"
				go c.Bot.Send(msg)
				continue
			}

			if c.CheckAuth(update.Message.From) {
				go func() {
					if err := c.processMessage(update.Message); err != nil {
						log.Println(err)
					}
				}()
			} else {
				c.Auth(update.Message)
			}
		}
	}
}
