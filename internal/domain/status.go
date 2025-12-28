package domain

type Status int

const (
	NotStarted Status = iota
	InProgress
	Done
)

func (s Status) ToString() string {
	switch s {
	case NotStarted:
		return "Not Started"
	case InProgress:
		return "In Progress"
	case Done:
		return "Done"
	default:
		return "Placeholder"
	}
}
