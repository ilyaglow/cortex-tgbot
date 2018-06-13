package cortexbot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/ilyaglow/telegram-bot-api"
	"gopkg.in/ilyaglow/go-cortex.v2"
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

		j, err = newFileArtifactFromURL(link, input.Document.FileName, c.TLP, c.Bot.Client)
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
	mul := c.Cortex.Analyzers.NewMultiRun(context.Background(), 5*time.Minute)
	mul.OnReport = func(r *cortex.Report) {
		if r.Status == "Failure" {
			log.Printf("Analyzer %s failed with error message: %s", r.AnalyzerName, r.ReportBody.ErrorMessage)
			return
		}

		// Send JSON file with full report and taxonomies
		tr, _ := json.MarshalIndent(r, "", "  ")

		fb := tgbotapi.FileBytes{
			Name:  fmt.Sprintf("%s-%s.json", r.AnalyzerName, r.ID),
			Bytes: tr,
		}

		attachment := tgbotapi.NewDocumentUpload(input.Chat.ID, fb)
		attachment.ReplyToMessageID = input.MessageID
		attachment.Caption = buildTaxonomies(r.Taxonomies())
		c.Bot.Send(attachment)

	}

	mul.OnError = func(e error, o cortex.Observable, a *cortex.Analyzer) {
		log.Println(fmt.Sprintf("Cortex analyzer %s failed on data %s with an error: %s", a.Name, o.Description(), e.Error()))
	}

	err = mul.Do(j)
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
	} else if valid.IsEmail(s) {
		dataType = "mail"
	} else if valid.IsURL(s) {
		dataType = "url"
	} else {
		dataType = "unknown"
		return nil, errors.New("Unknown data type")
	}

	j := &cortex.Task{
		Data:     s,
		DataType: dataType,
		TLP:      &tlp,
	}

	return j, nil
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
