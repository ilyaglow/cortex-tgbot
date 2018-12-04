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
		go func(upd *tgbotapi.Update) {
			if err := c.processUpdate(upd); err != nil {
				log.Println(err)
			}
		}(&update)
	}
}
