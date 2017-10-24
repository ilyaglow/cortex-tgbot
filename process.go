package cortexbot

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	valid "github.com/asaskevich/govalidator"
	"github.com/ilyaglow/go-cortex"
	"gopkg.in/telegram-bot-api.v4"
)

// ProcessCortex asks Cortex about data submitted by a user
func (c *Client) ProcessCortex(input *tgbotapi.Message) error {
	j, err := c.constructJob(input.Text)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// Run all analyzers over it with 1 minute timeout
	reports, err := c.Cortex.AnalyzeData(j, "1minute")
	if err != nil {
		msg := tgbotapi.NewMessage(input.Chat.ID, "")
		msg.ReplyToMessageID = input.MessageID
		msg.Text = fmt.Sprintf("Cortex failed with an error: %s", err.Error())
		c.Bot.Send(msg)
	}

	// Iterate over channel with reports and get taxonomies
	for m := range reports {
		if m.Status == "Failure" {
			continue
		}

		txs := m.Taxonomies()
		// Send every taxonomy as a message
		for _, t := range txs {
			msg := tgbotapi.NewMessage(input.Chat.ID, "")
			msg.ReplyToMessageID = input.MessageID
			msg.Text = fmt.Sprintf("%s:%s = %s", t.Namespace, t.Predicate, t.Value)
			c.Bot.Send(msg)
		}

		// Send JSON file with full report
		tr, _ := json.MarshalIndent(m, "", "  ")

		fb := tgbotapi.FileBytes{
			Name:  fmt.Sprintf("%s-%s.json", m.AnalyzerID, m.ID),
			Bytes: tr,
		}

		attachment := tgbotapi.NewDocumentUpload(input.Chat.ID, fb)
		attachment.ReplyToMessageID = input.MessageID
		c.Bot.Send(attachment)
	}

	return nil
}

// constructJob make a JobBody depends on its type
func (c *Client) constructJob(s string) (*gocortex.JobBody, error) {
	var dataType string

	if valid.IsIP(s) {
		dataType = "ip"
	} else if IsDNSName(s) {
		dataType = "domain"
	} else if IsHash(s) {
		dataType = "hash"
	} else if valid.IsURL(s) {
		dataType = "url"
	} else {
		dataType = "unknown"
		return nil, errors.New("Unknown data type")
	}

	j := &gocortex.JobBody{
		Data: s,
		Attributes: gocortex.ArtifactAttributes{
			DataType: dataType,
			TLP:      3,
		},
	}

	return j, nil
}
