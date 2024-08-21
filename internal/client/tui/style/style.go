package style

import "github.com/charmbracelet/lipgloss"

var (
	ErrorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff1744"))

	NoStyle      = lipgloss.NewStyle()
	FocusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	BlurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	BorderColor = lipgloss.Color("56")
)

func RenderBox(body string) string {
	borderBox := NoStyle.
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(BorderColor)

	return borderBox.Render(body)
}
