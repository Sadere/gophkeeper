package screens

import (
	"context"
	"strings"

	"github.com/Sadere/gophkeeper/internal/client/tui/components"
	"github.com/Sadere/gophkeeper/internal/client/tui/style"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	errPasswordMismatch     = "Passwords don't match"
	errLoginEmpty           = "Please enter login"
	errPasswordEmpty        = "Please enter password"
	errPasswordConfirmEmpty = "Please enter password second time"
)

const (
	registerLogin = iota
	registerPassword
	registerConfirmPassword
)

type RegisterModel struct {
	state      *State
	inputGroup components.InputGroup
	errorMsg   string
}

func NewRegisterModel(state *State) *RegisterModel {
	inputs := make([]textinput.Model, 3)

	m := RegisterModel{
		state: state,
	}

	var t textinput.Model
	for i := range inputs {
		t = textinput.New()
		t.CharLimit = 32

		switch i {
		case registerLogin:
			t.Placeholder = "Login"
			t.Focus()
		case registerPassword:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		case registerConfirmPassword:
			t.Placeholder = "Confirm password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		inputs[i] = t
	}

	m.inputGroup = components.NewInputGroup(inputs)

	return &m
}

func (m RegisterModel) Init() tea.Cmd {
	return m.inputGroup.Init()
}

func (m RegisterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m RegisterModel) Submit() (tea.Model, tea.Cmd) {
	login := m.inputGroup.Inputs[registerLogin].Value()
	password := m.inputGroup.Inputs[registerPassword].Value()
	confirmPassword := m.inputGroup.Inputs[registerConfirmPassword].Value()

	// Validate inputs
	if len(login) == 0 {
		m.errorMsg = errLoginEmpty
		return m, nil
	}

	if len(password) == 0 {
		m.errorMsg = errPasswordEmpty
		return m, nil
	}

	if len(confirmPassword) == 0 {
		m.errorMsg = errPasswordConfirmEmpty
		return m, nil
	}

	if password != confirmPassword {
		m.errorMsg = errPasswordMismatch
		return m, nil
	}

	// Register
	accessToken, err := m.state.client.Register(context.Background(), login, password)
	if err != nil {
		m.errorMsg = err.Error()
		return m, nil
	}

	// Proceed to main screen
	m.state.accessToken = accessToken

	mainScreen := NewSecretListModel(m.state)
	return NewRootModel(m.state).SwitchScreen(mainScreen)
}

func (m RegisterModel) View() string {
	var b strings.Builder

	b.WriteString("Enter new user credentials\n")

	// View inputs
	b.WriteString(m.inputGroup.View())

	body := style.RenderBox(b.String())

	if len(m.errorMsg) > 0 {
		errorBox := style.ErrorStyle.Render(m.errorMsg)

		body = lipgloss.JoinVertical(lipgloss.Top, body, errorBox)
	}

	return body
}
