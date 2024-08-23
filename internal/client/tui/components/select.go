package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Sadere/gophkeeper/internal/client/tui/style"
)

type SelectModel struct {
	choices    []string
	focusIndex int
}

func NewSelectModel(choices []string) SelectModel {
	return SelectModel{
		choices: choices,
	}
}

func (m SelectModel) Init() tea.Cmd {
	return nil
}

func (m SelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.focusIndex > 0 {
				m.focusIndex--
			}

		case "down":
			if m.focusIndex < len(m.choices)-1 {
				m.focusIndex++
			}
		}
	}

	return m, nil
}

func (m SelectModel) Selected() int {
	return m.focusIndex
}

func (m SelectModel) View() string {
	var b strings.Builder

	for i, choice := range m.choices {
		cursor := " "
		if m.Selected() == i {
			cursor = ">"
		}

		// Render the row
		row := fmt.Sprintf("%s %s\n", cursor, choice)

		if m.Selected() == i {
			b.WriteString(style.FocusedStyle.Render(row))
		} else {
			b.WriteString(style.BlurredStyle.Render(row))
		}

		b.WriteString("\n")
	}

	return b.String()
}
