package main

import (
	"fmt"
	"os"

	"github.com/mcpherrinm/teyes/model"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(model.Model{}, tea.WithMouseAllMotion(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("error: %s", err)
		os.Exit(1)
	}
}
