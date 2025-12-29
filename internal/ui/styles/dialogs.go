package styles

import (
	"github.com/charmbracelet/lipgloss"
	"kahn/internal/ui/colors"
)

// DialogStyles provides styles for modal dialogs and forms
type DialogStyles struct {
	Title       lipgloss.Style
	Message     lipgloss.Style
	Instruction lipgloss.Style
	Form        lipgloss.Style
	Error       lipgloss.Style
}

// GetDialogStyles returns consistent dialog styling
func GetDialogStyles() DialogStyles {
	return DialogStyles{
		Title: lipgloss.NewStyle().
			Foreground(lipgloss.Color(colors.Mauve)).
			Bold(true).
			Align(lipgloss.Center),
		Message: lipgloss.NewStyle().
			Foreground(lipgloss.Color(colors.Text)).
			Align(lipgloss.Center),
		Instruction: lipgloss.NewStyle().
			Foreground(lipgloss.Color(colors.Subtext1)).
			Align(lipgloss.Center),
		Form: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(colors.Green)).
			Padding(2, 3),
		Error: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(colors.Red)).
			Padding(2, 3),
	}
}

// ProjectItemStyle provides styling for project items in switcher
type ProjectItemStyle struct {
	Normal lipgloss.Style
	Active lipgloss.Style
}

// GetProjectItemStyle returns styling for project list items
func GetProjectItemStyle(color string) ProjectItemStyle {
	return ProjectItemStyle{
		Normal: lipgloss.NewStyle().
			Foreground(lipgloss.Color(color)),
		Active: lipgloss.NewStyle().
			Foreground(lipgloss.Color(color)).
			Background(lipgloss.Color(colors.Surface1)),
	}
}

// DeleteConfirmStyles provides styling for delete confirmation dialogs
type DeleteConfirmStyles struct {
	Title   lipgloss.Style
	Warning lipgloss.Style
	Message lipgloss.Style
	Form    lipgloss.Style
}

// GetDeleteConfirmStyles returns styling for delete confirmations
func GetDeleteConfirmStyles() DeleteConfirmStyles {
	return DeleteConfirmStyles{
		Title: lipgloss.NewStyle().
			Foreground(lipgloss.Color(colors.Red)).
			Bold(true).
			Align(lipgloss.Center),
		Warning: lipgloss.NewStyle().
			Foreground(lipgloss.Color(colors.Text)).
			Bold(true).
			Align(lipgloss.Center),
		Message: lipgloss.NewStyle().
			Foreground(lipgloss.Color(colors.Subtext1)).
			Align(lipgloss.Center),
		Form: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(colors.Red)).
			Padding(2, 3),
	}
}
