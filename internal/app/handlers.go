package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"kahn/internal/ui/input"
	"kahn/internal/ui/styles"
)

// handleFormInput processes input when a form is active
func (km *KahnModel) handleFormInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	action := km.inputHandler.HandleKeyMsg(msg, km)
	if action.Handled {
		return km, action.Cmd
	}

	comps := km.formState.GetActiveInputComponents()
	if comps.FocusedField == 0 {
		updatedName, cmd := comps.NameInput.Update(msg)
		comps.NameInput = updatedName
		return km, cmd
	} else {
		updatedDesc, cmd := comps.DescInput.Update(msg)
		comps.DescInput = updatedDesc
		return km, cmd
	}
}

// handleProjectSwitch processes input during project switching
func (km *KahnModel) handleProjectSwitch(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if km.navState.IsShowingProjectSwitch() && km.inputHandler.GetMode() != input.ProjectSwitchMode {
		km.inputHandler.SetMode(input.ProjectSwitchMode)
	}
	if km.confirmState.IsShowingProjectDeleteConfirm() && km.inputHandler.GetMode() != input.ProjectDeleteConfirmMode {
		km.inputHandler.SetMode(input.ProjectDeleteConfirmMode)
	}

	if km.confirmState.IsShowingProjectDeleteConfirm() {
		switch msg.String() {
		case "y", "Y":
			return km.executeProjectDeletion(), nil
		case "n", "N", "esc":
			km.confirmState.ClearProjectDelete()
			return km, nil
		}
		return km, nil
	}

	switch msg.String() {
	case "esc":
		km.navState.HideProjectSwitch()
		return km, nil
	case "d":
		if len(km.Projects) > 0 {
			km.confirmState.ShowProjectDeleteConfirm(km.ActiveProjectID)
		}
		return km, nil
	case "n":
		km.navState.HideProjectSwitch()
		km.formState.ShowProjectForm()
		km.inputHandler.SetMode(input.ProjectFormMode)
		km.inputHandler.SetFocusType(input.NameFocus)
		return km, nil
	case "j", "down":
		for i, proj := range km.Projects {
			if proj.ID == km.ActiveProjectID {
				nextIndex := (i + 1) % len(km.Projects)
				km.ActiveProjectID = km.Projects[nextIndex].ID
				km.navState.UpdateTaskLists(km.GetActiveProject(), km.taskService)
				return km, nil
			}
		}
	case "k", "up":
		for i, proj := range km.Projects {
			if proj.ID == km.ActiveProjectID {
				prevIndex := (i - 1 + len(km.Projects)) % len(km.Projects)
				km.ActiveProjectID = km.Projects[prevIndex].ID
				km.navState.UpdateTaskLists(km.GetActiveProject(), km.taskService)
				return km, nil
			}
		}
	case "enter":
		km.navState.HideProjectSwitch()
		return km, nil
	default:
		if len(msg.String()) == 1 && msg.String()[0] >= '1' && msg.String()[0] <= '9' {
			index := int(msg.String()[0] - '1')
			if index < len(km.Projects) {
				km.ActiveProjectID = km.Projects[index].ID
				km.navState.UpdateTaskLists(km.GetActiveProject(), km.taskService)
				km.navState.HideProjectSwitch()
			}
		}
		return km, nil
	}
	return km, nil
}

// handleTaskDeleteConfirm processes input during task deletion confirmation
func (km *KahnModel) handleTaskDeleteConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if km.inputHandler.GetMode() != input.TaskDeleteConfirmMode {
		km.inputHandler.SetMode(input.TaskDeleteConfirmMode)
	}

	switch msg.String() {
	case "y", "Y":
		return km.executeTaskDeletion(), nil
	case "n", "N", "esc":
		km.confirmState.ClearTaskDelete()
		return km, nil
	}
	return km, nil
}

// handleNormalMode processes input in normal browsing mode
func (km *KahnModel) handleNormalMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if km.inputHandler.GetMode() != input.NormalMode {
		km.inputHandler.SetMode(input.NormalMode)
	}

	// Let input handler try to handle the key first
	result := km.inputHandler.HandleKeyMsg(msg, km)

	// If input handler handled it, return its result
	if result.Handled {
		return km, result.Cmd
	}

	// Otherwise, handle navigation keys by updating the active list
	keyStr := msg.String()
	if keyStr == "up" || keyStr == "down" || keyStr == "j" || keyStr == "k" {
		cmd := km.navState.UpdateActiveList(msg)
		return km, cmd
	}

	// Handle other hotkeys
	switch keyStr {
	case "q":
		return km, tea.Quit
	case "n":
		km.formState.ShowTaskForm()
		return km, nil
	case "p":
		km.navState.ShowProjectSwitch()
		return km, nil
	case "e":
		if selectedItem := km.navState.GetActiveList().SelectedItem(); selectedItem != nil {
			if taskWrapper, ok := selectedItem.(styles.TaskWithTitle); ok {
				km.formState.ShowTaskEditForm(taskWrapper.ID, taskWrapper.Name, taskWrapper.Desc, taskWrapper.Priority)
			}
		}
		return km, nil
	case "d":
		if selectedItem := km.navState.GetActiveList().SelectedItem(); selectedItem != nil {
			if taskWrapper, ok := selectedItem.(styles.TaskWithTitle); ok {
				km.confirmState.ShowTaskDeleteConfirm(taskWrapper.ID)
			}
		}
		return km, nil
	case " ":
		if selectedItem := km.navState.GetActiveList().SelectedItem(); selectedItem != nil {
			if taskWrapper, ok := selectedItem.(styles.TaskWithTitle); ok {
				km.MoveTaskToNextStatus(taskWrapper.ID)
			}
		}
		return km, nil
	case "backspace":
		if selectedItem := km.navState.GetActiveList().SelectedItem(); selectedItem != nil {
			if taskWrapper, ok := selectedItem.(styles.TaskWithTitle); ok {
				km.MoveTaskToPreviousStatus(taskWrapper.ID)
			}
		}
		return km, nil
	}
	return km, nil
}

// handleResize processes window resize events
func (km *KahnModel) handleResize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	km.width = msg.Width
	km.height = msg.Height

	frameWidth, frameHeight := styles.DefaultStyle.GetFrameSize()
	availableWidth := msg.Width - frameWidth/2
	availableHeight := msg.Height - frameHeight

	// Calculate project header height dynamically
	activeProj := km.GetActiveProject()
	projectHeaderHeight := 0
	if activeProj != nil {
		projectHeader := km.board.GetRenderer().RenderProjectHeader(activeProj, availableWidth-3)
		projectHeaderHeight = lipgloss.Height(projectHeader)
	}

	// Set list heights accounting for both frame and project header
	listHeight := availableHeight - projectHeaderHeight
	km.navState.UpdateListSizes(availableWidth, listHeight)
	return km, nil
}
