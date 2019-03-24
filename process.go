package cortexbot

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	valid "github.com/asaskevich/govalidator"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	cortex "github.com/ilyaglow/go-cortex/v2"
)

type updateMeta struct {
	from   *tgbotapi.User
	chatID int64
	msgID  int
	data   string
	typ    string
}

func newUpdateMeta(update *tgbotapi.Update) (*updateMeta, error) {
	if update.CallbackQuery != nil {
		return &updateMeta{
			from:   update.CallbackQuery.From,
			chatID: update.CallbackQuery.Message.Chat.ID,
			msgID:  update.CallbackQuery.Message.MessageID,
			data:   update.CallbackQuery.Data,
			typ:    "callback",
		}, nil
	}

	if update.Message != nil {
		return &updateMeta{
			from:   update.Message.From,
			chatID: update.Message.Chat.ID,
			msgID:  update.Message.MessageID,
			data:   update.Message.Text,
			typ:    "message",
		}, nil
	}

	return nil, fmt.Errorf("unknown update type: %v", update)
}

// sendReport sends a report depends on success
func (c *Cortexbot) sendReport(r *cortex.Report, callback *tgbotapi.CallbackQuery) error {
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

	_, err := c.Bot.Send(attachment)
	return err
}

func (c *Cortexbot) processUpdate(update *tgbotapi.Update) error {
	meta, err := newUpdateMeta(update)
	if err != nil {
		return err
	}

	log.Printf(
		"%s: [%s] %s",
		meta.typ,
		meta.from.UserName,
		meta.data,
	)

	msg := tgbotapi.NewMessage(meta.chatID, "")
	msg.ReplyToMessageID = meta.msgID

	isCmd := update.Message != nil && update.Message.IsCommand()
	var cmd string
	if isCmd {
		cmd = update.Message.Command()
	}

	if isCmd && cmd == "login" && c.noAdmin() {
		if update.Message.CommandArguments() != c.Password {
			msg.Text = "Wrong password"
			_, err = c.Bot.Send(msg)
			return err
		}

		u := &User{
			ID:    meta.from.ID,
			Admin: 1,
			About: meta.from.String(),
		}
		err = c.addUser(u)
		if err != nil {
			msg.Text = fmt.Sprintf("Can't register: %s", err.Error())
			_, err = c.Bot.Send(msg)
			return err
		}
		log.Printf("Registered new user %d", meta.from.ID)
		return nil
	}

	// user initiates interaction.
	if isCmd && cmd == "start" {
		if c.noAdmin() {
			msg.Text = "No admin assigned. Try /login password"
			_, err = c.Bot.Send(msg)
			return err
		}

		if c.CheckAuth(meta.from) {
			msg.Text = "Already logged in, you can send indicators here"
			_, err := c.Bot.Send(msg)
			return err
		}

		msg.Text = "Forward the next message to the bot admin"
		_, err := c.Bot.Send(msg)
		if err != nil {
			return err
		}
		msg = tgbotapi.NewMessage(meta.chatID, "")
		msg.Text = fmt.Sprintf("/approve %d %s", meta.from.ID, meta.from.String())
		_, err = c.Bot.Send(msg)
		return err
	}

	if c.CheckAdmin(meta.from) && isCmd && cmd == "approve" {
		parts := strings.SplitN(update.Message.Text, " ", 3)
		if len(parts) < 3 {
			msg.Text = fmt.Sprintf("Not enough parameters to approve: %s", update.Message.Text)
			_, err = c.Bot.Send(msg)
			return err
		}
		id, err := strconv.ParseInt(parts[1], 10, 32)
		u := &User{
			ID:    int(id),
			Admin: 0,
			About: parts[2],
		}
		err = c.addUser(u)
		if err != nil {
			msg.Text = fmt.Sprintf("can't create user: %s", err.Error())
			_, err = c.Bot.Send(msg)
			if err != nil {
				return err
			}
		}
		msg.Text = "Successfully approved"
		_, err = c.Bot.Send(msg)
		return err
	}

	// auth guard
	if !c.CheckAuth(meta.from) {
		msg.Text = "Non-authorized action, type /start to login"
		_, err := c.Bot.Send(msg)
		return err
	}

	switch meta.typ {
	case "callback":
		if err := c.processCallback(update.CallbackQuery); err != nil {
			return err
		}

		cbcfg := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
		_, err := c.Bot.AnswerCallbackQuery(cbcfg)
		return err
	case "message":
		return c.processMessage(update.Message)
	}

	return nil
}

// processCallback analyzes observables with a selected set of analyzers
func (c *Cortexbot) processCallback(callback *tgbotapi.CallbackQuery) error {
	var j cortex.Observable
	var err error

	if callback.Message.ReplyToMessage.Document != nil {
		link, err := c.Bot.GetFileDirectURL(callback.Message.ReplyToMessage.Document.FileID)
		if err != nil {
			log.Println("Can't get direct link to file")
			return err
		}

		j, err = newFileArtifactFromURL(link, callback.Message.ReplyToMessage.Document.FileName, c.TLP, c.PAP, c.Bot.Client)
		if err != nil {
			log.Println(err.Error())
			return err
		}
	} else {
		j = newArtifact(callback.Message.ReplyToMessage.Text, c.TLP, c.PAP)
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
		edit := tgbotapi.NewEditMessageReplyMarkup(callback.Message.Chat.ID, callback.Message.MessageID, *kb)
		if _, err := c.Bot.Send(edit); err != nil {
			return err
		}
	case "show":
		kb, err := c.analyzersButtons(dataType(callback.Message.ReplyToMessage.Text))
		if err != nil {
			return err
		}
		edit := tgbotapi.NewEditMessageReplyMarkup(callback.Message.Chat.ID, callback.Message.MessageID, *kb)
		if _, err := c.Bot.Send(edit); err != nil {
			return err
		}
	default:
		r, err := c.Cortex.Analyzers.Run(context.Background(), callback.Data, j, c.Timeout)
		if err != nil {
			msg := tgbotapi.NewMessage(callback.Message.Chat.ID, fmt.Sprintf("%s failed: %s", callback.Data, err.Error()))
			msg.ReplyToMessageID = callback.Message.MessageID
			if _, err := c.Bot.Send(msg); err != nil {
				return err
			}
		} else {
			err = c.sendReport(r, callback)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// analyzersButtons returns a markup of analyzers as buttons
func (c *Cortexbot) analyzersButtons(datatype string) (*tgbotapi.InlineKeyboardMarkup, error) {
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
func (c *Cortexbot) processMessage(input *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(input.Chat.ID, "Select analyzer to run. Choose <All> to run all of them.")
	msg.ReplyToMessageID = input.MessageID

	var dt string

	if input.Document != nil {
		dt = "file"
	} else {
		dt = dataType(input.Text)
	}

	bmarkup, err := c.analyzersButtons(dt)
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
func newArtifact(s string, tlp cortex.TLP, pap cortex.PAP) cortex.Observable {
	return &cortex.Task{
		Data:     s,
		DataType: dataType(s),
		TLP:      &tlp,
		PAP:      &pap,
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
func newFileArtifactFromURL(link, fname string, tlp cortex.TLP, pap cortex.PAP, client *http.Client) (cortex.Observable, error) {
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
			PAP:      &pap,
		},
		FileName: fname,
		Reader:   resp.Body,
	}, nil
}
