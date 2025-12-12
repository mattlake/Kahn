package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

func initializeInputs() (textinput.Model, textinput.Model) {
	name := textinput.New()
	name.Placeholder = "Task name"
	name.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorSubtext0))
	name.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorText))
	name.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorMauve))
	name.Focus()
	name.CharLimit = 50
	name.Width = 40

	desc := textinput.New()
	desc.Placeholder = "Task description"
	desc.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorSubtext0))
	desc.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorText))
	desc.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorMauve))
	desc.CharLimit = 100
	desc.Width = 40

	return name, desc
}

func (m Model) renderForm() string {
	formTitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorMauve)).
		Bold(true).
		Align(lipgloss.Center).
		Width(50).
		Render("Add New Task")

	nameLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorText)).
		Bold(true).
		Render("Task Name:")

	descLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorText)).
		Bold(true).
		Render("Description:")

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorSubtext1)).
		Align(lipgloss.Center).
		Width(50).
		Render("Tab: Switch fields | Enter: Submit | Esc: Cancel")

	// Highlight focused input
	nameField := m.nameInput.View()
	descField := m.descInput.View()

	if m.focusedInput == 0 {
		nameField = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorMauve)).
			Padding(0, 1).
			Render(nameField)
		descField = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorOverlay1)).
			Padding(0, 1).
			Render(descField)
	} else {
		nameField = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorOverlay1)).
			Padding(0, 1).
			Render(nameField)
		descField = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorMauve)).
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
		BorderForeground(lipgloss.Color(ColorMauve)).
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

func (m Model) renderTaskEditForm() string {
	formTitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorMauve)).
		Bold(true).
		Align(lipgloss.Center).
		Width(50).
		Render("Edit Task")

	nameLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorText)).
		Bold(true).
		Render("Task Name:")

	descLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorText)).
		Bold(true).
		Render("Description:")

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorSubtext1)).
		Align(lipgloss.Center).
		Width(50).
		Render("Tab: Switch fields • Enter: Save • Esc: Cancel")

	// Highlight focused input
	nameField := m.nameInput.View()
	descField := m.descInput.View()

	if m.focusedInput == 0 {
		nameField = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorMauve)).
			Padding(0, 1).
			Render(nameField)
		descField = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorOverlay1)).
			Padding(0, 1).
			Render(descField)
	} else {
		nameField = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorOverlay1)).
			Padding(0, 1).
			Render(nameField)
		descField = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorMauve)).
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
		BorderForeground(lipgloss.Color(ColorMauve)).
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
