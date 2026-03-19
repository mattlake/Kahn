package domain

type Status int

const (
	NotStarted Status = iota
	InProgress
	Done
)

// IsValidStatus returns true if v is a valid Status value.
func IsValidStatus(v int) bool {
	return v >= int(NotStarted) && v <= int(Done)
}

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
