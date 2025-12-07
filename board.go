package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

func NewKahnModel() *KahnModel {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), 100, 0)
	defaultList.SetShowHelp(false)
	taskLists := []list.Model{defaultList, defaultList, defaultList}

	nameInput, descInput := initializeInputs()

	taskLists[NotStarted].Title = NotStarted.ToString()
	taskLists[NotStarted].SetItems(
		[]list.Item{
			Task{Name: "Task1", Desc: "Task1 Description", Status: NotStarted},
			Task{Name: "Task2", Desc: "Task2 Description", Status: NotStarted},
			Task{Name: "Task3", Desc: "Task3 Description", Status: NotStarted},
		})
	taskLists[NotStarted].Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorBlue)).
		Bold(true).
		Align(lipgloss.Center)

	taskLists[InProgress].Title = InProgress.ToString()
	taskLists[InProgress].SetItems(
		[]list.Item{
			Task{Name: "Task1", Desc: "Task1 Description", Status: InProgress},
		})
	taskLists[InProgress].Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorYellow)).
		Bold(true).
		Align(lipgloss.Center)

	taskLists[Done].Title = Done.ToString()
	taskLists[Done].SetItems(
		[]list.Item{
			Task{Name: "Task1", Desc: "Task1 Description", Status: Done},
			Task{Name: "Task2", Desc: "Task2 Description", Status: Done},
		})
	taskLists[Done].Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorGreen)).
		Bold(true).
		Align(lipgloss.Center)

	return &KahnModel{
		Tasks:     taskLists,
		nameInput: nameInput,
		descInput: descInput,
		width:     80, // default
		height:    24, // default
	}
}

func (m KahnModel) View() string {
	if m.showForm {
		return m.renderForm()
	}

	// Get current width of first list to use as column width
	columnWidth := m.Tasks[0].Width()

	// Apply width constraints to ensure equal columns
	notStartedView := defaultStyle.Width(columnWidth).Render(m.Tasks[NotStarted].View())
	inProgressView := defaultStyle.Width(columnWidth).Render(m.Tasks[InProgress].View())
	doneView := defaultStyle.Width(columnWidth).Render(m.Tasks[Done].View())

	focusedNotStartedView := focusedStyle.Width(columnWidth).Render(m.Tasks[NotStarted].View())
	focusedInProgressView := focusedStyle.Width(columnWidth).Render(m.Tasks[InProgress].View())
	focusedDoneView := focusedStyle.Width(columnWidth).Render(m.Tasks[Done].View())

	switch m.activeListIndex {
	case InProgress:
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			notStartedView,
			focusedInProgressView,
			doneView,
		)
	case Done:
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			notStartedView,
			inProgressView,
			focusedDoneView,
		)
	default:
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			focusedNotStartedView,
			inProgressView,
			doneView,
		)
	}
}
