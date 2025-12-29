package styles

import (
	"github.com/charmbracelet/lipgloss"
	"kahn/internal/ui/colors"
)

// Lipgloss styles with Catppuccin colors - minimal viewport edge margins
var DefaultStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color(colors.Text)).
	Padding(1, 2)

var FocusedStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color(colors.Green)).
	Padding(1, 2)
