package screens

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Sadere/gophkeeper/internal/client/tui/components"
	"github.com/Sadere/gophkeeper/internal/client/tui/style"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type status int

const (
	// Statuses
	fileStart status = iota
	filePicking
	fileUpload
	fileComplete
	fileError
)

type fileCompleteMsg struct{}

type errMsg struct {
	msg string
}

type FileModel struct {
	state         *State
	secretID      uint64
	metadataInput textinput.Model
	filepicker    filepicker.Model
	status        status
	selectedFile  string
	errorMsg      string
}

func NewFileModel(state *State, ID uint64) *FileModel {
	var err error

	m := FileModel{
		state:         state,
		secretID:      ID,
		metadataInput: components.NewMetaDataInput(),
		status:        fileStart,
	}

	fp := filepicker.New()
	fp.AutoHeight = false

	fp.CurrentDirectory, err = os.UserHomeDir()
	if err != nil {
		m.errorMsg = fmt.Sprintf("failed to get user home dir: %s", err.Error())
	}

	m.filepicker = fp

	// Load secret if in view/edit mode
	if ID > 0 {
		secret, err := m.state.client.LoadSecret(context.Background(), ID)
		if err != nil {
			m.errorMsg = err.Error()
			return &m
		}

		m.metadataInput.SetValue(secret.Metadata)
	}

	return &m
}

func (m FileModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m FileModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			mainScreen := NewSecretListModel(m.state)
			return NewRootModel(m.state).SwitchScreen(mainScreen)
		case "p":
			// Unfocus metadata
			m.metadataInput.Blur()
			m.metadataInput.PromptStyle = style.BlurredStyle
			m.metadataInput.TextStyle = style.BlurredStyle

			m.status = filePicking
			return m, m.filepicker.Init()
		case "b":
			// Clear errors
			m.errorMsg = ""

			// Focus metadata input
			m.metadataInput.PromptStyle = style.FocusedStyle
			m.metadataInput.TextStyle = style.FocusedStyle

			m.status = fileStart
			return m, m.metadataInput.Focus()
		}
	case fileCompleteMsg:
		m.status = fileComplete
	case errMsg:
		m.errorMsg = msg.msg
		m.status = fileError
	}

	// Handle metadata input
	if m.status == fileStart {
		m.metadataInput, cmd = m.metadataInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	// Set filepicker size
	m.filepicker.Height = m.state.windowHeight - 10

	// Update file picker if it's in focus
	if m.status == filePicking {
		m.filepicker, cmd = m.filepicker.Update(msg)
		cmds = append(cmds, cmd)

		if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
			m.selectedFile = path

			return m, m.uploadStart()
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *FileModel) uploadStart() tea.Cmd {
	m.status = fileUpload

	return func() tea.Msg {
		// validate inputs
		metadata := m.metadataInput.Value()
		if len(metadata) == 0 {
			return errMsg{msg: errMetadataEmpty}
		}

		// upload file
		err := m.state.client.UploadFile(context.Background(), metadata, m.selectedFile)
		if err != nil {
			return errMsg{msg: err.Error()}
		}

		return fileCompleteMsg{}
	}
}

func (m FileModel) View() string {
	var b strings.Builder

	b.WriteString("Pick a file and enter metadata\n")

	b.WriteString(m.metadataInput.View())
	b.WriteString("\n\n")

	switch m.status {
	case fileStart:
		b.WriteString("Press p to pick file")
	case filePicking:
		b.WriteString("Press b to edit metadata\n")
		b.WriteString(m.filepicker.View())
	case fileUpload:
		b.WriteString("File upload in progress, please wait...")
	case fileComplete:
		b.WriteString("File transfer is done, press esc to go back")
	case fileError:
		b.WriteString("Error occured during upload, press b to start over")
	}

	body := style.RenderBox(b.String())

	if len(m.errorMsg) > 0 {
		errorBox := style.ErrorStyle.Render(m.errorMsg)

		body = lipgloss.JoinVertical(lipgloss.Top, body, errorBox)
	}

	return body
}
