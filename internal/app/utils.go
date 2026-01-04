package app

import (
	"github.com/charmbracelet/bubbles/list"
	"kahn/internal/domain"
	"kahn/internal/ui/styles"
)

func convertTasksToListItems(tasks []domain.Task) []list.Item {
	items := make([]list.Item, len(tasks))
	for i, task := range tasks {

		items[i] = styles.NewTaskWithTitle(task)
	}
	return items
}
