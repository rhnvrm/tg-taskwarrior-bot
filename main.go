package main

import (
	"github.com/gammazero/workerpool"
)

func main() {
	wp := workerpool.New(20)
	cfg := readConfig("config.toml")
	manager := setupBot(cfg.Bot.Token, cfg.Bot.Debug)

	for update := range manager.drainAndListen() {
		wp.Submit(func() {
			manager.handle(update)
		})
	}
}
