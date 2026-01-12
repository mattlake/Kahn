package components

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"kahn/internal/domain"
	"kahn/internal/ui/colors"
	"kahn/internal/ui/styles"
)

type Board struct {
	renderer BoardRenderer
}

func NewBoard() *Board {
	return &Board{
		renderer: &BoardComponent{},
	}
}

type BoardComponent struct{}

func (b *BoardComponent) RenderProjectFooter(project *domain.Project, width int, version string) string {
	if project == nil {
		return ""
	}

	projectLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Subtext1)).
		Render("Project:")

	projectNameText := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Green)).
		Bold(true).
		Render(project.Name)

	helpText := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Subtext1)).
		Render(fmt.Sprintf("Kahn %s | Nav: ←→/h/l | Move: space | Project: p | Add: n | Edit: e | Delete: d | Search: / | Quit: q", version))

	footerContent := lipgloss.JoinHorizontal(
		lipgloss.Left,
		projectLabel,
		lipgloss.NewStyle().Render(" "),
		projectNameText,
		lipgloss.NewStyle().Render(" | "),
		helpText,
	)

	return lipgloss.NewStyle().
		Margin(0, 0).
		Padding(0, 1).
		Width(width).
		Render(footerContent)
}

// RenderSearchBar renders the search input bar at the bottom when search is active
func (b *BoardComponent) RenderSearchBar(query string, matchCount int, width int) string {
	searchLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Mauve)).
		Bold(true).
		Render("Search:")

	queryText := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Text)).
		Render(query)

	// Add cursor indicator
	cursor := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Mauve)).
		Render("▊")

	matchInfo := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Subtext1)).
		Render(fmt.Sprintf("(%d matches)", matchCount))

	helpText := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Subtext1)).
		Render("[ESC] Clear search")

	searchContent := lipgloss.JoinHorizontal(
		lipgloss.Left,
		searchLabel,
		lipgloss.NewStyle().Render(" "),
		queryText,
		cursor,
		lipgloss.NewStyle().Render(" "),
		matchInfo,
		lipgloss.NewStyle().Render(" | "),
		helpText,
	)

	return lipgloss.NewStyle().
		Margin(0, 0).
		Padding(0, 1).
		Width(width).
		Render(searchContent)
}

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

func (b *BoardComponent) RenderTaskDeleteConfirmWithError(task *domain.Task, errorMessage string, width, height int) string {
	if task == nil {
		return ""
	}

	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Red)).
		Bold(true).
		Align(lipgloss.Center).
		Width(60).
		Render("⚠️  Delete Task")

	var content string

	if errorMessage != "" {
		// Show error message
		errorText := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff5555")).
			Align(lipgloss.Center).
			Width(60).
			Render("❌ " + errorMessage)

		okText := lipgloss.NewStyle().
			Foreground(lipgloss.Color(colors.Subtext1)).
			Align(lipgloss.Center).
			Width(60).
			Render("[ESC] Continue")

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

func (b *BoardComponent) RenderBoard(project *domain.Project, taskLists [3]list.Model, activeListIndex domain.Status, width int, version string, searchActive bool, searchQuery string, searchMatchCount int) string {
	if project == nil {
		return ""
	}

	// Render footer or search bar depending on search state
	var footer string
	if searchActive {
		footer = b.RenderSearchBar(searchQuery, searchMatchCount, width)
	} else {
		footer = b.RenderProjectFooter(project, width, version)
	}

	columnWidth := taskLists[0].Width()

	notStartedContent := taskLists[domain.NotStarted].View()
	inProgressContent := taskLists[domain.InProgress].View()
	doneContent := taskLists[domain.Done].View()

	notStartedView := styles.DefaultStyle.Width(columnWidth).Render(notStartedContent)
	inProgressView := styles.DefaultStyle.Width(columnWidth).Render(inProgressContent)
	doneView := styles.DefaultStyle.Width(columnWidth).Render(doneContent)

	focusedNotStartedView := styles.FocusedStyle.Width(columnWidth).Render(notStartedContent)
	focusedInProgressView := styles.FocusedStyle.Width(columnWidth).Render(inProgressContent)
	focusedDoneView := styles.FocusedStyle.Width(columnWidth).Render(doneContent)

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
		boardContent,
		footer,
	)
}

func (b *Board) GetRenderer() BoardRenderer {
	return b.renderer
}
