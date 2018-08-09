package cortexbot

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/ilyaglow/go-cortex"
	"github.com/ilyaglow/telegram-bot-api"
)

// sendReport sends a report depends on success
func (c *Client) sendReport(r *cortex.Report, callback *tgbotapi.CallbackQuery) error {
	if r.Status == "Failure" {
		return fmt.Errorf("Analyzer %s failed with error message: %s", r.AnalyzerName, r.ReportBody.ErrorMessage)
	}

	// Send JSON file with full report and taxonomies
	tr, _ := json.MarshalIndent(r, "", "  ")

	fb := tgbotapi.FileBytes{
		Name:  fmt.Sprintf("%s-%s.json", r.AnalyzerName, r.ID),
		Bytes: tr,
	}

	attachment := tgbotapi.NewDocumentUpload(callback.Message.Chat.ID, fb)
	attachment.ReplyToMessageID = callback.Message.ReplyToMessage.MessageID
	attachment.Caption = buildTaxonomies(r.Taxonomies())
	go c.Bot.Send(attachment)

	return nil
}

// processCallback analyzes observables with a selected set of analyzers
func (c *Client) processCallback(callback *tgbotapi.CallbackQuery) error {
	var j cortex.Observable
	var err error

	if callback.Message.ReplyToMessage.Document != nil {
		link, err := c.Bot.GetFileDirectURL(callback.Message.ReplyToMessage.Document.FileID)
		if err != nil {
			log.Println("Can't get direct link to file")
			return err
		}

		j, err = newFileArtifactFromURL(link, callback.Message.ReplyToMessage.Document.FileName, c.TLP, c.Bot.Client)
		if err != nil {
			log.Println(err.Error())
			return err
		}
	} else {
		j = newArtifact(callback.Message.ReplyToMessage.Text, c.TLP)
	}

	switch callback.Data {
	case "all":
		mul := c.Cortex.Analyzers.NewMultiRun(context.Background(), 5*time.Minute)
		mul.OnReport = func(r *cortex.Report) {
			err = c.sendReport(r, callback)
			log.Println(err)
		}

		mul.OnError = func(e error, o cortex.Observable, a *cortex.Analyzer) {
			log.Println(fmt.Sprintf("Cortex analyzer %s failed on data %s with an error: %s", a.Name, o.Description(), e.Error()))
		}

		err = mul.Do(j)
		if err != nil {
			return err
		}
	case "close":
		kb := showButton()
		edit := tgbotapi.NewEditMessageReplyMarkup(callback.Message.Chat.ID, callback.Message.MessageID, *kb)
		go c.Bot.Send(edit)
	case "show":
		kb, err := c.analyzersButtons(dataType(callback.Message.ReplyToMessage.Text))
		if err != nil {
			return err
		}
		edit := tgbotapi.NewEditMessageReplyMarkup(callback.Message.Chat.ID, callback.Message.MessageID, *kb)
		go c.Bot.Send(edit)
	default:
		r, err := c.Cortex.Analyzers.Run(context.Background(), callback.Data, j, time.Minute*5)
		if err != nil {
			msg := tgbotapi.NewMessage(callback.Message.Chat.ID, fmt.Sprintf("%s failed: %s", callback.Data, err.Error()))
			msg.ReplyToMessageID = callback.Message.MessageID
			c.Bot.Send(msg)
		} else {
			err = c.sendReport(r, callback)
			if err != nil {
				return err
			}
		}
	}

	cbcfg := tgbotapi.NewCallback(callback.ID, "")
	go c.Bot.AnswerCallbackQuery(cbcfg)

	return nil
}

// analyzersButtons returns a markup of analyzers as buttons
func (c *Client) analyzersButtons(datatype string) (*tgbotapi.InlineKeyboardMarkup, error) {
	analyzers, _, err := c.Cortex.Analyzers.ListByType(context.Background(), datatype)
	if err != nil {
		return nil, err
	}

	var names []string
	for i := range analyzers {
		names = append(names, analyzers[i].Name)
	}
	sort.Strings(names)

	var rows [][]tgbotapi.InlineKeyboardButton
	rows = append(rows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("All", "all"),
	})

	for _, n := range names {
		var buttons []tgbotapi.InlineKeyboardButton
		b := tgbotapi.NewInlineKeyboardButtonData(n, n)
		buttons = append(buttons, b)
		rows = append(rows, buttons)
	}

	rows = append(rows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Close", "close"),
	})

	markup := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &markup, nil
}

// showButton returns only one button, used on close callback
func showButton() *tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	rows = append(rows, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Show", "show"),
	})

	markup := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &markup
}

// processMessage asks Cortex about data submitted by a user
func (c *Client) processMessage(input *tgbotapi.Message) error {
	var err error

	msg := tgbotapi.NewMessage(input.Chat.ID, "Select analyzer to run. Choose <All> to run all of them.")
	msg.ReplyToMessageID = input.MessageID

	bmarkup := &tgbotapi.InlineKeyboardMarkup{}
	var dt string

	if input.Document != nil {
		dt = "file"
	} else {
		dt = dataType(input.Text)
	}

	bmarkup, err = c.analyzersButtons(dt)
	if err != nil {
		return err
	}
	msg.ReplyMarkup = bmarkup

	_, err = c.Bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

// buildTaxonomies joins taxonomies in one formatted string
// Every taxonomy is separated with two spaces from each other
func buildTaxonomies(txs []cortex.Taxonomy) string {
	var stats []string

	for _, t := range txs {
		stats = append(stats, fmt.Sprintf("%s:%s = %v", t.Namespace, t.Predicate, t.Value))
	}
	return strings.Join(stats, "\n")
}

// newArtifact makes an Artifact depends on its type
func newArtifact(s string, tlp int) cortex.Observable {
	return &cortex.Task{
		Data:     s,
		DataType: dataType(s),
		TLP:      &tlp,
	}
}

func dataType(d string) (t string) {
	if valid.IsIP(d) {
		t = "ip"
	} else if IsDNSName(d) {
		t = "domain"
	} else if IsHash(d) {
		t = "hash"
	} else if valid.IsEmail(d) {
		t = "mail"
	} else if valid.IsURL(d) {
		t = "url"
	} else {
		t = "other"
	}

	return t
}

// newFileArtifactFromURL makes a FileArtifact from URL
func newFileArtifactFromURL(link string, fname string, tlp int, client *http.Client) (cortex.Observable, error) {
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return &cortex.FileTask{
		FileTaskMeta: cortex.FileTaskMeta{
			DataType: "file",
			TLP:      &tlp,
		},
		FileName: fname,
		Reader:   resp.Body,
	}, nil
}
