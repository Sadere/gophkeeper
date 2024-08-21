package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Sadere/gophkeeper/internal/client"
	"github.com/Sadere/gophkeeper/internal/client/config"
)

func main() {
	cfg, err := config.NewConfig(os.Args[1:])
	if err != nil {
		log.Fatal("failed to initialize config: ", err)
	}

	app, err := client.NewKeeperClient(cfg)
	if err != nil {
		log.Fatal("failed to initialize app: ", err)
	}

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal("failed to setup bubbletea debug log: ", err)
	}

	defer f.Close()

	p := tea.NewProgram(app.Root, tea.WithAltScreen())

	_, err = p.Run()
	if err != nil {
		log.Fatal("failed to run bubbletea program: ", err)
	}
}
