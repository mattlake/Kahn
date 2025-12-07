package main

type Task struct {
	Name   string
	Desc   string
	Status Status
}

func (t Task) Title() string       { return t.Name }
func (t Task) Description() string { return t.Desc }
func (t Task) FilterValue() string { return t.Name }
