package model

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

var Options = []tea.ProgramOption{tea.WithAltScreen(), tea.WithMouseAllMotion()}

type Model struct {
	mouseX, mouseY, winWidth, winHeight int
	debug                               bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "d":
			// Pressing d turns some debug output on
			m.debug = !m.debug
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

func (m Model) View() string {
	eyeLX := m.winWidth / 4
	eyeY := m.winHeight / 2
	eyeRadius := m.winHeight / 2
	if eyeRadius > m.winWidth/2 {
		eyeRadius = m.winWidth / 2
	}

	var pupilVecX, pupilVecY int
	if m.winWidth == 0 || m.winHeight == 0 {
		pupilVecX = 0
		pupilVecY = 0
	} else {
		pupilVecX = int((float32(m.mouseX-(m.winWidth/2)) / float32(m.winWidth/2)) * float32(eyeRadius/2))
		pupilVecY = int((float32(m.mouseY-(m.winHeight/2)) / float32(m.winHeight/2)) * float32(eyeRadius/2))
	}

	s := state{
		eyeLX:     eyeLX,
		eyeY:      eyeY,
		eyeRadius: eyeRadius,

		pupilVecX: pupilVecX,
		pupilVecY: pupilVecY,
	}

	var art strings.Builder
	if m.debug {
		art.WriteString(fmt.Sprintf("%d,%d in %dx%d: %+v", m.mouseX, m.mouseY, m.winWidth, m.winHeight, s))
	}
	art.WriteString("\n")

	for y := 1; y < m.winHeight-1; y++ {
		for x := 0; x < m.winWidth; x++ {
			art.WriteByte(m.eye(x, y, s))
		}
		art.WriteString("\n")
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

type state struct {
	eyeLX     int
	eyeY      int
	eyeRadius int

	pupilVecX int
	pupilVecY int
}

// eye returns the character to render at an x,y point
func (m Model) eye(x, y int, s state) byte {
	// Debug mouse position:
	if m.debug {
		if x == m.mouseX && y == m.mouseY {
			return 'M'
		}
	}

	// the eyes are symmetric
	if x > m.winWidth/2 {
		x = x - m.winWidth/2
	}

	// Outside the circle is background
	if !inCircle(x, y, s.eyeLX, s.eyeY, s.eyeRadius) {
		return ' '
	}

	// Pupil:
	if inCircle(x, y, s.eyeLX+s.pupilVecX, s.eyeY+s.pupilVecY, s.eyeRadius/3) {
		return 'X'
	}

	// Iris:
	if inCircle(x, y, s.eyeLX, s.eyeY, s.eyeRadius-2) {
		return ' '
	}

	// Border:
	return 'X'
}
