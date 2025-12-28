package main

import (
	"github.com/charmbracelet/lipgloss"
	"kahn/pkg/colors"
)

// Lipgloss styles with Catppuccin colors
var defaultStyle = lipgloss.NewStyle().
	Margin(1, 2).
	Border(lipgloss.HiddenBorder()).
	Padding(1, 2)

var focusedStyle = lipgloss.NewStyle().
	Margin(1, 2).
	Border(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color(colors.Mauve)).
	Padding(1, 2)
