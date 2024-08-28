package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

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

	// Catch quit signals
	quit := make(chan os.Signal, 1)

	p := tea.NewProgram(app.Root, tea.WithAltScreen())

	// Run tea program
	go func() {
		_, err = p.Run()
		if err != nil {
			log.Fatal("failed to run bubbletea program: ", err)
		}

		quit <- syscall.SIGINT
	}()

	// Run notification monitor
	go app.Client.Notifications(p)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-quit

	log.Println("client shutdown ...")
}
