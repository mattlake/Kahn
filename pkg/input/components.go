package input

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"kahn/pkg/colors"
)

// InputComponents holds the text input components for forms
type InputComponents struct {
	NameInput textinput.Model
	DescInput textinput.Model
}

// NewInputComponents creates and initializes input components for task forms
func NewInputComponents() InputComponents {
	name := textinput.New()
	name.Placeholder = "Task name"
	name.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Subtext0))
	name.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Text))
	name.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Mauve))
	name.Focus()
	name.CharLimit = 50
	name.Width = 40

	desc := textinput.New()
	desc.Placeholder = "Task description"
	desc.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Subtext0))
	desc.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Text))
	desc.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Mauve))
	desc.CharLimit = 100
	desc.Width = 40

	return InputComponents{
		NameInput: name,
		DescInput: desc,
	}
}

// Reset resets all input components
func (ic *InputComponents) Reset() {
	ic.NameInput.Reset()
	ic.DescInput.Reset()
}

// FocusName focuses the name input
func (ic *InputComponents) FocusName() {
	ic.NameInput.Focus()
	ic.DescInput.Blur()
}

// FocusDesc focuses the description input
func (ic *InputComponents) FocusDesc() {
	ic.DescInput.Focus()
	ic.NameInput.Blur()
}

// Blur blurs all inputs
func (ic *InputComponents) Blur() {
	ic.NameInput.Blur()
	ic.DescInput.Blur()
}
