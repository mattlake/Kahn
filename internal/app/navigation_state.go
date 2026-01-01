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

	// Convert tasks to list items and update
	notStartedTasks, err := taskService.GetTasksByStatus(project.ID, domain.NotStarted)
	if err != nil {
		// Handle error - maybe set empty lists
		ns.Tasks[domain.NotStarted].SetItems([]list.Item{})
	} else {
		project.Tasks = notStartedTasks
		ns.Tasks[domain.NotStarted].SetItems(convertTasksToListItems(project.GetTasksByStatus(domain.NotStarted)))
	}

	inProgressTasks, err := taskService.GetTasksByStatus(project.ID, domain.InProgress)
	if err != nil {
		ns.Tasks[domain.InProgress].SetItems([]list.Item{})
	} else {
		project.Tasks = inProgressTasks
		ns.Tasks[domain.InProgress].SetItems(convertTasksToListItems(project.GetTasksByStatus(domain.InProgress)))
	}

	doneTasks, err := taskService.GetTasksByStatus(project.ID, domain.Done)
	if err != nil {
		ns.Tasks[domain.Done].SetItems([]list.Item{})
	} else {
		project.Tasks = doneTasks
		ns.Tasks[domain.Done].SetItems(convertTasksToListItems(project.GetTasksByStatus(domain.Done)))
	}

	// Update selection states after refresh
	ns.Tasks[domain.NotStarted].SetItems(styles.UpdateTaskSelection(ns.Tasks[domain.NotStarted].Items(), notStartedIndex, ns.activeListIndex == domain.NotStarted))
	ns.Tasks[domain.InProgress].SetItems(styles.UpdateTaskSelection(ns.Tasks[domain.InProgress].Items(), inProgressIndex, ns.activeListIndex == domain.InProgress))
	ns.Tasks[domain.Done].SetItems(styles.UpdateTaskSelection(ns.Tasks[domain.Done].Items(), doneIndex, ns.activeListIndex == domain.Done))
}

// GetActiveList returns the currently active list model
func (ns *NavigationState) GetActiveList() *list.Model {
	return &ns.Tasks[ns.activeListIndex]
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
