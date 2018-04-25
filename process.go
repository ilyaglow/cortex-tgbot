package cortexbot

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	valid "github.com/asaskevich/govalidator"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ilyaglow/go-cortex"
)

// ProcessCortex asks Cortex about data submitted by a user
func (c *Client) ProcessCortex(input *tgbotapi.Message) error {
	var j cortex.Observable
	var err error

	if input.Document != nil {
		link, err := c.Bot.GetFileDirectURL(input.Document.FileID)
		if err != nil {
			log.Println("Can't get direct link to file")
			return err
		}

		j, err = newFileArtifactFromURL(link, c.TLP)
		if err != nil {
			log.Println(err.Error())
			return err
		}
	} else {
		j, err = newArtifact(input.Text, c.TLP)
		if err != nil {
			log.Println(err.Error())
			return err
		}
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
			log.Printf("Analyzer %s failed with error message: %s", m.AnalyzerID, m.Report.ErrorMessage)
			continue
		}

		// Send JSON file with full report and taxonomies
		tr, _ := json.MarshalIndent(m, "", "  ")

		fb := tgbotapi.FileBytes{
			Name:  fmt.Sprintf("%s-%s.json", m.AnalyzerID, m.ID),
			Bytes: tr,
		}

		attachment := tgbotapi.NewDocumentUpload(input.Chat.ID, fb)
		attachment.ReplyToMessageID = input.MessageID
		attachment.Caption = buildTaxonomies(m.Taxonomies())
		c.Bot.Send(attachment)
	}

	return nil
}

// buildTaxonomies joins taxonomies in one formatted string
// Every taxonomy is separated with two spaces from each other
func buildTaxonomies(txs []cortex.Taxonomy) string {
	var stats []string

	for _, t := range txs {
		stats = append(stats, fmt.Sprintf("%s:%s = %s", t.Namespace, t.Predicate, t.Value))

	}
	return strings.Join(stats, ", ")
}

// newArtifact makes an Artifact depends on its type
func newArtifact(s string, tlp int) (cortex.Observable, error) {
	var dataType string

	if valid.IsIP(s) {
		dataType = "ip"
	} else if IsDNSName(s) {
		dataType = "domain"
	} else if IsHash(s) {
		dataType = "hash"
	} else if valid.IsURL(s) {
		dataType = "url"
	} else if valid.IsEmail(s) {
		dataType = "mail"
	} else {
		dataType = "unknown"
		return nil, errors.New("Unknown data type")
	}

	j := &cortex.Artifact{
		Data: s,
		Attributes: cortex.ArtifactAttributes{
			DataType: dataType,
			TLP:      tlp,
		},
	}

	return j, nil
}

// newFileArtifactFromURL makes a FileArtifact from URL
func newFileArtifactFromURL(link string, tlp int) (cortex.Observable, error) {
	resp, err := http.Get(link)
	if err != nil {
		return nil, err
	}

	return &cortex.FileArtifact{
		FileArtifactMeta: cortex.FileArtifactMeta{
			DataType: "file",
			TLP:      tlp,
		},
		FileName: link,
		Reader:   resp.Body,
	}, nil
}
