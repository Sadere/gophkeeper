package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Sadere/gophkeeper/internal/client/tui/style"
	"github.com/charmbracelet/bubbles/textinput"
)

type InputGroup struct {
	Inputs     []textinput.Model
	InputNum   int
	FocusIndex int
}

func NewInputGroup(inputs []textinput.Model) InputGroup {
	// Set styles
	for i, input := range inputs {
		input.Cursor.Style = style.FocusedStyle
		input.PromptStyle = style.FocusedStyle
		input.TextStyle = style.FocusedStyle

		inputs[i] = input
	}

	return InputGroup{
		Inputs:   inputs,
		InputNum: len(inputs),
	}
}

func (m InputGroup) Init() tea.Cmd {
	return textinput.Blink
}

func (m InputGroup) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.FocusIndex--
			} else {
				m.FocusIndex++
			}

			if m.FocusIndex > m.InputNum {
				m.FocusIndex = 0
			} else if m.FocusIndex < 0 {
				m.FocusIndex = m.InputNum
			}

			cmds := make([]tea.Cmd, m.InputNum)
			for i := 0; i <= m.InputNum-1; i++ {
				if i == m.FocusIndex {
					// Set focused state
					cmds[i] = m.Inputs[i].Focus()
					m.Inputs[i].PromptStyle = style.FocusedStyle
					m.Inputs[i].TextStyle = style.FocusedStyle
					continue
				}
				// Remove focused state
				m.Inputs[i].Blur()
				m.Inputs[i].PromptStyle = style.NoStyle
				m.Inputs[i].TextStyle = style.NoStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *InputGroup) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, m.InputNum)

	for i := range m.Inputs {
		m.Inputs[i], cmds[i] = m.Inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m InputGroup) View() string {
	var b strings.Builder

	for i := range m.Inputs {
		b.WriteString(m.Inputs[i].View())
		if i < m.InputNum-1 {
			b.WriteRune('\n')
		}
	}

	button := style.BlurredStyle.Render("[ Submit ]")
	if m.FocusIndex == m.InputNum {
		button = style.FocusedStyle.Render("[ Submit ]")
	}
	fmt.Fprintf(&b, "\n\n%s\n", button)

	return b.String()
}
