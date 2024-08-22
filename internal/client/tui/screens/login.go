package screens

import (
	"context"
	"strings"

	"github.com/Sadere/gophkeeper/internal/client/tui/components"
	"github.com/Sadere/gophkeeper/internal/client/tui/style"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
)

type LoginModel struct {
	state      *State
	inputGroup components.InputGroup
	errorMsg   string
}

func NewLoginModel(state *State) *LoginModel {
	inputs := make([]textinput.Model, 2)

	m := LoginModel{
		state: state,
	}

	var t textinput.Model
	for i := range inputs {
		t = textinput.New()
		t.Cursor.Style = style.FocusedStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Login"
			t.Focus()
			t.PromptStyle = style.FocusedStyle
			t.TextStyle = style.FocusedStyle
		case 1:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		inputs[i] = t
	}

	m.inputGroup = components.NewInputGroup(inputs)

	return &m
}

func (m LoginModel) Init() tea.Cmd {
	return nil
}

func (m LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.inputGroup.FocusIndex == m.inputGroup.InputNum {
				return m.Submit()
			}
		}
	}

	// Handle input group
	mm, cmd := m.inputGroup.Update(msg)

	m.inputGroup = mm.(components.InputGroup)

	return m, cmd
}

func (m LoginModel) Submit() (tea.Model, tea.Cmd) {
	login := m.inputGroup.Inputs[0].Value()
	password := m.inputGroup.Inputs[1].Value()

	// Validate inputs
	if len(login) == 0 {
		m.errorMsg = errLoginEmpty
		return m, nil
	}

	if len(password) == 0 {
		m.errorMsg = errPasswordEmpty
		return m, nil
	}

	// Login
	accessToken, err := m.state.client.Login(context.Background(), login, password)
	if err != nil {
		m.errorMsg = err.Error()
		return m, nil
	}

	// Proceed to main screen
	m.state.accessToken = accessToken

	mainScreen := NewSecretListModel(m.state)
	return NewRootModel(m.state).SwitchScreen(mainScreen)
}

func (m LoginModel) View() string {
	var b strings.Builder

	b.WriteString("Please provide correct credentials\n")

	// View inputs
	b.WriteString(m.inputGroup.View())

	body := style.RenderBox(b.String())

	if len(m.errorMsg) > 0 {
		errorBox := style.ErrorStyle.Render(m.errorMsg)

		body = lipgloss.JoinVertical(lipgloss.Top, body, errorBox)
	}

	return body
}
