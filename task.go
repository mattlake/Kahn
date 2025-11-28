package main

type taskStatus int

const (
	todo taskStatus = iota
	inProgress
	done
)

/* CUSTOM ITEM */

type Task struct {
	status      taskStatus
	title       string
	description string
}

// Mkke Task implement the list.Item interface
func (t Task) FilterValue() string {
	return t.title
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}
