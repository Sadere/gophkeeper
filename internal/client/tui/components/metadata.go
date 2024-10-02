package components

import (
	"github.com/Sadere/gophkeeper/internal/client/tui/style"
	"github.com/charmbracelet/bubbles/textinput"
)

// Create input field for metadata
func NewMetaDataInput() textinput.Model {
	input := textinput.New()

	input.Cursor.Style = style.FocusedStyle
	input.PromptStyle = style.FocusedStyle
	input.TextStyle = style.FocusedStyle
	input.Focus()
	input.Placeholder = "Metadata"

	return input
}
