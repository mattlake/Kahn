package app

import (
	"kahn/internal/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (km *KahnModel) handleFormInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	comps := km.uiStateManager.FormState().GetActiveInputComponents()

	switch msg.String() {
	case "esc":
		km.uiStateManager.HideAllStates()
		return km, nil
	case "tab":
		return km.handleTabKey(), nil
	case "ctrl+enter":
		// Ctrl+Enter always submits from any field
		if err := km.SubmitCurrentForm(); err != nil {
			// Validation failed - stay in form mode, show inline error
			return km, nil
		}

		// Success - exit form mode
		km.CancelCurrentForm()
		return km, nil
	case "enter":
		// If description field is focused, allow newlines in textarea
		if comps.FocusedField == 1 { // Description field (field index 1)
			updatedDesc, cmd := comps.DescInput.Update(msg)
			comps.DescInput = updatedDesc
			return km, cmd
		}

		// Otherwise, Enter from other fields means submit
		if err := km.SubmitCurrentForm(); err != nil {
			// Validation failed - stay in form mode, show inline error
			return km, nil
		}

		// Success - exit form mode
		km.CancelCurrentForm()
		return km, nil
	case "up", "down":
		// Handle priority/type/blocked_by cycling when those fields are focused
		if comps.IsTaskForm() {
			if comps.FocusedField == 2 { // Priority field focused
				if msg.String() == "up" {
					comps.CyclePriorityUp()
				} else {
					comps.CyclePriorityDown()
				}
				return km, nil
			} else if comps.FocusedField == 3 { // Type field focused
				if msg.String() == "up" {
					comps.CycleTypeUp()
				} else {
					comps.CycleTypeDown()
				}
				return km, nil
			} else if comps.FocusedField == 4 { // BlockedBy field focused
				if msg.String() == "up" {
					comps.CycleBlockedByUp()
				} else {
					comps.CycleBlockedByDown()
				}
				return km, nil
			}
		}
		// Let textinput/textarea handle for other fields
	default:
		// Clear any previous errors when user types
		km.ClearFormError()
	}

	// Update the focused input field
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

func (km *KahnModel) handleTabKey() tea.Model {
	comps := km.uiStateManager.FormState().GetActiveInputComponents()

	switch comps.FocusedField {
	case 0: // Name -> Description
		comps.FocusDesc()
		comps.BlurName()
	case 1: // Description -> Priority (for task forms) or Name (for project forms)
		if comps.IsTaskForm() {
			// Task forms: Description -> Priority
			comps.FocusPriority()
			comps.BlurDesc()
		} else {
			// Project forms: Description -> Name (cycle back)
			comps.FocusName()
			comps.BlurDesc()
		}
	case 2: // Priority -> Type (only for task forms)
		comps.FocusType()
		comps.BlurPriority()
	case 3: // Type -> BlockedBy (only for task forms)
		comps.FocusBlockedBy()
		comps.BlurType()
	case 4: // BlockedBy -> Name (only for task forms, cycle back)
		comps.FocusName()
		comps.BlurBlockedBy()
	default:
		// Fallback to name focus
		comps.FocusName()
	}
	return km
}

func (km *KahnModel) handleProjectSwitch(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	navState := km.uiStateManager.NavigationState()
	confirmState := km.uiStateManager.ConfirmationState()

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
	// Handle navigation keys by updating the active list
	keyStr := msg.String()
	if keyStr == "up" || keyStr == "down" || keyStr == "j" || keyStr == "k" {
		cmd := km.navState.UpdateActiveList(msg)
		return km, cmd
	}

	// Handle left/right navigation
	switch keyStr {
	case "l", "right":
		km.navState.NextList()
		return km, nil
	case "h", "left":
		km.navState.PrevList()
		return km, nil
	}

	// Handle other hotkeys
	switch keyStr {
	case "q":
		return km, tea.Quit
	case "n":
		km.ShowTaskForm()
		return km, nil
	case "p":
		km.uiStateManager.ShowProjectSwitcher()
		return km, nil
	case "e":
		if selectedItem := km.navState.GetActiveList().SelectedItem(); selectedItem != nil {
			if taskWrapper, ok := selectedItem.(styles.TaskWithTitle); ok {
				km.ShowTaskEditForm(taskWrapper.ID, taskWrapper.Name, taskWrapper.Desc, taskWrapper.Priority, taskWrapper.Type, taskWrapper.BlockedBy)
			}
		}
		return km, nil
	case "d":
		if selectedItem := km.navState.GetActiveList().SelectedItem(); selectedItem != nil {
			if taskWrapper, ok := selectedItem.(styles.TaskWithTitle); ok {
				km.uiStateManager.ShowTaskDeleteConfirm(taskWrapper.ID)
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
	km.navState.UpdateListSizes(availableWidth, listHeight)
	return km, nil
}
