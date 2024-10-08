// Provides lipgloss styles and render helpers for TUI app
package style

import "github.com/charmbracelet/lipgloss"

var (
	ErrorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff1744"))

	NoStyle      = lipgloss.NewStyle()
	FocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	BlurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	HelpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	NewSecretStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#388E3C"))
	UpdatedSecretStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFEB3B"))

	BorderColor = lipgloss.Color("56")
)

// Render body in a box with border
func RenderBox(body string) string {
	borderBox := NoStyle.
		Padding(0, 1).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(BorderColor)

	return borderBox.Render(body)
}
