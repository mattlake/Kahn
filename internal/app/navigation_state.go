package app

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"kahn/internal/domain"
	"kahn/internal/services"
	"kahn/internal/ui/styles"
)

// NavigationState manages navigation and project switching
type NavigationState struct {
	showProjectSwitch bool
	activeListIndex   domain.Status
	Tasks             []list.Model
	// PERFORMANCE: Dirty flag tracking for incremental updates
	dirtyFlags map[domain.Status]bool
}

// NewNavigationState creates a new navigation state
func NewNavigationState(tasks []list.Model) *NavigationState {
	return &NavigationState{
		Tasks:           tasks,
		activeListIndex: domain.NotStarted,
	}
}

// ShowProjectSwitch shows the project switcher
func (ns *NavigationState) ShowProjectSwitch() {
	ns.showProjectSwitch = true
}

// HideProjectSwitch hides the project switcher
func (ns *NavigationState) HideProjectSwitch() {
	ns.showProjectSwitch = false
}

// IsShowingProjectSwitch returns whether project switcher is shown
func (ns *NavigationState) IsShowingProjectSwitch() bool {
	return ns.showProjectSwitch
}

// GetActiveListIndex returns the currently active list index
func (ns *NavigationState) GetActiveListIndex() domain.Status {
	return ns.activeListIndex
}

// SetActiveListIndex sets the active list index
func (ns *NavigationState) SetActiveListIndex(index domain.Status) {
	ns.activeListIndex = index
}

// NextList moves focus to the next list
func (ns *NavigationState) NextList() {
	// Switch old active list to inactive delegate
	oldActiveIndex := ns.activeListIndex
	ns.Tasks[oldActiveIndex].SetDelegate(styles.NewInactiveListDelegate())

	// Update selection state for old list (no longer active)
	oldItems := ns.Tasks[oldActiveIndex].Items()
	ns.Tasks[oldActiveIndex].SetItems(styles.UpdateTaskSelection(oldItems, ns.Tasks[oldActiveIndex].Index(), false))

	// Move to next list
	if ns.activeListIndex == domain.Done {
		ns.activeListIndex = domain.NotStarted
	} else {
		ns.activeListIndex++
	}

	// Switch new active list to active delegate
	ns.Tasks[ns.activeListIndex].SetDelegate(styles.NewActiveListDelegate())

	// Update selection state for new list (now active)
	newItems := ns.Tasks[ns.activeListIndex].Items()
	ns.Tasks[ns.activeListIndex].SetItems(styles.UpdateTaskSelection(newItems, ns.Tasks[ns.activeListIndex].Index(), true))

	// Update title styles to reflect new focus
	styles.ApplyFocusedTitleStyles(ns.Tasks[:], ns.activeListIndex)
}

// PrevList moves focus to the previous list
func (ns *NavigationState) PrevList() {
	// Switch old active list to inactive delegate
	oldActiveIndex := ns.activeListIndex
	ns.Tasks[oldActiveIndex].SetDelegate(styles.NewInactiveListDelegate())

	// Update selection state for old list (no longer active)
	oldItems := ns.Tasks[oldActiveIndex].Items()
	ns.Tasks[oldActiveIndex].SetItems(styles.UpdateTaskSelection(oldItems, ns.Tasks[oldActiveIndex].Index(), false))

	// Move to previous list
	if ns.activeListIndex == domain.NotStarted {
		ns.activeListIndex = domain.Done
	} else {
		ns.activeListIndex--
	}

	// Switch new active list to active delegate
	ns.Tasks[ns.activeListIndex].SetDelegate(styles.NewActiveListDelegate())

	// Update selection state for new list (now active)
	newItems := ns.Tasks[ns.activeListIndex].Items()
	ns.Tasks[ns.activeListIndex].SetItems(styles.UpdateTaskSelection(newItems, ns.Tasks[ns.activeListIndex].Index(), true))

	// Update title styles to reflect new focus
	styles.ApplyFocusedTitleStyles(ns.Tasks[:], ns.activeListIndex)
}

