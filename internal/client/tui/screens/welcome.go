package screens

import (
	"strings"

	"github.com/Sadere/gophkeeper/internal/client/tui/components"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	selectLogin = iota
	selectRegister
)

type WelcomeModel struct {
	state       *State
	selectModel components.SelectModel
}

func NewWelcomeModel(state *State) *WelcomeModel {
	choices := []string{
		selectLogin:    "Login",
		selectRegister: "Register new user",
	}

	return &WelcomeModel{
		state:       state,
		selectModel: components.NewSelectModel(choices),
	}
}

func (m WelcomeModel) Init() tea.Cmd {
	return tea.SetWindowTitle("GophKeeper client")
}

func (m *WelcomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			switch m.selectModel.Selected() {
			case selectLogin:
				loginModel := NewLoginModel(m.state)
				return NewRootModel(m.state).SwitchScreen(loginModel)
			case selectRegister:
				registerModel := NewRegisterModel(m.state)
				return NewRootModel(m.state).SwitchScreen(registerModel)
			}
		}
	}

	// Update select
	mm, cmd := m.selectModel.Update(msg)

	m.selectModel = mm.(components.SelectModel)

	return m, cmd
}

func (m *WelcomeModel) View() string {
	var b strings.Builder

	b.WriteString("What would you like to do?\n\n")

	// Render select
	b.WriteString(m.selectModel.View())

	return b.String()
}
