package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"kahn/internal/ui/input"
	"kahn/internal/ui/styles"
)

func (km *KahnModel) handleFormInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	action := km.inputHandler.HandleKeyMsg(msg, km)
	if action.Handled {
		return km, action.Cmd
	}

	comps := km.uiStateManager.FormState().GetActiveInputComponents()
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

func (km *KahnModel) handleProjectSwitch(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	navState := km.uiStateManager.NavigationState()
	confirmState := km.uiStateManager.ConfirmationState()

	if navState.IsShowingProjectSwitch() && km.inputHandler.GetMode() != input.ProjectSwitchMode {
		km.inputHandler.SetMode(input.ProjectSwitchMode)
	}
	if confirmState.IsShowingProjectDeleteConfirm() && km.inputHandler.GetMode() != input.ProjectDeleteConfirmMode {
		km.inputHandler.SetMode(input.ProjectDeleteConfirmMode)
	}

	if confirmState.IsShowingProjectDeleteConfirm() {
		switch msg.String() {
		case "y", "Y":
			return km.executeProjectDeletion(), nil
		case "n", "N", "esc":
			confirmState.ClearProjectDelete()
			return km, nil
		}
		return km, nil
	}

	switch msg.String() {
	case "esc":
		navState.HideProjectSwitch()
		return km, nil
	case "d":
		if km.projectManager.HasProjects() {
			confirmState.ShowProjectDeleteConfirm(km.projectManager.GetActiveProjectID())
		}
		return km, nil
	case "n":
		navState.HideProjectSwitch()
		km.uiStateManager.ShowProjectForm()
		km.inputHandler.SetMode(input.ProjectFormMode)
		km.inputHandler.SetFocusType(input.NameFocus)
		return km, nil
	case "j", "down":
		projects := km.projectManager.GetProjectsAsDomain()
		activeID := km.projectManager.GetActiveProjectID()
		for i, proj := range projects {
			if proj.ID == activeID {
				nextIndex := (i + 1) % len(projects)
				km.projectManager.SwitchToProject(projects[nextIndex].ID)
				return km, nil
			}
		}
	case "k", "up":
		projects := km.projectManager.GetProjectsAsDomain()
		activeID := km.projectManager.GetActiveProjectID()
		for i, proj := range projects {
			if proj.ID == activeID {
				prevIndex := (i - 1 + len(projects)) % len(projects)
				km.projectManager.SwitchToProject(projects[prevIndex].ID)
				return km, nil
			}
		}
	case "enter":
		navState.HideProjectSwitch()
		return km, nil
	default:
		if len(msg.String()) == 1 && msg.String()[0] >= '1' && msg.String()[0] <= '9' {
			projects := km.projectManager.GetProjectsAsDomain()
			index := int(msg.String()[0] - '1')
			if index < len(projects) {
				km.projectManager.SwitchToProject(projects[index].ID)
				navState.HideProjectSwitch()
			}
		}
		return km, nil
	}
	return km, nil
}

func (km *KahnModel) handleTaskDeleteConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	confirmState := km.uiStateManager.ConfirmationState()

	if km.inputHandler.GetMode() != input.TaskDeleteConfirmMode {
		km.inputHandler.SetMode(input.TaskDeleteConfirmMode)
	}

	switch msg.String() {
	case "y", "Y":
		return km.executeTaskDeletion(), nil
	case "n", "N", "esc":
		confirmState.ClearTaskDelete()
		return km, nil
	}
	return km, nil
}

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
		cmd := km.taskListManager.UpdateActiveList(msg)
		return km, cmd
	}

	// Handle other hotkeys
	switch keyStr {
	case "q":
		return km, tea.Quit
	case "n":
		km.uiStateManager.ShowTaskForm()
		return km, nil
	case "p":
		km.uiStateManager.ShowProjectSwitcher()
		return km, nil
	case "e":
		if selectedItem := km.taskListManager.GetActiveList().SelectedItem(); selectedItem != nil {
			if taskWrapper, ok := selectedItem.(styles.TaskWithTitle); ok {
				km.uiStateManager.ShowTaskEditForm(taskWrapper.ID, taskWrapper.Name, taskWrapper.Desc, taskWrapper.Priority, taskWrapper.Type)
			}
		}
		return km, nil
	case "d":
		if selectedItem := km.taskListManager.GetActiveList().SelectedItem(); selectedItem != nil {
			if taskWrapper, ok := selectedItem.(styles.TaskWithTitle); ok {
				km.uiStateManager.ShowTaskDeleteConfirm(taskWrapper.ID)
			}
		}
		return km, nil
	case " ":
		if selectedItem := km.taskListManager.GetActiveList().SelectedItem(); selectedItem != nil {
			if taskWrapper, ok := selectedItem.(styles.TaskWithTitle); ok {
				km.MoveTaskToNextStatus(taskWrapper.ID)
			}
		}
		return km, nil
	case "backspace":
		if selectedItem := km.taskListManager.GetActiveList().SelectedItem(); selectedItem != nil {
			if taskWrapper, ok := selectedItem.(styles.TaskWithTitle); ok {
				km.MoveTaskToPreviousStatus(taskWrapper.ID)
			}
		}
		return km, nil
	}
	return km, nil
}

func (km *KahnModel) handleResize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	km.width = msg.Width
	km.height = msg.Height

	frameWidth, frameHeight := styles.DefaultStyle.GetFrameSize()
	availableWidth := msg.Width - frameWidth/2
	availableHeight := msg.Height - frameHeight

	// Calculate project footer height dynamically
	activeProj := km.GetActiveProject()
	projectFooterHeight := 0
	if activeProj != nil {
		projectFooter := km.board.GetRenderer().RenderProjectFooter(activeProj, availableWidth-3, km.version)
		projectFooterHeight = lipgloss.Height(projectFooter)
	}

	// Set list heights accounting for both frame and project footer
	listHeight := availableHeight - projectFooterHeight
	km.taskListManager.UpdateListSizes(availableWidth, listHeight)
	return km, nil
}
