package irpbot

import (
	"encoding/json"
	"fmt"

	"github.com/ilyaglow/go-cortex"
	"gopkg.in/telegram-bot-api.v4"
)

// ProcessCommand works with the command user submitted
func (c *Client) ProcessCommand(input *tgbotapi.Message) error {
	if input.Command() == "password" {
		c.auth(input)
		return nil
	}

	msg := tgbotapi.NewMessage(input.Chat.ID, "")
	msg.ReplyToMessageID = input.MessageID

	if c.checkAuth(input.From.UserName) {
		if input.Command() == "interact" {
			c.processInteraction(input)
		} else if input.Command() == "hive" {
			msg.Text = "Hive is not available by now"
		} else if input.Command() == "cortex" {
			c.processCortex(input)
		} else if input.Command() == "help" {
			msg.Text = "Help is not available by now"
		} else {
			msg.Text = "Unknown command, check your spelling"
		}
		c.Bot.Send(msg)
	} else {
		msg.Text = "Not authorized"
		c.Bot.Send(msg)
	}
	return nil
}

// processInteraction aimed to deal with Inline things
// Work in progress
func (c *Client) processInteraction(input *tgbotapi.Message) {
	msg := interactionMain(&input.Chat.ID)
	c.Bot.Send(msg)
}

// ProcessCallback works with interaction callbacks
// Work in progress
func (c *Client) ProcessCallback(input *tgbotapi.CallbackQuery) {
	if c.checkAuth(input.From.UserName) {
		e := editMessage{
			ChatID:    &input.Message.Chat.ID,
			MessageID: &input.Message.MessageID,
		}

		var mtc tgbotapi.EditMessageTextConfig

		switch option := input.Data; option {
		case "HiveMenu":
			mtc = e.interactionHive()
		case "HiveCases":
			mtc = e.interactionHiveCases()
		case "HiveTasks":
			mtc = e.interactionHiveTasks()
		case "HiveObservables":
			mtc = e.interactionHiveObservables()
		case "CortexMenu":
			mtc = e.interactionCortex()
		case "CortexTasksAdd":
			mtc = e.interactionCortexTasksAdd()
		case "CortexTasksList":
			mtc = e.interactionCortexTasksList()
		}
		c.Bot.Send(mtc)
	}
}

func (c *Client) processCortex(input *tgbotapi.Message) {
	// Fill the JobBody struct
	// TLP is hardcoded by now
	j := &gocortex.JobBody{
		Data: input.CommandArguments(),
		Attributes: gocortex.ArtifactAttributes{
			DataType: "ip",
			TLP:      3,
		},
	}

	// Run all analyzers over it with 1 minute timeout
	reports, err := c.Cortex.AnalyzeData(j, "1minute")
	if err != nil {
		emsg := tgbotapi.NewMessage(input.Chat.ID, "")
		emsg.ReplyToMessageID = input.MessageID
		emsg.Text = fmt.Sprintf("Cortex failed with an error: %s", err.Error())
		c.Bot.Send(emsg)
	}

	// Iterate over channel with reports and get taxonomies
	for m := range reports {
		if m.Status == "Failure" {
			continue
		}

		txs := m.Taxonomies()
		for _, t := range txs {
			msg := tgbotapi.NewMessage(input.Chat.ID, "")
			msg.ReplyToMessageID = input.MessageID
			msg.Text = fmt.Sprintf("%s:%s = %s", t.Namespace, t.Predicate, t.Value)
			c.Bot.Send(msg)
		}

		tr, _ := json.MarshalIndent(m, "", "  ")

		fb := tgbotapi.FileBytes{
			Name:  fmt.Sprintf("%s-%s.json", m.AnalyzerID, m.ID),
			Bytes: tr,
		}

		attachment := tgbotapi.NewDocumentUpload(input.Chat.ID, fb)
		attachment.ReplyToMessageID = input.MessageID
		c.Bot.Send(attachment)
	}
}
