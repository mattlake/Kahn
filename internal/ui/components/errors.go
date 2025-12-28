package components

import (
	"kahn/internal/domain"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"kahn/pkg/colors"
)

type ErrorMessage struct {
	Title   string
	Message string
}

type ErrorService struct {
	errors chan ErrorMessage
}

func NewErrorService() *ErrorService {
	return &ErrorService{
		errors: make(chan ErrorMessage, 10),
	}
}

func (es *ErrorService) HandleError(err error) ErrorMessage {
	var title, message string

	switch e := err.(type) {
	case *domain.ValidationError:
		title = "Validation Error"
		message = e.Message
	case *domain.RepositoryError:
		title = "Database Error"
		message = "An error occurred while accessing the database"
	default:
		title = "Error"
		if strings.Contains(err.Error(), "not found") {
			message = "The requested item was not found"
		} else if strings.Contains(err.Error(), "database") || strings.Contains(err.Error(), "sql") {
			message = "A database error occurred"
		} else {
			message = "An unexpected error occurred"
		}
	}

	return ErrorMessage{Title: title, Message: message}
}

func (es *ErrorService) ShowError(title, message string) tea.Cmd {
	return func() tea.Msg {
		return ErrorMessage{Title: title, Message: message}
	}
}

type ErrorDisplayModel struct {
	error   ErrorMessage
	visible bool
	width   int
	height  int
}

func NewErrorDisplayModel() *ErrorDisplayModel {
	return &ErrorDisplayModel{
		visible: false,
	}
}

func (edm *ErrorDisplayModel) Show(error ErrorMessage) {
	edm.error = error
	edm.visible = true
}

func (edm *ErrorDisplayModel) Hide() {
	edm.visible = false
}

func (edm *ErrorDisplayModel) IsVisible() bool {
	return edm.visible
}

func (edm *ErrorDisplayModel) Update(msg tea.Msg) tea.Cmd {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.Type {
		case tea.KeyEnter, tea.KeyEscape, tea.KeyCtrlC:
			edm.Hide()
			return nil
		}
	}
	return nil
}

func (edm *ErrorDisplayModel) View() string {
	if !edm.visible {
		return ""
	}

	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Red)).
		Bold(true).
		Align(lipgloss.Center).
		Width(60).
		Render(edm.error.Title)

	message := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Text)).
		Align(lipgloss.Center).
		Width(60).
		Render(edm.error.Message)

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Subtext1)).
		Align(lipgloss.Center).
		Width(60).
		Render("[Enter] Close")

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
		BorderForeground(lipgloss.Color(colors.Red)).
		Padding(2, 3).
		Width(70).
		Height(10).
		Render(content)

	return lipgloss.Place(
		edm.width, edm.height,
		lipgloss.Center, lipgloss.Center,
		form,
	)
}
