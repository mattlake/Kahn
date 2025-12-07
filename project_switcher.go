package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderProjectSwitcher() string {
	if len(m.Projects) == 0 {
		return m.renderNoProjectsMessage()
	}

	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorMauve)).
		Bold(true).
		Align(lipgloss.Center).
		Width(50).
		Render("Select Project")

	var projectItems []string
	for i, proj := range m.Projects {
		prefix := " "
		if proj.ID == m.ActiveProjectID {
			prefix = "►"
		}

		color := proj.Color
		if color == "" {
			color = ColorBlue
		}

		itemStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(color))

		if proj.ID == m.ActiveProjectID {
			itemStyle = itemStyle.Background(lipgloss.Color(ColorSurface1))
		}

		item := itemStyle.Render(fmt.Sprintf("%s %d. %s", prefix, i+1, proj.Name))

		projectItems = append(projectItems, item)
	}

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorSubtext1)).
		Align(lipgloss.Center).
		Width(50).
		Render("[↑/↓] Navigate • [Enter] Select • [n] New • [d] Delete • [esc] Cancel")

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

	form := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(ColorMauve)).
		Padding(2, 3).
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
	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorMauve)).
		Bold(true).
		Align(lipgloss.Center).
		Width(50).
		Render("No Projects")

	message := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorText)).
		Align(lipgloss.Center).
		Width(50).
		Render("Create your first project to get started")

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorSubtext1)).
		Align(lipgloss.Center).
		Width(50).
		Render("[n] New Project • [esc] Cancel")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		title,
		"",
		message,
		"",
		instructions,
	)

	form := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(ColorMauve)).
		Padding(2, 3).
		Width(60).
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
