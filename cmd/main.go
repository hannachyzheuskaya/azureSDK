package main

import (
	"log/slog"
	"os"
	"x/internal/app/apiserver"
	"x/internal/app/config"
)

const (
	jsonLogger = "json"
	textLogger = "text"
)

func main() {

	cfg := config.MustLoad()
	log := setupLogger(cfg.LoggerType)
	log.Info("starting azure-connector")

	if err := apiserver.Start(log, cfg); err != nil {
		log.Error("failed to start azure-connector")
	}
}

func setupLogger(typeLogger string) *slog.Logger {
	var createdLogger *slog.Logger

	switch typeLogger {
	case textLogger:
		createdLogger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case jsonLogger:
		createdLogger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return createdLogger
}
