package main

import (
	"errors"
	"log"
	"os"

	"github.com/Sadere/gophkeeper/internal/server"
	"github.com/Sadere/gophkeeper/internal/server/config"
	"github.com/Sadere/gophkeeper/internal/server/utils"
	"github.com/spf13/pflag"
)

func main() {
	// Read config
	cfg, err := config.NewConfig(os.Args[1:])

	// Help requested, do nothing
	if errors.Is(err, pflag.ErrHelp) {
		return
	}

	if err != nil {
		log.Fatal("failed to initialize config: ", err)
	}

	// Init logs
	logger, err := utils.NewZapLogger(cfg.LogLevel)
	if err != nil {
		log.Fatal("failed to initialize zap logger: ", err)
	}

	// Create app instance
	app := server.NewApp(cfg, logger)

	// Run app
	err = app.Start()
	if err != nil {
		log.Fatal("failed to start app: ", err)
	}
}
