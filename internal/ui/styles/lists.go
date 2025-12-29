package styles

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"kahn/internal/domain"
	"kahn/internal/ui/colors"
)

// ListTitleStyles provides styles for kanban list titles
type ListTitleStyles struct {
	NotStarted lipgloss.Style
	InProgress lipgloss.Style
	Done       lipgloss.Style
}

// GetListTitleStyles returns styled titles for all three kanban lists
func GetListTitleStyles() ListTitleStyles {
	return ListTitleStyles{
		NotStarted: lipgloss.NewStyle().
			Foreground(lipgloss.Color(colors.Blue)).
			Bold(true).
			Align(lipgloss.Center),
		InProgress: lipgloss.NewStyle().
			Foreground(lipgloss.Color(colors.Yellow)).
			Bold(true).
			Align(lipgloss.Center),
		Done: lipgloss.NewStyle().
			Foreground(lipgloss.Color(colors.Green)).
			Bold(true).
			Align(lipgloss.Center),
	}
}

// GetListTitleStyle returns style for specific status
func GetListTitleStyle(status domain.Status) lipgloss.Style {
	styles := GetListTitleStyles()
	switch status {
	case domain.NotStarted:
		return styles.NotStarted
	case domain.InProgress:
		return styles.InProgress
	case domain.Done:
		return styles.Done
	default:
		return styles.NotStarted
	}
}

// ApplyListTitleStyles applies proper title styles to task lists
func ApplyListTitleStyles(taskLists []list.Model) {
	ApplyTitleStyles(taskLists, domain.NotStarted)
}

// ApplyFocusedTitleStyles applies focused title styles to task lists
func ApplyFocusedTitleStyles(taskLists []list.Model, activeListIndex domain.Status) {
	if len(taskLists) < 3 {
		return
	}

	for i := range taskLists {
		status := domain.Status(i)
		if status == activeListIndex {
			// Focused list gets green title
			taskLists[i].Styles.Title = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colors.Green)).
				Bold(true).
				Align(lipgloss.Center)
		} else {
			// Unfocused lists get white title
			taskLists[i].Styles.Title = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colors.Text)).
				Bold(true).
				Align(lipgloss.Center)
		}
	}
}

// ApplyTitleStyles applies basic title styles (for backward compatibility)
func ApplyTitleStyles(taskLists []list.Model, activeListIndex domain.Status) {
	if len(taskLists) < 3 {
		return
	}

	for i := range taskLists {
		status := domain.Status(i)
		if status == activeListIndex {
			// Focused list gets green title
			taskLists[i].Styles.Title = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colors.Green)).
				Bold(true).
				Align(lipgloss.Center)
		} else {
			// Unfocused lists get white title
			taskLists[i].Styles.Title = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colors.Text)).
				Bold(true).
				Align(lipgloss.Center)
		}
	}
}
