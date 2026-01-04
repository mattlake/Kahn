package components

import (
	"github.com/charmbracelet/bubbles/list"
	"kahn/internal/domain"
)

// BoardRenderer defines the interface for board-related UI rendering
type BoardRenderer interface {
	// RenderProjectFooter renders the bottom project footer with name and help text
	RenderProjectFooter(project *domain.Project, width int, version string) string

	// RenderNoProjectsBoard renders the empty state when no projects exist
	RenderNoProjectsBoard(width, height int) string

	// RenderTaskDeleteConfirm renders the task deletion confirmation dialog
	RenderTaskDeleteConfirm(task *domain.Task, width, height int) string

	// RenderBoard renders the main kanban board with three columns
	RenderBoard(project *domain.Project, taskLists [3]list.Model, activeListIndex domain.Status, width int, version string) string
}
