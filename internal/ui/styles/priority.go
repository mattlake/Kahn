package styles

import (
	"kahn/internal/domain"
	"kahn/internal/ui/colors"

	"github.com/charmbracelet/lipgloss"
)

// priorityIndicatorWidth defines the target width for all priority indicators (in visible characters)
const priorityIndicatorWidth = 2 // 2 visible chars + 1 trailing space = 3 total

// GetPriorityIndicator returns formatted priority indicator with color
// High: »» (red), Medium: » (peach), Low:  (padded spaces)
func GetPriorityIndicator(priority domain.Priority) string {
	var indicator, color string

	switch priority {
	case domain.High:
		indicator = "»»"
		color = colors.Red
	case domain.Medium:
		indicator = "»"
		color = colors.Peach
	case domain.Low:
		indicator = "  "     // 2 spaces for low priority
		color = colors.Green // Available for future use
	default:
		indicator = "»"
		color = colors.Peach
	}

	// Apply color to indicator and add trailing space for separation
	coloredIndicator := lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Render(indicator)

	// Add trailing space for visual separation from task name
	return coloredIndicator + " "
}

// getPriorityColor returns the color for a given priority
func getPriorityColor(priority domain.Priority) string {
	switch priority {
	case domain.High:
		return colors.Red
	case domain.Medium:
		return colors.Peach
	case domain.Low:
		return colors.Green // Available for future use
	default:
		return colors.Peach
	}
}

// GetPriorityIndicatorWidth returns the fixed width for alignment calculations
func GetPriorityIndicatorWidth() int {
	return priorityIndicatorWidth
}

// GetPriorityIndicatorUncolored returns raw priority indicator without styling
// High: »» , Medium: »  , Low:
// Matches GetPriorityIndicator structure but without lipgloss styling
func GetPriorityIndicatorUncolored(priority domain.Priority) string {
	switch priority {
	case domain.High:
		return "»» " // 2 chevrons + space
	case domain.Medium:
		return "»  " // 1 chevron + 2 spaces + space
	case domain.Low:
		return "   " // 2 spaces + space
	default:
		return "   " // 1 chevron + 2 spaces + space
	}
}

// FormatTaskWithPriority applies priority indicator to pure task title
// This keeps Task.Title() pure while adding UI-specific formatting
func FormatTaskWithPriority(task domain.Task) string {
	indicator := GetPriorityIndicator(task.Priority)
	return indicator + task.Title()
}
