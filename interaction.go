package irpbot

import (
	"gopkg.in/telegram-bot-api.v4"
)

// editMessage struct used for editing message
type editMessage struct {
	ChatID    *int64
	MessageID *int
	Text      string
}

// interactionMain shows menu for the /interact
func interactionMain(id *int64) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(*id, "Choose what you want me to do")
	keyboard := tgbotapi.InlineKeyboardMarkup{}
	var row []tgbotapi.InlineKeyboardButton

	row = append(row, tgbotapi.NewInlineKeyboardButtonData("Hive", "HiveMenu"))
	row = append(row, tgbotapi.NewInlineKeyboardButtonData("Cortex", "CortexMenu"))
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

	msg.ReplyMarkup = keyboard

	return msg
}

// interactionHive shows menu for The Hive
func (e *editMessage) interactionHive() tgbotapi.EditMessageTextConfig {
	msg := tgbotapi.NewEditMessageText(*e.ChatID, *e.MessageID, "Choose what you want me to do")
	keyboard := tgbotapi.InlineKeyboardMarkup{}
	var row []tgbotapi.InlineKeyboardButton

	row = append(row, tgbotapi.NewInlineKeyboardButtonData("Cases", "HiveCases"))
	row = append(row, tgbotapi.NewInlineKeyboardButtonData("Tasks", "HiveTasks"))
	row = append(row, tgbotapi.NewInlineKeyboardButtonData("Observables", "HiveObservables"))
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

	msg.ReplyMarkup = &keyboard

	return msg
}

// interactionHiveTasks shows menu for the Hive's tasks
func (e *editMessage) interactionHiveTasks() tgbotapi.EditMessageTextConfig {
	msg := tgbotapi.NewEditMessageText(*e.ChatID, *e.MessageID, "Sorry! Tasks handler is under construction...")
	return msg
}

// interactionHiveCases shows menu for the Hive's cases info
func (e *editMessage) interactionHiveCases() tgbotapi.EditMessageTextConfig {
	msg := tgbotapi.NewEditMessageText(*e.ChatID, *e.MessageID, "Sorry! Cases handler is under construction...")
	return msg
}

// interactionHiveObservables shows menu for the Hive's cases info
func (e *editMessage) interactionHiveObservables() tgbotapi.EditMessageTextConfig {
	msg := tgbotapi.NewEditMessageText(*e.ChatID, *e.MessageID, "Sorry! Cases handler is under construction...")
	return msg
}

// interactionCortex shows menu for the Cortex
func (e *editMessage) interactionCortex() tgbotapi.EditMessageTextConfig {
	msg := tgbotapi.NewEditMessageText(*e.ChatID, *e.MessageID, "Choose what you want me to do")
	keyboard := tgbotapi.InlineKeyboardMarkup{}
	var row []tgbotapi.InlineKeyboardButton

	row = append(row, tgbotapi.NewInlineKeyboardButtonData("Start job", "CortexTasksAdd"))
	row = append(row, tgbotapi.NewInlineKeyboardButtonData("List jobs", "CortexTasksList"))
	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

	msg.ReplyMarkup = &keyboard

	return msg
}

// interactionCortexTasksAdd shows menu for adding Cortex jobs
func (e *editMessage) interactionCortexTasksAdd() tgbotapi.EditMessageTextConfig {
	msg := tgbotapi.NewEditMessageText(*e.ChatID, *e.MessageID, "Sorry! Cortex tasks handler is under construction...")
	return msg
}

// interactionCortexTasksList shows list of Cortex jobs
func (e *editMessage) interactionCortexTasksList() tgbotapi.EditMessageTextConfig {
	msg := tgbotapi.NewEditMessageText(*e.ChatID, *e.MessageID, "Sorry! Cortex tasks handler is under construction...")
	return msg
}
