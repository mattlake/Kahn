package styles

import (
	"github.com/charmbracelet/lipgloss"
	"kahn/internal/ui/colors"
)

// Lipgloss styles with Catppuccin colors
var DefaultStyle = lipgloss.NewStyle().
	Margin(1, 2).
	Border(lipgloss.HiddenBorder()).
	Padding(1, 2)

var FocusedStyle = lipgloss.NewStyle().
	Margin(1, 2).
	Border(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color(colors.Mauve)).
	Padding(1, 2)
