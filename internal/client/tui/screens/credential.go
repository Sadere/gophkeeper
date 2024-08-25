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

var errMetadataEmpty = "Please enter metadata"

const (
	credMetadata = iota
	credLogin
	credPassword
)

type CredentialModel struct {
	state      *State
	credID     uint64
	inputGroup components.InputGroup
	errorMsg   string
}

func NewCredentialModel(state *State, ID uint64) *CredentialModel {
	inputs := make([]textinput.Model, 3)

	m := CredentialModel{
		state:  state,
		credID: ID,
	}

	var t textinput.Model
	for i := range inputs {
		t = textinput.New()
		t.CharLimit = 100

		switch i {
		case credMetadata:
			t.Placeholder = "Metadata"
			t.Focus()
		case credLogin:
			t.Placeholder = "Login"
		case credPassword:
			t.Placeholder = "Password"
		}

		inputs[i] = t
	}

	// Load secret if in view/edit mode
	if ID > 0 {
		secret, err := m.state.client.LoadSecret(context.Background(), ID)
		if err != nil {
			m.errorMsg = err.Error()
			return &m
		}

		inputs[credMetadata].SetValue(secret.Metadata)
		inputs[credLogin].SetValue(secret.Creds.Login)
		inputs[credPassword].SetValue(secret.Creds.Password)
	}

	m.inputGroup = components.NewInputGroup(inputs)

	return &m
}

func (m CredentialModel) Init() tea.Cmd {
	return m.inputGroup.Init()
}

func (m CredentialModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m CredentialModel) Submit() (tea.Model, tea.Cmd) {
	metadata := m.inputGroup.Inputs[credMetadata].Value()
	login := m.inputGroup.Inputs[credLogin].Value()
	password := m.inputGroup.Inputs[credPassword].Value()

	// Validate inputs
	if len(metadata) == 0 {
		m.errorMsg = errMetadataEmpty
		return m, nil
	}

	if len(login) == 0 {
		m.errorMsg = errLoginEmpty
		return m, nil
	}

	if len(password) == 0 {
		m.errorMsg = errPasswordEmpty
		return m, nil
	}

	// Save credential
	err := m.state.client.SaveCredential(context.Background(), m.credID, metadata, login, password)
	if err != nil {
		m.errorMsg = err.Error()
		return m, nil
	}

	// Return to list
	mainScreen := NewSecretListModel(m.state)
	return NewRootModel(m.state).SwitchScreen(mainScreen)
}

func (m CredentialModel) View() string {
	var b strings.Builder

	b.WriteString("Fill in credential details and metadata\n")

	// View inputs
	b.WriteString(m.inputGroup.View())

	body := style.RenderBox(b.String())

	if len(m.errorMsg) > 0 {
		errorBox := style.ErrorStyle.Render(m.errorMsg)

		body = lipgloss.JoinVertical(lipgloss.Top, body, errorBox)
	}

	return body
}
