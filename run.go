package cortexbot

import (
	"log"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

const defaultPollTimeout = 20

// Run represents infinite function that waits for a message,
// authenticate user and process task
func (c *Client) Run() {
	defer c.DB.Close()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = defaultPollTimeout

	updates, err := c.Bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Authorized on account %s", c.Bot.Self.UserName)
	log.Printf("Users in database: %s", strings.Join(c.listUsers(), ","))

	for update := range updates {

		if update.CallbackQuery != nil {
			log.Printf("username: %s, id: %d, text: %s", update.CallbackQuery.Message.From.UserName, update.CallbackQuery.Message.From.ID, update.CallbackQuery.Message.Text)
			go func() {
				if err := c.processCallback(update.CallbackQuery); err != nil {
					log.Println(err)
				}
				cbcfg := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
				if _, err := c.Bot.AnswerCallbackQuery(cbcfg); err != nil {
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
				c.log("run: new unauthorized user")
				msg.Text = "Enter your password"
				if _, err := c.Bot.Send(msg); err != nil {
					log.Println(err)
				}
				continue
			}

			if c.CheckAuth(update.Message.From) {
				c.log("run: auth succeeded")
				go func() {
					if err := c.processMessage(update.Message); err != nil {
						log.Println(err)
					}
				}()
			} else {
				c.log("run: wrong password from the user")
				if err := c.Auth(update.Message); err != nil {
					log.Println(err)
				}
			}
		}
	}
}
