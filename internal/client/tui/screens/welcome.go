package screens

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type WelcomeModel struct {
	state   *State
	choices []string
	cursor  int
}

func NewWelcomeModel(state *State) *WelcomeModel {
	return &WelcomeModel{
		state:   state,
		choices: []string{"Login", "Register new user"},
	}
}

func (m WelcomeModel) Init() tea.Cmd {
	return tea.SetWindowTitle("GophKeeper client")
}

func (m *WelcomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			switch m.cursor {
			case 0:
				loginModel := NewLoginModel(m.state)
				return NewRootModel(m.state).SwitchScreen(loginModel)
			case 1:
				registerModel := NewRegisterModel(m.state)
				return NewRootModel(m.state).SwitchScreen(registerModel)
			}
		}
	}

	return m, nil
}

func (m *WelcomeModel) View() string {
	var b strings.Builder

	b.WriteString("What would you like to do?\n\n")

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		// Render the row
		fmt.Fprintf(&b, "%s %s\n", cursor, choice)
	}

	return b.String()
}
