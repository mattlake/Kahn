package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := NewModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
