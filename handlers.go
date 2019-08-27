package main

import (
	"log"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rhnvrm/tgtask/task"
)

var preRegistrationKeyboard = tg.NewReplyKeyboard(
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton("/help"),
	),
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton("/register"),
	),
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton("/close"),
	),
)

var postRegistrationKeyboard = tg.NewReplyKeyboard(
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton("/help"),
		tg.NewKeyboardButton("/register"),
	),
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton("/list"),
	),
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton("/close"),
	),
)

func (m *manager) handle(update tg.Update) {
	if update.Message == nil {
		return
	}

	username := update.Message.From.UserName
	status := m.task.GetStatus(username)
	log.Println("current_status", status)

	if status == task.DoesNotExist {
		status = task.PendingRegistration
		m.task.Add(username)
	}

	if update.Message.IsCommand() {
		msg := tg.NewMessage(update.Message.Chat.ID, "Okay")

		switch update.Message.Command() {
		case "start":
			switch status {
			case task.Registered:
				msg.ReplyMarkup = postRegistrationKeyboard
			default:
				msg.ReplyMarkup = preRegistrationKeyboard
			}
		case "help":
			msg.Text = helpText
			msg.ParseMode = tg.ModeMarkdown
		case "register":
			msg.Text = "Send me your api key for https://inthe.am"
			m.task.SetStatus(username, task.WaitForAPIKey)
			msg.ReplyMarkup = tg.NewRemoveKeyboard(true)
		case "list":
			msg.Text = ""
		case "close":
			msg.ReplyMarkup = tg.NewRemoveKeyboard(true)
		default:
			msg.Text = couldNotUnderstand
		}

		m.bot.Send(msg)
		return
	}

	msg := tg.NewMessage(update.Message.Chat.ID, couldNotUnderstand)
	msg.ReplyToMessageID = update.Message.MessageID

	switch status {
	case task.WaitForAPIKey:
		m.task.Register(username, update.Message.Text)
		msg.Text = "Registered api key."
		msg.ReplyMarkup = postRegistrationKeyboard
	}

	m.bot.Send(msg)
}
