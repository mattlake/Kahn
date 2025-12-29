package styles

import (
	"github.com/charmbracelet/bubbles/textinput"
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

// ConfigureNameInput configures a text input for name fields
func ConfigureNameInput(styles FormFieldStyles) textinput.Model {
	input := textinput.New()
	input.Placeholder = "Name *"
	input.PlaceholderStyle = styles.Placeholder
	input.TextStyle = styles.Text
	input.Cursor.Style = styles.Cursor
	input.CharLimit = 50
	input.Width = 40
	return input
}

// ConfigureDescriptionInput configures a text input for description fields
func ConfigureDescriptionInput(styles FormFieldStyles) textinput.Model {
	input := textinput.New()
	input.Placeholder = "Description (optional)"
	input.PlaceholderStyle = styles.Placeholder
	input.TextStyle = styles.Text
	input.Cursor.Style = styles.Cursor
	input.CharLimit = 200 // Larger for project descriptions
	input.Width = 40
	return input
}

// ConfigureTaskNameInput configures a text input specifically for task names
func ConfigureTaskNameInput(styles FormFieldStyles) textinput.Model {
	input := textinput.New()
	input.Placeholder = "Task name *"
	input.PlaceholderStyle = styles.Placeholder
	input.TextStyle = styles.Text
	input.Cursor.Style = styles.Cursor
	input.CharLimit = 50
	input.Width = 40
	return input
}

// ConfigureTaskDescriptionInput configures a text input for task descriptions
func ConfigureTaskDescriptionInput(styles FormFieldStyles) textinput.Model {
	input := textinput.New()
	input.Placeholder = "Task description (optional)"
	input.PlaceholderStyle = styles.Placeholder
	input.TextStyle = styles.Text
	input.Cursor.Style = styles.Cursor
	input.CharLimit = 200
	input.Width = 40
	return input
}

// ConfigureProjectNameInput configures a text input for project names
func ConfigureProjectNameInput(styles FormFieldStyles) textinput.Model {
	input := textinput.New()
	input.Placeholder = "Project name *"
	input.PlaceholderStyle = styles.Placeholder
	input.TextStyle = styles.Text
	input.Cursor.Style = styles.Cursor
	input.CharLimit = 50
	input.Width = 40
	return input
}

// ConfigureProjectDescriptionInput configures a text input for project descriptions
func ConfigureProjectDescriptionInput(styles FormFieldStyles) textinput.Model {
	input := textinput.New()
	input.Placeholder = "Project description (optional)"
	input.PlaceholderStyle = styles.Placeholder
	input.TextStyle = styles.Text
	input.Cursor.Style = styles.Cursor
	input.CharLimit = 200
	input.Width = 40
	return input
}
