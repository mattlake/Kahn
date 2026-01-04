package components

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"kahn/internal/domain"
	"kahn/internal/ui/colors"
	"kahn/internal/ui/styles"
)

// ProjectSwitcher handles project selection and switching UI
type ProjectSwitcher struct {
	renderer ProjectSwitcherRenderer
}

// NewProjectSwitcher creates a new project switcher component
func NewProjectSwitcher() *ProjectSwitcher {
	return &ProjectSwitcher{
		renderer: &ProjectSwitcherComponent{},
	}
}

// ProjectSwitcherRenderer defines the interface for project switcher UI rendering
type ProjectSwitcherRenderer interface {
	RenderProjectSwitcher(projects []domain.Project, activeProjectID string, showDeleteConfirm bool, projectToDelete string, width, height int) string
	RenderNoProjectsMessage(width, height int) string
	RenderProjectDeleteConfirm(project *domain.Project, errorMessage string, width, height int) string
}

// ProjectSwitcherComponent implements ProjectSwitcherRenderer interface
type ProjectSwitcherComponent struct{}

// RenderProjectSwitcher renders the project selection dialog
func (psc *ProjectSwitcherComponent) RenderProjectSwitcher(projects []domain.Project, activeProjectID string, showDeleteConfirm bool, projectToDelete string, width, height int) string {
	// Show confirmation dialog if active
	if showDeleteConfirm {
		// Find the project to delete
		var projectToDeletePtr *domain.Project
		for _, proj := range projects {
			if proj.ID == projectToDelete {
				projectToDeletePtr = &proj
				break
			}
		}
		return psc.RenderProjectDeleteConfirm(projectToDeletePtr, "", width, height)
	}

	if len(projects) == 0 {
		return psc.RenderNoProjectsMessage(width, height)
	}

	dialogStyles := styles.GetDialogStyles()
	title := dialogStyles.Title.Width(50).Render("Select Project")

	var projectItems []string
	for i, proj := range projects {
		prefix := " "
		if proj.ID == activeProjectID {
			prefix = "►"
		}

		color := proj.Color
		if color == "" {
			color = colors.Blue
		}

		itemStyles := styles.GetProjectItemStyle(color)

		var itemStyle lipgloss.Style
		if proj.ID == activeProjectID {
			itemStyle = itemStyles.Active
		} else {
			itemStyle = itemStyles.Normal
		}

		item := itemStyle.Render(fmt.Sprintf("%s %d. %s", prefix, i+1, proj.Name))

		projectItems = append(projectItems, item)
	}

	instructions := dialogStyles.Instruction.Width(50).Render("[↑/↓] Navigate • [Enter] Select • [n] New • [d] Delete • [esc] Cancel")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		title,
		"",
	)

	for _, item := range projectItems {
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			content,
			item,
		)
	}

	content = lipgloss.JoinVertical(
		lipgloss.Left,
		content,
		"",
		instructions,
	)

	form := dialogStyles.Form.
		Width(60).
		Height(min(20, len(projectItems)+8)).
		Render(content)

	return lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		form,
	)
}

// RenderNoProjectsMessage renders the empty state when no projects exist
func (psc *ProjectSwitcherComponent) RenderNoProjectsMessage(width, height int) string {
	dialogStyles := styles.GetDialogStyles()

	title := dialogStyles.Title.Width(50).Render("No Projects")
	message := dialogStyles.Message.Width(50).Render("Create your first project to get started")
	instructions := dialogStyles.Instruction.Width(50).Render("[n] New Project • [esc] Cancel")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		title,
		"",
		message,
		"",
		instructions,
	)

	form := dialogStyles.Form.
		Width(60).
		Height(12).
		Render(content)

	return lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		form,
	)
}

// RenderProjectDeleteConfirm renders the project deletion confirmation dialog
func (psc *ProjectSwitcherComponent) RenderProjectDeleteConfirm(project *domain.Project, errorMessage string, width, height int) string {
	if project == nil {
		return ""
	}

	deleteStyles := styles.GetDeleteConfirmStyles()
	dialogStyles := styles.GetDialogStyles()

	title := deleteStyles.Title.Width(60).Render("⚠️  Delete Project")

	projectName := dialogStyles.Message.
		Foreground(lipgloss.Color(project.Color)).
		Bold(true).
		Render(project.Name)

	var content string

	if errorMessage != "" {
		// Show error message
		errorText := deleteStyles.Message.Width(60).Foreground(lipgloss.Color("#ff5555")).Render("❌ " + errorMessage)
		okText := deleteStyles.Message.Width(60).Render("[ESC] Continue")

		content = lipgloss.JoinVertical(
			lipgloss.Left,
			"",
			title,
			"",
			errorText,
			"",
			okText,
		)
	} else {
		// Show normal confirmation dialog
		warningMessage := deleteStyles.Message.Width(60).Render(fmt.Sprintf("Delete project \"%s\" and ALL its tasks?", projectName))
		subWarning := deleteStyles.Message.Width(60).Render("This action cannot be undone.")
		instructions := deleteStyles.Message.Width(60).Render("[y] Yes, Delete • [n] No, Cancel")

		content = lipgloss.JoinVertical(
			lipgloss.Left,
			"",
			title,
			"",
			warningMessage,
			"",
			subWarning,
			"",
			instructions,
		)
	}

	form := deleteStyles.Form.
		Width(70).
		Height(12).
		Render(content)

	return lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		form,
	)
}

// GetRenderer returns the project switcher renderer interface
func (ps *ProjectSwitcher) GetRenderer() ProjectSwitcherRenderer {
	return ps.renderer
}

// RenderSwitcher is a convenience method that matches the original Model method signature
func (ps *ProjectSwitcher) RenderSwitcher(projects []domain.Project, activeProjectID string, showDeleteConfirm bool, projectToDelete string, width, height int) string {
	return ps.renderer.RenderProjectSwitcher(projects, activeProjectID, showDeleteConfirm, projectToDelete, width, height)
}

// RenderSwitcherWithError renders the project switcher with error information
func (ps *ProjectSwitcher) RenderSwitcherWithError(projects []domain.Project, activeProjectID string, showDeleteConfirm bool, projectToDelete, errorMessage string, width, height int) string {
	// Show confirmation dialog if active
	if showDeleteConfirm {
		// Find the project to delete
		var projectToDeletePtr *domain.Project
		for _, proj := range projects {
			if proj.ID == projectToDelete {
				projectToDeletePtr = &proj
				break
			}
		}
		return ps.renderer.RenderProjectDeleteConfirm(projectToDeletePtr, errorMessage, width, height)
	}

	if len(projects) == 0 {
		return ps.renderer.RenderNoProjectsMessage(width, height)
	}

	return ps.renderer.RenderProjectSwitcher(projects, activeProjectID, showDeleteConfirm, projectToDelete, width, height)
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
