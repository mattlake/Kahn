package styles

import (
	"kahn/internal/ui/colors"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

// NewActiveListDelegate creates a list delegate for the active list with prominent selection styling
func NewActiveListDelegate() list.DefaultDelegate {
	delegate := list.NewDefaultDelegate()

	// Disable description rendering completely
	delegate.ShowDescription = false
	delegate.SetSpacing(0)

	// Normal title styling
	delegate.Styles.NormalTitle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Text))

	// Prominent selection styling for active list
	delegate.Styles.SelectedTitle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Surface0)).
		Background(lipgloss.Color(colors.Green)).
		Bold(true)

	return delegate
}

// NewInactiveListDelegate creates a list delegate for inactive lists with minimal selection styling
func NewInactiveListDelegate() list.DefaultDelegate {
	delegate := list.NewDefaultDelegate()

	// Disable description rendering completely
	delegate.ShowDescription = false
	delegate.SetSpacing(0)

	// Normal title styling
	delegate.Styles.NormalTitle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Text))

	// Minimal selection styling - same as normal to hide selection in inactive lists
	delegate.Styles.SelectedTitle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Text))

	return delegate
}

// NewTitleOnlyDelegate creates a list delegate that only displays task titles
// without descriptions, providing a cleaner visual foundation for future UI improvements
// Kept for backward compatibility
func NewTitleOnlyDelegate() list.DefaultDelegate {
	return NewActiveListDelegate()
}
