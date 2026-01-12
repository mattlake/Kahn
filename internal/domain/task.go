package domain

import (
	"fmt"
	"sort"
	"time"
)

type Task struct {
	IntID     int       `json:"int_id"`
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	Name      string    `json:"name"`
	Desc      string    `json:"desc"`
	Status    Status    `json:"status"`
	Type      TaskType  `json:"type"`
	BlockedBy *int      `json:"blocked_by,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Priority  Priority  `json:"priority,omitempty"`
}

type Priority int

const (
	Low Priority = iota
	Medium
	High
)

func (p Priority) String() string {
	switch p {
	case Low:
		return "Low"
	case Medium:
		return "Medium"
	case High:
		return "High"
	default:
		return "Medium"
	}
}

type TaskType int

const (
	RegularTask TaskType = iota
	Bug
	Feature
)

func (tt TaskType) String() string {
	switch tt {
	case RegularTask:
		return "Task"
	case Bug:
		return "Bug"
	case Feature:
		return "Feature"
	default:
		return "Task"
	}
}

// Validation length constants for tasks
const (
	MaxTaskNameLength        = 100
	MaxTaskDescriptionLength = 500
)

func NewTask(name, description, projectID string) *Task {
	now := time.Now()
	return &Task{
		ID:        generateTaskID(),
		ProjectID: projectID,
		Name:      name,
		Desc:      description,
		Status:    NotStarted,
		Type:      RegularTask,
		CreatedAt: now,
		UpdatedAt: now,
		Priority:  Low,
	}
}

func generateTaskID() string {
	return fmt.Sprintf("task_%d", time.Now().UnixNano())
}

func (t *Task) Validate() error {
	validator := NewFieldValidator()

	if err := validator.ValidateNotEmpty("name", t.Name, "task"); err != nil {
		return err
	}
	if err := validator.ValidateMaxLength("name", t.Name, MaxTaskNameLength, "task"); err != nil {
		return err
	}
	if err := validator.ValidateMaxLength("description", t.Desc, MaxTaskDescriptionLength, "task"); err != nil {
		return err
	}
	if err := validator.ValidateRequiredID(t.ProjectID, "task"); err != nil {
		return NewValidationError("project_id", "project ID cannot be empty")
	}
	if err := validator.ValidateEnum("priority", int(t.Priority), int(Low), int(High), "task"); err != nil {
		return err
	}
	if err := validator.ValidateEnum("status", int(t.Status), int(NotStarted), int(Done), "task"); err != nil {
		return err
	}
	if err := validator.ValidateEnum("type", int(t.Type), int(RegularTask), int(Feature), "task"); err != nil {
		return err
	}
	// Validate that a task cannot block itself
	if t.BlockedBy != nil && t.IntID != 0 && *t.BlockedBy == t.IntID {
		return NewValidationError("blocked_by", "task cannot block itself")
	}
	return nil
}

func (t Task) Title() string         { return t.Name }
func (t Task) Description() string   { return t.Desc }
func (t Task) GetPriority() Priority { return t.Priority }
func (t Task) FilterValue() string   { return t.Name }

// SortTasks sorts a slice of tasks based on the status
// For NotStarted: priority DESC, then created_at ASC (oldest highest priority first)
// For InProgress and Done: updated_at DESC (newest changes first)
func SortTasks(tasks []Task, status Status) []Task {
	sorted := make([]Task, len(tasks))
	copy(sorted, tasks)

	if status == NotStarted {
		// Sort by priority DESC, then created_at ASC
		sort.Slice(sorted, func(i, j int) bool {
			if sorted[i].Priority != sorted[j].Priority {
				return sorted[i].Priority > sorted[j].Priority // Higher priority first
			}
			return sorted[i].CreatedAt.Before(sorted[j].CreatedAt) // Older first
		})
	} else {
		// Sort by updated_at DESC
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].UpdatedAt.After(sorted[j].UpdatedAt) // Newest first
		})
	}

	return sorted
}
