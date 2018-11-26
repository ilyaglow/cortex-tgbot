package cortexbot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"

	valid "github.com/asaskevich/govalidator"
	"github.com/ilyaglow/go-cortex"
	tb "gopkg.in/tucnak/telebot.v2"
)

// sendReport sends a report depends on success
func (c *Client) sendReport(r *cortex.Report, callback *tb.Callback) error {
	if r.Status == "Failure" {
		return fmt.Errorf("Analyzer %s failed with error message: %s", r.AnalyzerName, r.ReportBody.ErrorMessage)
	}

	// Send JSON file with full report and taxonomies
	tr, _ := json.MarshalIndent(r, "", "  ")

	fb := &tb.Document{
		File:     tb.FromReader(bytes.NewReader(tr)),
		FileName: fmt.Sprintf("%s-%s.json", r.AnalyzerName, r.ID),
		Caption:  buildTaxonomies(r.Taxonomies()),
	}

	_, err := c.Bot.Reply(callback.Message.ReplyTo, fb)
	return err
}

// processCallback analyzes observables with a selected set of analyzers
func (c *Client) processCallback(callback *tb.Callback) error {
	var j cortex.Observable
	var err error

	if callback.Message.ReplyTo.Document != nil {
		link, err := c.Bot.FileURLByID(callback.Message.ReplyTo.Document.FileID)
		if err != nil {
			log.Println("Can't get direct link to file")
			return err
		}

		j, err = newFileArtifactFromURL(link, callback.Message.ReplyTo.Document.FileName, c.TLP, c.BotSettings.Client)
		if err != nil {
			log.Println(err.Error())
			return err
		}
	} else {
		j = newArtifact(callback.Message.ReplyTo.Text, c.TLP)
	}

	switch callback.Data {
	case "all":
		mul := c.Cortex.Analyzers.NewMultiRun(context.Background(), c.Timeout)
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
		if _, err := c.Bot.Edit(callback.Message, kb); err != nil {
			return err
		}
	case "show":
		kb, err := c.analyzersButtons(dataType(callback.Message.ReplyTo.Text))
		if err != nil {
			return err
		}
		if _, err := c.Bot.Edit(callback.Message, kb); err != nil {
			return err
		}
	default:
		r, err := c.Cortex.Analyzers.Run(context.Background(), callback.Data, j, c.Timeout)
		if err != nil {
			return fmt.Errorf("%s failed: %s", callback.Data, err.Error())
		}

		err = c.sendReport(r, callback)
		if err != nil {
			return err
		}
	}

	return nil
}

// analyzersButtons returns a markup of analyzers as buttons
func (c *Client) analyzersButtons(datatype string) (*tb.InlineKeyboardMarkup, error) {
	analyzers, _, err := c.Cortex.Analyzers.ListByType(context.Background(), datatype)
	if err != nil {
		return nil, err
	}

	var names []string
	for i := range analyzers {
		names = append(names, analyzers[i].Name)
	}
	sort.Strings(names)

	var rows [][]tb.InlineButton
	rows = append(rows, []tb.InlineButton{
		tb.InlineButton{
			Text: "All",
			Data: "all",
		},
	})

	for _, n := range names {
		rows = append(rows, []tb.InlineButton{
			tb.InlineButton{
				Text: n,
				Data: n,
			},
		})
	}

	// rows = append(rows, []tb.InlineButton{
	// 	tb.InlineButton{
	// 		Text: "Close",
	// 		Data: "close",
	// 	},
	// })

	return &tb.InlineKeyboardMarkup{rows}, nil
}

// showButton returns only one button, used on close callback
func showButton() *tb.InlineKeyboardMarkup {
	var rows [][]tb.InlineButton
	return &tb.InlineKeyboardMarkup{
		InlineKeyboard: append(rows, []tb.InlineButton{
			tb.InlineButton{
				Text: "Show",
				Data: "show",
			},
		}),
	}
}

// processMessage asks Cortex about data submitted by a user
func (c *Client) processMessage(input *tb.Message) error {
	var (
		err error
		dt  string
	)

	bmarkup := &tb.InlineKeyboardMarkup{}

	if input.Document != nil {
		dt = "file"
	} else {
		dt = dataType(input.Text)
	}

	bmarkup, err = c.analyzersButtons(dt)
	if err != nil {
		return err
	}

	_, err = c.Bot.Reply(input, "Select analyzer to run. Choose <All> to run all of them.", &tb.ReplyMarkup{
		InlineKeyboard: bmarkup.InlineKeyboard,
	})
	return err
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

func (c *Client) availableDataTypes() (*tb.InlineKeyboardMarkup, error) {
	dataTypes, err := c.Cortex.Analyzers.DataTypes(context.Background())
	if err != nil {
		return nil, err
	}
	sort.Strings(dataTypes)

	var rows [][]tb.InlineButton
	for _, n := range dataTypes {
		rows = append(rows, []tb.InlineButton{
			tb.InlineButton{
				Text: n,
				Data: n,
			},
		})
	}

	rows = append(rows, []tb.InlineButton{
		tb.InlineButton{
			Text: "Close",
			Data: "close",
		},
	})

	return &tb.InlineKeyboardMarkup{rows}, nil
}
