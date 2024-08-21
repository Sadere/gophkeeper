package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type RootModel struct {
	state         *State
	windowHeight  int
	windowWidth   int
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
		m.windowHeight = msg.Height
		m.windowWidth = msg.Width

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

	middleBox := lipgloss.NewStyle().Width(m.windowWidth).Height(m.windowHeight - 1)

	body := middleBox.AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render(content)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		body,
	)
}

func (m RootModel) SwitchScreen(model tea.Model) (tea.Model, tea.Cmd) {
	m.currentWindow = model
	return m.currentWindow, m.currentWindow.Init()
}
