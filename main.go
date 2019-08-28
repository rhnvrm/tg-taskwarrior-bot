package main

import (
	"time"

	"github.com/gammazero/workerpool"
)

func main() {
	wp := workerpool.New(20)
	cfg := readConfig("config.toml")
	manager := setupBot(cfg.Bot.Token, cfg.Bot.Debug)

	dumpTimer := time.NewTicker(30 * time.Second)
	go func() {
		for range dumpTimer.C {
			manager.task.Dump()
		}
	}()

	for update := range manager.drainAndListen() {
		wp.Submit(func() {
			manager.handle(update)
		})
	}
}
