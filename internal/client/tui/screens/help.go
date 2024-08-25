package screens

import (
	"fmt"

	"github.com/Sadere/gophkeeper/internal/client/tui/style"
	tea "github.com/charmbracelet/bubbletea"
)

func RenderHelpForModel(model tea.Model) string {
	help := "↑/↓: navigate • ctrl+c: quit"

	backHelp := " • esc: back"

	switch model.(type) {
	case LoginModel, RegisterModel:
		help += " • enter: submit"
	case *SecretListModel:
		help += " • r: refresh list • a: add new • enter: view/edit"
	case CredentialModel, CardModel:
		help += backHelp
	case TextModel:
		help += fmt.Sprintf("%s • tab: next input • shit+tab: previous input", backHelp)
	}

	return style.HelpStyle.Render(help)
}
