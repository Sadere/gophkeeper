package screens

import (
	"github.com/Sadere/gophkeeper/internal/client/tui/style"
	tea "github.com/charmbracelet/bubbletea"
)

func RenderHelpForModel(model tea.Model) string {
	help := "↑/↓: navigate • ctrl+c: quit"

	switch model.(type) {
	case LoginModel, RegisterModel:
		help += " • enter: submit"
	case *SecretListModel:
		help += " • r: refresh list • a: add new • enter: view/edit"
	}

	return style.HelpStyle.Render(help)
}
