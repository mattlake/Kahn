package styles

import (
	"kahn/internal/domain"
	"kahn/internal/ui/colors"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

// TaskWithTitle wraps a domain.Task with priority-formatted title
// This allows us to keep Task.Title() pure while modifying display
type TaskWithTitle struct {
	domain.Task
	priorityText string
	isSelected   bool
	isActiveList bool
}

// Title returns the priority-formatted title for display
// Smart rendering based on selection state and list activity
func (t TaskWithTitle) Title() string {
	if t.isSelected && t.isActiveList {
		// Blue selection styling for both priority and title
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(colors.Blue)).
			Bold(true).
			Render(t.priorityText + t.Task.Title())
	} else {
		// Individual priority colors for unselected, default text for title
		priorityStyled := lipgloss.NewStyle().
			Foreground(lipgloss.Color(getPriorityColor(t.Task.Priority))).
			Render(t.priorityText)
		return priorityStyled + t.Task.Title()
	}
}

// NewTaskWithTitle creates a wrapper for display purposes
func NewTaskWithTitle(task domain.Task) TaskWithTitle {
	return TaskWithTitle{
		Task:         task,
		priorityText: GetPriorityIndicatorUncolored(task.Priority),
	}
}

// UpdateTaskSelection updates selection state for all items in a list
func UpdateTaskSelection(items []list.Item, selectedIndex int, isActiveList bool) []list.Item {
	updatedItems := make([]list.Item, len(items))
	for i, item := range items {
		if taskItem, ok := item.(TaskWithTitle); ok {
			taskItem.isSelected = (i == selectedIndex)
			taskItem.isActiveList = isActiveList
			updatedItems[i] = taskItem
		} else {
			updatedItems[i] = item
		}
	}
	return updatedItems
}

// NewActiveListDelegate creates a list delegate for active list with priority indicators
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
		Foreground(lipgloss.Color(colors.Blue)).
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
