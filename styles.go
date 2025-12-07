package main

import "github.com/charmbracelet/lipgloss"

// Lipgloss styles with Catppuccin colors
var defaultStyle = lipgloss.NewStyle().
	Margin(1, 2).
	Border(lipgloss.HiddenBorder()).
	Padding(1, 2)

var focusedStyle = lipgloss.NewStyle().
	Margin(1, 2).
	Border(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color(ColorMauve)).
	Padding(1, 2)
