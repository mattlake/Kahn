package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"kahn/pkg/colors"
)

func initializeProjectInputs() (textinput.Model, textinput.Model) {
	name := textinput.New()
	name.Placeholder = "Project name"
	name.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Subtext0))
	name.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Text))
	name.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Mauve))
	name.Focus()
	name.CharLimit = 50
	name.Width = 40

	desc := textinput.New()
	desc.Placeholder = "Project description"
	desc.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Subtext0))
	desc.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Text))
	desc.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Mauve))
	desc.CharLimit = 200
	desc.Width = 40

	return name, desc
}

func (m Model) renderProjectForm() string {
	formTitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Mauve)).
		Bold(true).
		Align(lipgloss.Center).
		Width(50).
		Render("New Project")

	nameLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Text)).
		Bold(true).
		Render("Project Name:")

	descLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Text)).
		Bold(true).
		Render("Description:")

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Subtext1)).
		Align(lipgloss.Center).
		Width(50).
		Render("Tab: Switch fields • Enter: Create • Esc: Cancel")

	// Highlight focused input
	nameField := m.projNameInput.View()
	descField := m.projDescInput.View()

	if m.focusedProjInput == 0 {
		nameField = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(colors.Mauve)).
			Padding(0, 1).
			Render(nameField)
		descField = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(colors.Overlay1)).
			Padding(0, 1).
			Render(descField)
	} else {
		nameField = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(colors.Overlay1)).
			Padding(0, 1).
			Render(nameField)
		descField = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(colors.Mauve)).
			Padding(0, 1).
			Render(descField)
	}

	formContent := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		formTitle,
		"",
		nameLabel,
		nameField,
		"",
		descLabel,
		descField,
		"",
		instructions,
	)

	form := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(colors.Mauve)).
		Padding(2, 3).
		Width(60).
		Height(20).
		Render(formContent)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		form,
	)
}
