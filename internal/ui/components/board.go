package components

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"kahn/internal/domain"
	"kahn/internal/ui/colors"
	"kahn/internal/ui/styles"
)

// Board handles main kanban board rendering
type Board struct {
	renderer BoardRenderer
}

// NewBoard creates a new board component
func NewBoard() *Board {
	return &Board{
		renderer: &BoardComponent{},
	}
}

// BoardComponent implements BoardRenderer interface
type BoardComponent struct{}

// RenderProjectHeader renders the top project header
func (b *BoardComponent) RenderProjectHeader(project *domain.Project, width int) string {
	if project == nil {
		return ""
	}

	// Create a more prominent project indicator
	projectLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Subtext1)).
		Render("Project:")

	projectNameText := lipgloss.NewStyle().
		Foreground(lipgloss.Color(project.Color)).
		Bold(true).
		Render(project.Name)

	helpText := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Subtext1)).
		Render("[p] Switch • [n] Add Task • [e] Edit Task • [d] Delete Task • [q] Quit")

	// Create a more prominent header with better visual hierarchy
	headerContent := lipgloss.JoinHorizontal(
		lipgloss.Left,
		projectLabel,
		lipgloss.NewStyle().Render(" "),
		projectNameText,
		lipgloss.NewStyle().Width(width-len(project.Name)-len("Project: ")-25).Render(""),
		helpText,
	)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(project.Color)).
		Padding(0, 1).
		Background(lipgloss.Color(colors.Surface0)).
		Width(width).
		Render(headerContent)
}

// RenderNoProjectsBoard renders empty state when no projects exist
func (b *BoardComponent) RenderNoProjectsBoard(width, height int) string {
	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Mauve)).
		Bold(true).
		Align(lipgloss.Center).
		Width(60).
		Render("No Projects")

	message := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Text)).
		Align(lipgloss.Center).
		Width(60).
		Render("Create your first project to get started")

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Subtext1)).
		Align(lipgloss.Center).
		Width(60).
		Render("[p] Create Project • [q] Quit")

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
		BorderForeground(lipgloss.Color(colors.Mauve)).
		Padding(2, 3).
		Width(70).
		Height(12).
		Render(content)

	return lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		form,
	)
}

// RenderTaskDeleteConfirm renders task deletion confirmation dialog
func (b *BoardComponent) RenderTaskDeleteConfirm(task *domain.Task, width, height int) string {
	if task == nil {
		return ""
	}

	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Red)).
		Bold(true).
		Align(lipgloss.Center).
		Width(60).
		Render("⚠️  Delete Task")

	taskName := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Text)).
		Bold(true).
		Render(task.Name)

	warningMessage := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Text)).
		Align(lipgloss.Center).
		Width(60).
		Render(fmt.Sprintf("Delete task \"%s\"?", taskName))

	subWarning := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Subtext1)).
		Align(lipgloss.Center).
		Width(60).
		Render("This action cannot be undone.")

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Subtext1)).
		Align(lipgloss.Center).
		Width(60).
		Render("[y] Yes, Delete • [n] No, Cancel")

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

	form := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(colors.Red)).
		Padding(2, 3).
		Width(70).
		Height(12).
		Render(content)

	return lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		form,
	)
}

// RenderBoard renders the main kanban board with three columns
func (b *BoardComponent) RenderBoard(project *domain.Project, taskLists [3]list.Model, activeListIndex domain.Status, width int) string {
	if project == nil {
		return ""
	}

	// Render project header
	projectHeader := b.RenderProjectHeader(project, width)

	columnWidth := taskLists[0].Width()

	notStartedView := styles.DefaultStyle.Width(columnWidth).Render(taskLists[domain.NotStarted].View())
	inProgressView := styles.DefaultStyle.Width(columnWidth).Render(taskLists[domain.InProgress].View())
	doneView := styles.DefaultStyle.Width(columnWidth).Render(taskLists[domain.Done].View())

	focusedNotStartedView := styles.FocusedStyle.Width(columnWidth).Render(taskLists[domain.NotStarted].View())
	focusedInProgressView := styles.FocusedStyle.Width(columnWidth).Render(taskLists[domain.InProgress].View())
	focusedDoneView := styles.FocusedStyle.Width(columnWidth).Render(taskLists[domain.Done].View())

	boardContent := ""
	switch activeListIndex {
	case domain.InProgress:
		boardContent = lipgloss.JoinHorizontal(
			lipgloss.Left,
			notStartedView,
			focusedInProgressView,
			doneView,
		)
	case domain.Done:
		boardContent = lipgloss.JoinHorizontal(
			lipgloss.Left,
			notStartedView,
			inProgressView,
			focusedDoneView,
		)
	default:
		boardContent = lipgloss.JoinHorizontal(
			lipgloss.Left,
			focusedNotStartedView,
			inProgressView,
			doneView,
		)
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		projectHeader,
		boardContent,
	)
}

// GetRenderer returns the board renderer interface
func (b *Board) GetRenderer() BoardRenderer {
	return b.renderer
}
