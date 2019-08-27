package main

import (
	"log"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
)

type config struct {
	Bot botCfg `koanf:"bot"`
}

type botCfg struct {
	Token string `koanf:"token"`
	Debug bool   `koanf:"debug"`
}

func readConfig(fname string) config {
	var (
		cfg = config{}
		k   = koanf.New(".")

		err error
	)

	err = k.Load(file.Provider(fname), toml.Parser())
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	err = k.Unmarshal("", &cfg)
	if err != nil {
		log.Fatalf("err: %v", err)
	}

	return cfg
}
