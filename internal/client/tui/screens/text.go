package screens

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Sadere/gophkeeper/internal/client/tui/style"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var errContentEmpty = "Please enter something in text area"

const (
	textInputNum = 3

	textMetadata = iota
	textContent
)

type TextModel struct {
	state         *State
	textID        uint64
	cursor        int
	metadataInput textinput.Model
	content       textarea.Model
	errorMsg      string
}

func NewTextModel(state *State, ID uint64) *TextModel {
	m := TextModel{
		state:         state,
		textID:        ID,
		metadataInput: textinput.New(),
		content:       textarea.New(),
		cursor:        1,
	}

	m.metadataInput.Cursor.Style = style.FocusedStyle
	m.metadataInput.PromptStyle = style.FocusedStyle
	m.metadataInput.TextStyle = style.FocusedStyle
	m.metadataInput.Focus()
	m.metadataInput.Placeholder = "Metadata"

	m.content.Placeholder = "Enter any text"

	// Load secret if in view/edit mode
	if ID > 0 {
		secret, err := m.state.client.LoadSecret(context.Background(), ID)
		if err != nil {
			m.errorMsg = err.Error()
			return &m
		}

		m.metadataInput.SetValue(secret.Metadata)
		m.content.SetValue(secret.Text.Content)
	}

	return &m
}

func (m TextModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m TextModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			mainScreen := NewSecretListModel(m.state)
			return NewRootModel(m.state).SwitchScreen(mainScreen)
		case "enter":
			if m.cursor == textInputNum {
				return m.Submit()
			}
		case "tab", "shift+tab":
			s := msg.String()

			// Cycle indexes
			if s == "shift+tab" {
				m.cursor--
			} else {
				m.cursor++
			}

			if m.cursor > textInputNum {
				m.cursor = 1
			} else if m.cursor < 0 {
				m.cursor = textInputNum
			}

			if m.cursor == textMetadata {
				cmds = append(cmds, m.metadataInput.Focus())
				m.metadataInput.PromptStyle = style.FocusedStyle
				m.metadataInput.TextStyle = style.FocusedStyle
			} else {
				m.metadataInput.Blur()
				m.metadataInput.PromptStyle = style.BlurredStyle
				m.metadataInput.TextStyle = style.BlurredStyle
			}

			if m.cursor == textContent {
				cmds = append(cmds, m.content.Focus())
			} else {
				m.content.Blur()
			}
			log.Printf("%d\n", m.cursor)

			return m, tea.Batch(cmds...)
		}
	}

	// Inputs
	m.metadataInput, cmd = m.metadataInput.Update(msg)
	cmds = append(cmds, cmd)

	m.content, cmd = m.content.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m TextModel) Submit() (tea.Model, tea.Cmd) {
	metadata := m.metadataInput.Value()
	content := m.content.Value()

	// Validate inputs
	if len(metadata) == 0 {
		m.errorMsg = errMetadataEmpty
		return m, nil
	}

	if len(content) == 0 {
		m.errorMsg = errContentEmpty
		return m, nil
	}

	// Save text
	err := m.state.client.SaveText(context.Background(), m.textID, metadata, content)
	if err != nil {
		m.errorMsg = err.Error()
		return m, nil
	}

	// Return to list
	mainScreen := NewSecretListModel(m.state)
	return NewRootModel(m.state).SwitchScreen(mainScreen)
}

func (m TextModel) View() string {
	var b strings.Builder

	b.WriteString("Enter some text\n")

	// View inputs
	b.WriteString(m.metadataInput.View())
	b.WriteRune('\n')
	b.WriteString(m.content.View())

	// View button
	button := style.BlurredStyle.Render("[ Submit ]")
	if m.cursor == textInputNum {
		button = style.FocusedStyle.Render("[ Submit ]")
	}
	fmt.Fprintf(&b, "\n\n%s\n", button)

	body := style.RenderBox(b.String())

	if len(m.errorMsg) > 0 {
		errorBox := style.ErrorStyle.Render(m.errorMsg)

		body = lipgloss.JoinVertical(lipgloss.Top, body, errorBox)
	}

	return body
}
