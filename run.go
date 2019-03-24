package cortexbot

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

const defaultPollTimeout = 20

// Run represents infinite function that waits for a message,
// authenticate user and process task
func (c *Cortexbot) Run(ctx context.Context) error {
	defer c.DB.Close()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = defaultPollTimeout

	updates, err := c.Bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}
	log.Printf("Authorized on account %s", c.Bot.Self.UserName)

	admins, err := c.getAdmins()
	if err != nil {
		return err
	}

	var adminsInfo []string
	for i := range admins {
		adminsInfo = append(adminsInfo, fmt.Sprintf(
			"%d (%s)",
			admins[i].ID,
			admins[i].About,
		))
	}
	if len(admins) == 0 {
		log.Printf("No administrators registered")
	} else {
		log.Printf("Administrators registered: %s", strings.Join(adminsInfo, ", "))
	}

	for {
		select {
		case <-ctx.Done():
			return errors.New("got ctrl-c")
		case update := <-updates:
			go func(upd *tgbotapi.Update) {
				if err := c.processUpdate(upd); err != nil {
					log.Println(err)
				}
			}(&update)
		}
	}

	return nil
}
