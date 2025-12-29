package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"kahn/internal/domain"
	"kahn/internal/ui/colors"
	"kahn/internal/ui/styles"
)

func (m Model) renderProjectSwitcher() string {
	// Show confirmation dialog if active
	if m.showProjectDeleteConfirm {
		return m.renderProjectDeleteConfirm()
	}

	if len(m.Projects) == 0 {
		return m.renderNoProjectsMessage()
	}

	dialogStyles := styles.GetDialogStyles()
	title := dialogStyles.Title.Width(50).Render("Select Project")

	var projectItems []string
	for i, proj := range m.Projects {
		prefix := " "
		if proj.ID == m.ActiveProjectID {
			prefix = "►"
		}

		color := proj.Color
		if color == "" {
			color = colors.Blue
		}

		itemStyles := styles.GetProjectItemStyle(color)

		var itemStyle lipgloss.Style
		if proj.ID == m.ActiveProjectID {
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
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		form,
	)
}

func (m Model) renderNoProjectsMessage() string {
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
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		form,
	)
}

func (m Model) renderProjectDeleteConfirm() string {
	// Find the project to delete
	var projectToDelete *domain.Project
	for _, proj := range m.Projects {
		if proj.ID == m.projectToDelete {
			projectToDelete = &proj
			break
		}
	}

	if projectToDelete == nil {
		// Fallback to active project if somehow projectToDelete is not set
		projectToDelete = m.GetActiveProject()
		if projectToDelete == nil {
			return m.renderProjectSwitcher() // Fallback to normal switcher
		}
	}

	deleteStyles := styles.GetDeleteConfirmStyles()
	dialogStyles := styles.GetDialogStyles()

	title := deleteStyles.Title.Width(60).Render("⚠️  Delete Project")

	projectName := dialogStyles.Message.
		Foreground(lipgloss.Color(projectToDelete.Color)).
		Bold(true).
		Render(projectToDelete.Name)

	warningMessage := deleteStyles.Message.Width(60).Render(fmt.Sprintf("Delete project \"%s\" and ALL its tasks?", projectName))
	subWarning := deleteStyles.Message.Width(60).Render("This action cannot be undone.")
	instructions := deleteStyles.Message.Width(60).Render("[y] Yes, Delete • [n] No, Cancel")

	content := lipgloss.JoinVertical(
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

	form := deleteStyles.Form.
		Width(70).
		Height(12).
		Render(content)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		form,
	)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
