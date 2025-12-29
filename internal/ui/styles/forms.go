package styles

import (
	"github.com/charmbracelet/lipgloss"
	"kahn/internal/ui/colors"
)

// FormFieldStyles provides styling for form input fields
type FormFieldStyles struct {
	Placeholder lipgloss.Style
	Text        lipgloss.Style
	Cursor      lipgloss.Style
	Border      lipgloss.Style
	ErrorBorder lipgloss.Style
	ErrorText   lipgloss.Style
}

// GetFormFieldStyles returns consistent form field styling
func GetFormFieldStyles() FormFieldStyles {
	return FormFieldStyles{
		Placeholder: lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Subtext0)),
		Text:        lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Text)),
		Cursor:      lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Mauve)),
		Border: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(colors.Overlay1)).
			Padding(0, 1),
		ErrorBorder: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(colors.Red)).
			Padding(0, 1),
		ErrorText: lipgloss.NewStyle().
			Foreground(lipgloss.Color(colors.Red)).
			Faint(true),
	}
}
