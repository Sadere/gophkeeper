package screens

import (
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
)

type RootModel struct {
	state         *State
	currentWindow tea.Model
}

func NewRootModel(state *State) *RootModel {
	m := &RootModel{
		state: state,
	}

	m.currentWindow = NewWelcomeModel(state)

	return m
}

func (m RootModel) Init() tea.Cmd {
	return tea.SetWindowTitle("GophKeeper client")
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		curCmd tea.Cmd
		cmds   []tea.Cmd
	)

	m.currentWindow, curCmd = m.currentWindow.Update(msg)

	cmds = append(cmds, curCmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.state.SetSize(msg.Width, msg.Height)

		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			cmds = append(cmds, tea.Quit)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m RootModel) View() string {
	content := m.currentWindow.View()

	middleBox := lipgloss.NewStyle().
		Width(m.state.Width()).
		Height(m.state.Height() - 1)

	body := middleBox.AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render(content)

	help := RenderHelpForModel(m.currentWindow)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		body,
		help,
	)
}

func (m RootModel) SwitchScreen(model tea.Model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.currentWindow = model

	// Update child model window size
	m.currentWindow, cmd = m.currentWindow.Update(tea.WindowSizeMsg{
		Width:  m.state.Width(),
		Height: m.state.Height(),
	})

	return m.currentWindow, tea.Batch(
		// Clear screen fixes issue when changing screens results in leftover view from previous model
		tea.ClearScreen,
		cmd,
		m.currentWindow.Init(),
	)
}
