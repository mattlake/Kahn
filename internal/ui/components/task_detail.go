package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"kahn/internal/domain"
	"kahn/internal/ui/colors"
)

// TaskDetailRenderer handles rendering of the task detail modal
type TaskDetailRenderer struct {
	borderStyle   lipgloss.Style
	titleStyle    lipgloss.Style
	labelStyle    lipgloss.Style
	valueStyle    lipgloss.Style
	descStyle     lipgloss.Style
	footerStyle   lipgloss.Style
	contentWidth  int
	contentHeight int
}

// NewTaskDetailRenderer creates a new task detail renderer with default styling
func NewTaskDetailRenderer() *TaskDetailRenderer {
	return &TaskDetailRenderer{
		borderStyle: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(colors.Blue)).
			Padding(1, 2).
			Width(54),
		titleStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(colors.Text)).
			Bold(true).
			MarginBottom(1),
		labelStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(colors.Subtext0)).
			Width(14),
		valueStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(colors.Text)),
		descStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(colors.Text)).
			MarginTop(1).
			MarginBottom(1),
		footerStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color(colors.Subtext0)).
			MarginTop(1).
			Align(lipgloss.Center),
		contentWidth:  50,
		contentHeight: 20,
	}
}

// Render renders the task detail modal
func (tdr *TaskDetailRenderer) Render(task *domain.Task, blockingTask *domain.Task) string {
	if task == nil {
		return tdr.borderStyle.Render("No task selected")
	}

	var content strings.Builder

	// Title
	title := tdr.titleStyle.Render(task.Name)
	content.WriteString(title)
	content.WriteString("\n\n")

	// Type
	typeIcon := tdr.getTypeIcon(task.Type)
	typeLabel := tdr.labelStyle.Render("Type:")
	typeValue := tdr.valueStyle.Render(fmt.Sprintf("%s %s", typeIcon, task.Type.String()))
	content.WriteString(typeLabel + " " + typeValue)
	content.WriteString("\n")

	// Priority
	priorityLabel := tdr.labelStyle.Render("Priority:")
	priorityValue := tdr.getPriorityStyled(task.Priority)
	content.WriteString(priorityLabel + " " + priorityValue)
	content.WriteString("\n")

	// Status
	statusLabel := tdr.labelStyle.Render("Status:")
	statusValue := tdr.valueStyle.Render(task.Status.ToString())
	content.WriteString(statusLabel + " " + statusValue)
	content.WriteString("\n")

	// Description (if present)
	if strings.TrimSpace(task.Desc) != "" {
		content.WriteString("\n")
		descLabel := tdr.labelStyle.Render("Description:")
		content.WriteString(descLabel)
		content.WriteString("\n")
		descText := tdr.wrapText(task.Desc, tdr.contentWidth)
		content.WriteString(tdr.descStyle.Render(descText))
	}

	// Blocked by (if present)
	content.WriteString("\n")
	blockedLabel := tdr.labelStyle.Render("Blocked by:")
	blockedValue := tdr.getBlockedByText(blockingTask)
	content.WriteString(blockedLabel + " " + blockedValue)
	content.WriteString("\n")

	// Footer with keybindings
	footer := tdr.footerStyle.Render("[e] Edit  [Esc] Close")
	content.WriteString(footer)

	// Wrap everything in border
	modalContent := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(colors.Blue)).
		Padding(1, 2).
		Width(54).
		Render(content.String())

	return modalContent
}

// getTypeIcon returns the icon for the task type
func (tdr *TaskDetailRenderer) getTypeIcon(taskType domain.TaskType) string {
	switch taskType {
	case domain.Bug:
		return "󰃤"
	case domain.Feature:
		return "󱕣"
	default: // RegularTask
		return "󰄬"
	}
}

// getPriorityStyled returns a styled priority string with color
func (tdr *TaskDetailRenderer) getPriorityStyled(priority domain.Priority) string {
	var color string
	switch priority {
	case domain.Low:
		color = colors.Green
	case domain.Medium:
		color = colors.Peach
	case domain.High:
		color = colors.Red
	default:
		color = colors.Text
	}

	style := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
	return style.Render("● " + priority.String())
}

// getBlockedByText returns the text to display for the blocked by field
func (tdr *TaskDetailRenderer) getBlockedByText(blockingTask *domain.Task) string {
	if blockingTask == nil {
		return tdr.valueStyle.Render("None")
	}
	return tdr.valueStyle.Render(blockingTask.Name)
}

// wrapText wraps text to fit within the specified width
func (tdr *TaskDetailRenderer) wrapText(text string, width int) string {
	if len(text) <= width {
		return text
	}

	var wrapped strings.Builder
	words := strings.Fields(text)
	lineLength := 0

	for i, word := range words {
		wordLen := len(word)

		if lineLength+wordLen+1 > width {
			wrapped.WriteString("\n")
			wrapped.WriteString(word)
			lineLength = wordLen
		} else {
			if i > 0 {
				wrapped.WriteString(" ")
				lineLength++
			}
			wrapped.WriteString(word)
			lineLength += wordLen
		}
	}

	return wrapped.String()
}
