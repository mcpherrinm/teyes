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
	eyeRX := eyeLX + m.winWidth/2
	eyeY := m.winHeight / 2
	eyeRadius := m.winHeight / 2
	if eyeRadius > m.winWidth/2 {
		eyeRadius = m.winWidth / 2
	}

	pupilVecLX, pupilVecLY := vec(eyeLX, eyeY, m.mouseX, m.mouseY)
	pupilVecLX, pupilVecLY = clamp(eyeRadius/3, pupilVecLX, pupilVecLY)
	pupilVecRX, pupilVecRY := vec(eyeRX, eyeY, m.mouseX, m.mouseY)
	pupilVecRX, pupilVecRY = clamp(eyeRadius/3, pupilVecRX, pupilVecRY)

	s := state{
		eyeLX:     eyeLX,
		eyeRX:     eyeRX,
		eyeY:      eyeY,
		eyeRadius: eyeRadius,

		pupilVecLX: pupilVecLX,
		pupilVecLY: pupilVecLY,
		pupilVecRX: pupilVecRX,
		pupilVecRY: pupilVecRY,
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

func vec(fromX, fromY, toX, toY int) (int, int) {
	return toX - fromX, toY - fromY
}

// clamp magnitude of vector
func clamp(maxMag, x, y int) (int, int) {
	// TODO: this should do a vector thing
	if x > maxMag {
		x = maxMag
	}
	if x < -maxMag {
		x = -maxMag
	}
	if y > maxMag {
		y = maxMag
	}
	if y < -maxMag {
		y = -maxMag
	}
	return x, y
}

type state struct {
	eyeLX     int
	eyeRX     int
	eyeY      int
	eyeRadius int

	pupilVecLX int
	pupilVecLY int
	pupilVecRX int
	pupilVecRY int
}

// eye returns the character to render at an x,y point
func (m Model) eye(x, y int, s state) byte {
	// Debug mouse position:
	if m.debug {
		if x == m.mouseX && y == m.mouseY {
			return 'M'
		}
	}

	// Left pupil:
	if inCircle(x, y, s.eyeLX+s.pupilVecLX, s.eyeY+s.pupilVecLY, s.eyeRadius/3) {
		return 'X'
	}

	// Right pupil:
	if inCircle(x, y, s.eyeRX+s.pupilVecRX, s.eyeY+s.pupilVecRY, s.eyeRadius/3) {
		return 'X'
	}

	// the rest of the eye is symmetric
	if x > m.winWidth/2 {
		x = x - m.winWidth/2
	}

	// Iris:
	if inCircle(x, y, s.eyeLX, s.eyeY, s.eyeRadius-2) {
		return ' '
	}

	// border:
	if inCircle(x, y, s.eyeLX, s.eyeY, s.eyeRadius) {
		return 'X'
	}

	// background:
	return ' '
}
