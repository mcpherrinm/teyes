package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(model{debug: true}, tea.WithMouseAllMotion(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("error: %s", err)
		os.Exit(1)
	}
}

type model struct {
	mouseX, mouseY, winWidth, winHeight int
	debug                               bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}

	case tea.MouseMsg:
		m.mouseX = msg.X
		m.mouseY = msg.Y

	case tea.WindowSizeMsg:
		m.winWidth = msg.Width
		m.winHeight = msg.Height
	}

	return m, nil
}

func (m model) View() string {
	var art strings.Builder

	for y := 0; y < m.winHeight; y++ {
		art.WriteString("\n")
		for x := 0; x < m.winWidth; x++ {
			art.WriteByte(m.eye(x, y))
		}
	}

	return art.String()
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func inCircle(x, y, centerX, centerY, radius int) bool {
	dY := abs(y - centerY)
	dX := abs(x - centerX)

	return dX*dX+dY*dY < radius*radius
}

// eye returns the character to render at an x,y point
func (m model) eye(x, y int) byte {
	// Debug mouse position:
	if m.debug {
		if x == m.mouseX && y == m.mouseY {
			return 'M'
		}
	}

	// handle symmetric left-right
	halfWidth := m.winWidth / 2
	if x >= halfWidth {
		x = x - halfWidth
	}

	eyeX := m.winWidth / 4
	eyeY := m.winHeight / 2
	eyeRadius := m.winHeight / 2

	// Pupil:
	padding := 3

	pupilXOff := m.mouseX - eyeX
	if pupilXOff > eyeRadius/padding {
		pupilXOff = eyeRadius / padding
	}
	if pupilXOff < -eyeRadius/padding {
		pupilXOff = -eyeRadius / padding
	}
	pupilYOff := m.mouseY - eyeY
	if pupilYOff > eyeRadius/padding {
		pupilYOff = eyeRadius / padding
	}
	if pupilYOff < -eyeRadius/padding {
		pupilYOff = -eyeRadius / padding
	}
	if inCircle(x, y, eyeX+pupilXOff, eyeY+pupilYOff, eyeRadius/3) {
		return 'X'
	}

	// Iris:
	if inCircle(x, y, eyeX, eyeY, eyeRadius-2) {
		return ' '
	}

	// border:
	if inCircle(x, y, eyeX, eyeY, eyeRadius) {
		return 'X'
	}

	// background:
	return ' '
}
