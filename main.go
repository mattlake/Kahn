package main

import (
	"kahn/internal/app"
	"kahn/internal/config"
	"kahn/internal/database"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	Version = "dev" // Set during build
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	database, err := database.NewDatabase(config)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	m := app.NewKahnModel(database, Version)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
