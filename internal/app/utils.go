package app

import (
	"github.com/charmbracelet/bubbles/list"
	"kahn/internal/domain"
	"kahn/internal/ui/styles"
)

// convertTasksToListItems converts domain tasks to list items with styling
func convertTasksToListItems(tasks []domain.Task) []list.Item {
	items := make([]list.Item, len(tasks))
	for i, task := range tasks {
		// Wrap task with priority-formatted title
		items[i] = styles.NewTaskWithTitle(task)
	}
	return items
}
