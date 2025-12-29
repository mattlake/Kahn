package main

import (
	"github.com/charmbracelet/bubbles/list"
	"kahn/internal/domain"
)

// convertTasksToListItems converts domain tasks to list items for Bubble Tea lists
func convertTasksToListItems(tasks []domain.Task) []list.Item {
	items := make([]list.Item, len(tasks))
	for i, task := range tasks {
		items[i] = task
	}
	return items
}
