package main

import (
	"log"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rhnvrm/tgtask/task"
)

type manager struct {
	bot  *tg.BotAPI
	task *task.Manager
}

func setupBot(token string, debug bool) manager {
	bot, err := tg.NewBotAPI(token)
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	bot.Debug = debug

	task, err := task.NewManager()
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	return manager{bot: bot, task: task}
}

func (m *manager) drainAndListen() tg.UpdatesChannel {
	u := tg.NewUpdate(0)
	u.Timeout = 60

	updates, err := m.bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	// wait for updates and clear them as we don't want to handle
	// a large backlog of old messages.
	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	return updates
}
