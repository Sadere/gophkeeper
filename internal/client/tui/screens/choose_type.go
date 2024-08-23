package screens

import (
	"strings"

	"github.com/Sadere/gophkeeper/internal/client/tui/components"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	selectBack = iota
	selectCredential
	selectText
	selectBlob
	selectCard
)

type ChooseTypeModel struct {
	state       *State
	selectModel components.SelectModel
}

func NewChooseTypeModel(state *State) *ChooseTypeModel {
	choices := []string{
		selectBack:       "Go back",
		selectCredential: "Add new credentials",
		selectText:       "Add text",
		selectBlob:       "Upload file",
		selectCard:       "Add card info",
	}

	return &ChooseTypeModel{
		state:       state,
		selectModel: components.NewSelectModel(choices),
	}
}

func (m ChooseTypeModel) Init() tea.Cmd {
	return tea.SetWindowTitle("GophKeeper client")
}

func (m *ChooseTypeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			switch m.selectModel.Selected() {
			case selectBack:
				// Back to secret list
				mainScreen := NewSecretListModel(m.state)
				return NewRootModel(m.state).SwitchScreen(mainScreen)
			}
		}
	}

	// Update select
	mm, cmd := m.selectModel.Update(msg)

	m.selectModel = mm.(components.SelectModel)

	return m, cmd
}

func (m *ChooseTypeModel) View() string {
	var b strings.Builder

	b.WriteString("What would you like to do?\n\n")

	// Render select
	b.WriteString(m.selectModel.View())

	return b.String()
}
