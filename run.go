package cortexbot

import (
	"log"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

func (c *Client) handleMessage(m *tb.Message) {
	if !c.CheckAuth(m.Sender.ID) {
		if err := c.Auth(m); err != nil {
			log.Println(err)
		}
		return
	}
	go func() {
		if err := c.processMessage(m); err != nil {
			log.Println(err)
		}
	}()
}

func (c *Client) handleCallback(m *tb.Callback) {
	if !c.CheckAuth(m.Sender.ID) {
		if err := c.Bot.Respond(m, &tb.CallbackResponse{
			CallbackID: m.ID,
			Text:       "Something fishy going on here :)",
		}); err != nil {
			log.Println(err)
		}
		return
	}

	go func() {
		msg := ""
		if err := c.processCallback(m); err != nil {
			log.Println(err)
			msg = err.Error()
		}
		if err := c.Bot.Respond(m, &tb.CallbackResponse{
			Text: msg,
		}); err != nil {
			log.Println(err)
		}
	}()
}

// Run represents infinite function that waits for a message,
// authenticate user and process task
func (c *Client) Run() {
	defer c.DB.Close()

	c.Bot.Handle(tb.OnCallback, c.handleCallback)
	c.Bot.Handle(tb.OnText, c.handleMessage)
	c.Bot.Handle(tb.OnDocument, c.handleMessage)

	log.Printf("Users in database: %s", strings.Join(c.listUsers(), ","))
	c.Bot.Start()
}