// UpdateTaskLists updates the task lists with new data
func (ns *NavigationState) UpdateTaskLists(project *domain.Project, taskService *services.TaskService) {
	if project == nil {
		return
	}

	// Save current selection states before refreshing
	notStartedIndex := ns.Tasks[domain.NotStarted].Index()
	inProgressIndex := ns.Tasks[domain.InProgress].Index()
	doneIndex := ns.Tasks[domain.Done].Index()

	// PERFORMANCE: Single database query to get all tasks for project
	// This replaces 3 separate GetTasksByStatus calls with 1 GetTasksByProject call (66% reduction in DB calls)
	allTasks, err := taskService.GetTasksByProject(project.ID)
	if err != nil {
		// Handle error - set empty lists
		ns.Tasks[domain.NotStarted].SetItems([]list.Item{})
		ns.Tasks[domain.InProgress].SetItems([]list.Item{})
		ns.Tasks[domain.Done].SetItems([]list.Item{})
		return
	}

	// Update project tasks with all tasks from database
	project.Tasks = allTasks

	// Convert tasks to list items and update each status list
	// This preserves the existing filtering and UI behavior while using single DB query
	ns.Tasks[domain.NotStarted].SetItems(convertTasksToListItems(project.GetTasksByStatus(domain.NotStarted)))
	ns.Tasks[domain.InProgress].SetItems(convertTasksToListItems(project.GetTasksByStatus(domain.InProgress)))
	ns.Tasks[domain.Done].SetItems(convertTasksToListItems(project.GetTasksByStatus(domain.Done)))

	// Update selection states after refresh
	ns.Tasks[domain.NotStarted].SetItems(styles.UpdateTaskSelection(ns.Tasks[domain.NotStarted].Items(), notStartedIndex, ns.activeListIndex == domain.NotStarted))
	ns.Tasks[domain.InProgress].SetItems(styles.UpdateTaskSelection(ns.Tasks[domain.InProgress].Items(), inProgressIndex, ns.activeListIndex == domain.InProgress))
	ns.Tasks[domain.Done].SetItems(styles.UpdateTaskSelection(ns.Tasks[domain.Done].Items(), doneIndex, ns.activeListIndex == domain.Done))

	// PERFORMANCE: Clear all dirty flags after full update
	ns.clearAllDirtyFlags()
}

// GetActiveList returns the currently active list model
func (ns *NavigationState) GetActiveList() *list.Model {
	return &ns.Tasks[ns.activeListIndex]
}

// PERFORMANCE: Dirty flag management methods for incremental updates

// MarkListDirty marks a specific status list as needing updates
func (ns *NavigationState) MarkListDirty(status domain.Status) {
	if ns.dirtyFlags == nil {
		ns.dirtyFlags = make(map[domain.Status]bool)
	}
	ns.dirtyFlags[status] = true
}

// MarkAllListsDirty marks all status lists as needing updates
func (ns *NavigationState) MarkAllListsDirty() {
	if ns.dirtyFlags == nil {
		ns.dirtyFlags = make(map[domain.Status]bool)
	}
	ns.dirtyFlags[domain.NotStarted] = true
	ns.dirtyFlags[domain.InProgress] = true
	ns.dirtyFlags[domain.Done] = true
}

// clearAllDirtyFlags clears all dirty flags
func (ns *NavigationState) clearAllDirtyFlags() {
	if ns.dirtyFlags != nil {
		ns.dirtyFlags = make(map[domain.Status]bool)
	}
}

// IsListDirty checks if a specific status list needs updates
func (ns *NavigationState) IsListDirty(status domain.Status) bool {
	if ns.dirtyFlags == nil {
		return false
	}
	return ns.dirtyFlags[status]
}

// UpdateDirtyLists updates only the lists marked as dirty (PERFORMANCE optimization)
func (ns *NavigationState) UpdateDirtyLists(project *domain.Project, taskService *services.TaskService) {
	if project == nil {
		return
	}

	// If no dirty flags, treat as all dirty (fallback behavior)
	if ns.dirtyFlags == nil || len(ns.dirtyFlags) == 0 {
		ns.UpdateTaskLists(project, taskService)
		return
	}

	// Get all tasks once ( PERFORMANCE: single DB query )
	allTasks, err := taskService.GetTasksByProject(project.ID)
	if err != nil {
		// Handle error - clear dirty flags to avoid infinite loops
		ns.clearAllDirtyFlags()
		return
	}

	// Update project tasks
	project.Tasks = allTasks

	// Save selection states for lists that will be updated
	selections := map[domain.Status]int{
		domain.NotStarted: ns.Tasks[domain.NotStarted].Index(),
		domain.InProgress: ns.Tasks[domain.InProgress].Index(),
		domain.Done:       ns.Tasks[domain.Done].Index(),
	}

	// Update only dirty lists
	for status, isDirty := range ns.dirtyFlags {
		if isDirty {
			ns.Tasks[status].SetItems(convertTasksToListItems(project.GetTasksByStatus(status)))
			// Update selection state for the updated list
			ns.Tasks[status].SetItems(styles.UpdateTaskSelection(
				ns.Tasks[status].Items(),
				selections[status],
				ns.activeListIndex == status,
			))
		}
	}

	// Clear dirty flags after processing
	ns.clearAllDirtyFlags()
}

// UpdateListSizes updates the sizes of all task lists
func (ns *NavigationState) UpdateListSizes(width, height int) {
	columnWidth := max(20, (width-3)/3)

	ns.Tasks[domain.NotStarted].SetSize(columnWidth, height)
	ns.Tasks[domain.InProgress].SetSize(columnWidth, height)
	ns.Tasks[domain.Done].SetSize(columnWidth, height)
}

// UpdateActiveList updates only the active list with the given message
func (ns *NavigationState) UpdateActiveList(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	ns.Tasks[ns.activeListIndex], cmd = ns.Tasks[ns.activeListIndex].Update(msg)

	// Update selection state after cursor movement (up/down navigation)
	newItems := ns.Tasks[ns.activeListIndex].Items()
	ns.Tasks[ns.activeListIndex].SetItems(styles.UpdateTaskSelection(newItems, ns.Tasks[ns.activeListIndex].Index(), true))

	return cmd
}
