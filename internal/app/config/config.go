package config

import (
	"github.com/BurntSushi/toml"
	"log"
	"time"
)

type Config struct {
	BindAddr    string        `toml:"bind_addr"`
	Timeout     time.Duration `toml:"timeout"`
	IdleTimeout time.Duration `toml:"idle_timeout"`
	LoggerType  string        `toml:"logger_type"`
	SessionKey  string        `toml:"session_key"`
}

func newConfig() *Config {
	return &Config{
		BindAddr:    ":8080",
		LoggerType:  "text",
		Timeout:     3 * time.Minute,
		IdleTimeout: 3 * time.Minute,
	}
}

func MustLoad() *Config {
	cfg := newConfig()
	_, err := toml.DecodeFile("configs/apiserver.toml", cfg)
	if err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return cfg
}
