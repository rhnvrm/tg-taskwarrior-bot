package main

func main() {
	cfg := readConfig("config.toml")
	manager := setupBot(cfg.Bot.Token, cfg.Bot.Debug)

	for update := range manager.drainAndListen() {
		manager.handle(update)
	}
}
