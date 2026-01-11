package app

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"kahn/internal/domain"
	"kahn/internal/services"
	"kahn/internal/ui/styles"
)

// TaskListManager handles all task list operations and state management
type TaskListManager struct {
	navState    *NavigationState
	taskService *services.TaskService
}

// NewTaskListManager creates a new task list manager
func NewTaskListManager(navState *NavigationState, taskService *services.TaskService) *TaskListManager {
	return &TaskListManager{
		navState:    navState,
		taskService: taskService,
	}
}

// UpdateTaskLists updates all task lists from the project
func (tlm *TaskListManager) UpdateTaskLists(project *domain.Project) {
	if project == nil {
		return
	}

	// Save current selection states before refreshing
	notStartedIndex := tlm.navState.Tasks[domain.NotStarted].Index()
	inProgressIndex := tlm.navState.Tasks[domain.InProgress].Index()
	doneIndex := tlm.navState.Tasks[domain.Done].Index()

	allTasks, err := tlm.taskService.GetTasksByProject(project.ID)
	if err != nil {
		// Handle error - set empty lists
		tlm.navState.Tasks[domain.NotStarted].SetItems([]list.Item{})
		tlm.navState.Tasks[domain.InProgress].SetItems([]list.Item{})
		tlm.navState.Tasks[domain.Done].SetItems([]list.Item{})
		return
	}

	project.Tasks = allTasks

	tlm.navState.Tasks[domain.NotStarted].SetItems(convertTasksToListItems(project.GetTasksByStatus(domain.NotStarted)))
	tlm.navState.Tasks[domain.InProgress].SetItems(convertTasksToListItems(project.GetTasksByStatus(domain.InProgress)))
	tlm.navState.Tasks[domain.Done].SetItems(convertTasksToListItems(project.GetTasksByStatus(domain.Done)))

	// Update selection states after refresh
	tlm.navState.Tasks[domain.NotStarted].SetItems(styles.UpdateTaskSelection(tlm.navState.Tasks[domain.NotStarted].Items(), notStartedIndex, tlm.navState.activeListIndex == domain.NotStarted))
	tlm.navState.Tasks[domain.InProgress].SetItems(styles.UpdateTaskSelection(tlm.navState.Tasks[domain.InProgress].Items(), inProgressIndex, tlm.navState.activeListIndex == domain.InProgress))
	tlm.navState.Tasks[domain.Done].SetItems(styles.UpdateTaskSelection(tlm.navState.Tasks[domain.Done].Items(), doneIndex, tlm.navState.activeListIndex == domain.Done))

	tlm.clearAllDirtyFlags()
}

// UpdateDirtyLists updates only the lists marked as dirty
func (tlm *TaskListManager) UpdateDirtyLists(project *domain.Project) {
	if project == nil {
		return
	}

	// If no dirty flags, treat as all dirty (fallback behavior)
	if tlm.navState.dirtyFlags == nil || len(tlm.navState.dirtyFlags) == 0 {
		tlm.UpdateTaskLists(project)
		return
	}

	allTasks, err := tlm.taskService.GetTasksByProject(project.ID)
	if err != nil {
		// Handle error - clear dirty flags to avoid infinite loops
		tlm.clearAllDirtyFlags()
		return
	}

	// Update project tasks
	project.Tasks = allTasks

	selections := map[domain.Status]int{
		domain.NotStarted: tlm.navState.Tasks[domain.NotStarted].Index(),
		domain.InProgress: tlm.navState.Tasks[domain.InProgress].Index(),
		domain.Done:       tlm.navState.Tasks[domain.Done].Index(),
	}

	for status, isDirty := range tlm.navState.dirtyFlags {
		if isDirty {
			tlm.navState.Tasks[status].SetItems(convertTasksToListItems(project.GetTasksByStatus(status)))
			// Update selection state for the updated list
			tlm.navState.Tasks[status].SetItems(styles.UpdateTaskSelection(
				tlm.navState.Tasks[status].Items(),
				selections[status],
				tlm.navState.activeListIndex == status,
			))
		}
	}

	// Clear dirty flags after processing
	tlm.clearAllDirtyFlags()
}

// UpdateTaskListsConditional updates lists using dirty flags if available, otherwise updates all
func (tlm *TaskListManager) UpdateTaskListsConditional(project *domain.Project) {
	if tlm.navState.dirtyFlags != nil && len(tlm.navState.dirtyFlags) > 0 {
		tlm.UpdateDirtyLists(project)
	} else {
		tlm.UpdateTaskLists(project)
	}
}

// MarkListDirty marks a specific list as needing refresh
func (tlm *TaskListManager) MarkListDirty(status domain.Status) {
	if tlm.navState.dirtyFlags == nil {
		tlm.navState.dirtyFlags = make(map[domain.Status]bool)
	}
	tlm.navState.dirtyFlags[status] = true
}

// MarkAllListsDirty marks all lists as needing refresh
func (tlm *TaskListManager) MarkAllListsDirty() {
	if tlm.navState.dirtyFlags == nil {
		tlm.navState.dirtyFlags = make(map[domain.Status]bool)
	}
	tlm.navState.dirtyFlags[domain.NotStarted] = true
	tlm.navState.dirtyFlags[domain.InProgress] = true
	tlm.navState.dirtyFlags[domain.Done] = true
}

// GetTaskItems returns the list items for a specific status
func (tlm *TaskListManager) GetTaskItems(status domain.Status) []list.Item {
	return tlm.navState.Tasks[status].Items()
}

// GetActiveListIndex returns the currently active list index
func (tlm *TaskListManager) GetActiveListIndex() domain.Status {
	return tlm.navState.GetActiveListIndex()
}

// GetActiveList returns the currently active list
func (tlm *TaskListManager) GetActiveList() *list.Model {
	return tlm.navState.GetActiveList()
}

// NextList switches to the next list in the navigation
func (tlm *TaskListManager) NextList() {
	tlm.navState.NextList()
}

// PrevList switches to the previous list in the navigation
func (tlm *TaskListManager) PrevList() {
	tlm.navState.PrevList()
}

func (tlm *TaskListManager) UpdateActiveList(msg tea.Msg) tea.Cmd {
	return tlm.navState.UpdateActiveList(msg)
}

func (tlm *TaskListManager) UpdateListSizes(width, height int) {
	tlm.navState.UpdateListSizes(width, height)
}

func (tlm *TaskListManager) clearAllDirtyFlags() {
	if tlm.navState.dirtyFlags != nil {
		tlm.navState.dirtyFlags = make(map[domain.Status]bool)
	}
}
