package components

import (
	"github.com/charmbracelet/bubbles/list"
	"kahn/internal/domain"
)

// BoardRenderer defines the interface for board-related UI rendering
type BoardRenderer interface {
	// RenderProjectFooter renders the bottom project footer with name and help text
	RenderProjectFooter(project *domain.Project, width int, version string) string

	// RenderSearchBar renders the search input bar at the bottom when search is active
	RenderSearchBar(query string, matchCount int, width int) string

	// RenderNoProjectsBoard renders the empty state when no projects exist
	RenderNoProjectsBoard(width, height int) string

	// RenderTaskDeleteConfirm renders the task deletion confirmation dialog
	RenderTaskDeleteConfirm(task *domain.Task, width, height int) string

	// RenderTaskDeleteConfirmWithError renders the task deletion confirmation with error information
	RenderTaskDeleteConfirmWithError(task *domain.Task, errorMessage string, width, height int) string

	// RenderBoard renders the main kanban board with three columns.
	// When searchActive is true, displays search bar instead of project footer.
	RenderBoard(project *domain.Project, taskLists [3]list.Model, activeListIndex domain.Status, width int, version string, searchActive bool, searchQuery string, searchMatchCount int) string
}
