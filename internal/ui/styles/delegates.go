package styles

import (
	"kahn/internal/domain"
	"kahn/internal/ui/colors"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

// PERFORMANCE: Pre-allocated style objects to avoid recreating on every render
var (
	// Selection styling
	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colors.Blue)).
			Bold(true)

	// Blocked task styling
	blockedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colors.Red)).
			Bold(false)

	blockedSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colors.Red)).
				Bold(true)

	// Priority color styles (cached) - using values instead of pointers
	priorityStyles = map[domain.Priority]lipgloss.Style{
		domain.Low:    lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Green)),
		domain.Medium: lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Peach)),
		domain.High:   lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Red)),
	}
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

func (t TaskWithTitle) Title() string {

	title := t.Task.Title()
	switch t.Task.Type {
	case domain.RegularTask:
		title = "󰄬 " + title
	case domain.Bug:
		title = "󰃤 " + title
	case domain.Feature:
		title = "󱕣 " + title
	}

	// Check if task is blocked - render in red to indicate it's blocked
	if t.Task.BlockedBy != nil {
		if t.isSelected && t.isActiveList {
			return blockedSelectedStyle.Render(t.priorityText + title)
		}
		return blockedStyle.Render(t.priorityText + title)
	}

	// Original behavior for non-blocked tasks
	if t.isSelected && t.isActiveList {

		return selectedStyle.Render(t.priorityText + title)
	} else {

		priorityStyled := priorityStyles[t.Task.Priority].Render(t.priorityText)
		return priorityStyled + title
	}
}

// GetTaskType returns the task type for interface compliance
func (t TaskWithTitle) GetTaskType() domain.TaskType {
	return t.Task.Type
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
