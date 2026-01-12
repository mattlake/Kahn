package app

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"kahn/internal/domain"
	"kahn/internal/services"
	"kahn/internal/ui/styles"
)

type NavigationState struct {
	showProjectSwitch bool
	activeListIndex   domain.Status
	Tasks             []list.Model

	dirtyFlags map[domain.Status]bool
}

func NewNavigationState(tasks []list.Model) *NavigationState {
	return &NavigationState{
		Tasks:           tasks,
		activeListIndex: domain.NotStarted,
	}
}

func (ns *NavigationState) ShowProjectSwitch() {
	ns.showProjectSwitch = true
}

func (ns *NavigationState) HideProjectSwitch() {
	ns.showProjectSwitch = false
}

// IsShowingProjectSwitch returns whether project switcher is shown
func (ns *NavigationState) IsShowingProjectSwitch() bool {
	return ns.showProjectSwitch
}

func (ns *NavigationState) GetActiveListIndex() domain.Status {
	return ns.activeListIndex
}

func (ns *NavigationState) SetActiveListIndex(index domain.Status) {
	ns.activeListIndex = index
}

// switchToList handles the common logic of switching from current list to a new list
func (ns *NavigationState) switchToList(newIndex domain.Status) {
	// Switch old active list to inactive delegate
	oldActiveIndex := ns.activeListIndex
	ns.Tasks[oldActiveIndex].SetDelegate(styles.NewInactiveListDelegate())

	// Update selection state for old list (no longer active)
	oldItems := ns.Tasks[oldActiveIndex].Items()
	ns.Tasks[oldActiveIndex].SetItems(styles.UpdateTaskSelection(oldItems, ns.Tasks[oldActiveIndex].Index(), false))

	// Move to new list
	ns.activeListIndex = newIndex

	// Switch new active list to active delegate
	ns.Tasks[ns.activeListIndex].SetDelegate(styles.NewActiveListDelegate())

	// Update selection state for new list (now active)
	newItems := ns.Tasks[ns.activeListIndex].Items()
	ns.Tasks[ns.activeListIndex].SetItems(styles.UpdateTaskSelection(newItems, ns.Tasks[ns.activeListIndex].Index(), true))

	// Update title styles to reflect new focus
	styles.ApplyFocusedTitleStyles(ns.Tasks[:], ns.activeListIndex)
}

func (ns *NavigationState) NextList() {
	var nextIndex domain.Status
	if ns.activeListIndex == domain.Done {
		nextIndex = domain.NotStarted
	} else {
		nextIndex = ns.activeListIndex + 1
	}
	ns.switchToList(nextIndex)
}

func (ns *NavigationState) PrevList() {
	var prevIndex domain.Status
	if ns.activeListIndex == domain.NotStarted {
		prevIndex = domain.Done
	} else {
		prevIndex = ns.activeListIndex - 1
	}
	ns.switchToList(prevIndex)
}

func (ns *NavigationState) UpdateTaskLists(project *domain.Project, taskService *services.TaskService) {
	if project == nil {
		return
	}

	// Save current selection states before refreshing
	notStartedIndex := ns.Tasks[domain.NotStarted].Index()
	inProgressIndex := ns.Tasks[domain.InProgress].Index()
	doneIndex := ns.Tasks[domain.Done].Index()

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

	ns.clearAllDirtyFlags()
}

func (ns *NavigationState) GetActiveList() *list.Model {
	return &ns.Tasks[ns.activeListIndex]
}

func (ns *NavigationState) MarkListDirty(status domain.Status) {
	if ns.dirtyFlags == nil {
		ns.dirtyFlags = make(map[domain.Status]bool)
	}
	ns.dirtyFlags[status] = true
}

func (ns *NavigationState) MarkAllListsDirty() {
	if ns.dirtyFlags == nil {
		ns.dirtyFlags = make(map[domain.Status]bool)
	}
	ns.dirtyFlags[domain.NotStarted] = true
	ns.dirtyFlags[domain.InProgress] = true
	ns.dirtyFlags[domain.Done] = true
}

func (ns *NavigationState) clearAllDirtyFlags() {
	if ns.dirtyFlags != nil {
		ns.dirtyFlags = make(map[domain.Status]bool)
	}
}

func (ns *NavigationState) IsListDirty(status domain.Status) bool {
	if ns.dirtyFlags == nil {
		return false
	}
	return ns.dirtyFlags[status]
}

func (ns *NavigationState) UpdateDirtyLists(project *domain.Project, taskService *services.TaskService) {
	if project == nil {
		return
	}

	// If no dirty flags, treat as all dirty (fallback behavior)
	if ns.dirtyFlags == nil || len(ns.dirtyFlags) == 0 {
		ns.UpdateTaskLists(project, taskService)
		return
	}

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

func (ns *NavigationState) UpdateListSizes(width, height int) {
	columnWidth := max(20, (width-3)/3)

	ns.Tasks[domain.NotStarted].SetSize(columnWidth, height)
	ns.Tasks[domain.InProgress].SetSize(columnWidth, height)
	ns.Tasks[domain.Done].SetSize(columnWidth, height)
}

func (ns *NavigationState) UpdateActiveList(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	ns.Tasks[ns.activeListIndex], cmd = ns.Tasks[ns.activeListIndex].Update(msg)

	// Update selection state after cursor movement (up/down navigation)
	newItems := ns.Tasks[ns.activeListIndex].Items()
	ns.Tasks[ns.activeListIndex].SetItems(styles.UpdateTaskSelection(newItems, ns.Tasks[ns.activeListIndex].Index(), true))

	return cmd
}

// UpdateTaskListsConditional updates lists using dirty flags if available, otherwise updates all
func (ns *NavigationState) UpdateTaskListsConditional(project *domain.Project, taskService *services.TaskService) {
	if ns.dirtyFlags != nil && len(ns.dirtyFlags) > 0 {
		ns.UpdateDirtyLists(project, taskService)
	} else {
		ns.UpdateTaskLists(project, taskService)
	}
}

// GetTaskItems returns the list items for a specific status
func (ns *NavigationState) GetTaskItems(status domain.Status) []list.Item {
	return ns.Tasks[status].Items()
}
