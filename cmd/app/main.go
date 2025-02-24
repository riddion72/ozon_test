package main

import (
	"flag"
	"log"
	"log/slog"

	"github.com/riddion72/ozon_test/internal/app"
	"github.com/riddion72/ozon_test/internal/config"
	"github.com/riddion72/ozon_test/internal/logger"
)

func main() {
	configPath := flag.String("config", "config/config.yaml", "config file path")
	flag.Parse()

	cfg, err := config.ParseConfig(*configPath)
	if err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}

	logger.MustInit(cfg.Logger.Level)
	logger.Info("Starting application", slog.String("version", "1.0.0"))

	app := app.NewApp(cfg)

	app.Server.MustRun()
}
