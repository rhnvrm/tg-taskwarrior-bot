package main

import (
	"bytes"
	"log"
	"text/template"

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
		tg.NewKeyboardButton("/new"),
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

	if status == task.DoesNotExist {
		status = task.PendingRegistration
		m.task.Add(username)
	}

	if update.Message.IsCommand() {
		msg := tg.NewMessage(update.Message.Chat.ID, ":+1:")

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
			switch status {
			case task.Registered:
				t := template.Must(template.New("tmpl").Parse(taskList))
				tasks := m.api.GetTriggers(m.task.GetAPIKey(username))
				var txt bytes.Buffer
				err := t.Execute(&txt, tasks)
				if err != nil {
					log.Printf("err: %v", err)
				}
				msg.Text = txt.String()
			default:
				msg.Text = "You need to /register your api key with me first"
				msg.ReplyMarkup = preRegistrationKeyboard
			}
		case "new":
			switch status {
			case task.Registered:
				msg.Text = "What is the description of the task?"
				m.task.SetStatus(username, task.WaitForTaskDescription)
			default:
				msg.Text = "You need to /register your api key with me first"
				msg.ReplyMarkup = preRegistrationKeyboard
			}
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
	case task.WaitForTaskDescription:
		m.api.NewTrigger(m.task.GetAPIKey(username), update.Message.Text)
		m.task.SetStatus(username, task.Registered)
		msg.Text = "Added task."
		msg.ReplyMarkup = postRegistrationKeyboard
	}

	m.bot.Send(msg)
}
