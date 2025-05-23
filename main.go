package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"

	"capuchinator/internal/application"
)

func main() {
	app, err := application.New()
	if err != nil {
		log.Fatalf("New: %v\n", err)
	}

	p := tea.NewProgram(app, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Run: %v\n", err)
	}
}
