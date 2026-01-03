package domain

import (
	"fmt"
	"strings"
	"time"
)

type Task struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"project_id"`
	Name      string    `json:"name"`
	Desc      string    `json:"desc"`
	Status    Status    `json:"status"`
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

func NewTask(name, description, projectID string) *Task {
	now := time.Now()
	return &Task{
		ID:        generateTaskID(),
		ProjectID: projectID,
		Name:      name,
		Desc:      description,
		Status:    NotStarted,
		CreatedAt: now,
		UpdatedAt: now,
		Priority:  Low, // Changed default from Medium to Low
	}
}

func generateTaskID() string {
	return fmt.Sprintf("task_%d", time.Now().UnixNano())
}

func (t *Task) Validate() error {
	if strings.TrimSpace(t.Name) == "" {
		return &ValidationError{Field: "name", Message: "task name cannot be empty"}
	}
	if len(t.Name) > 100 {
		return &ValidationError{Field: "name", Message: "task name too long (max 100 characters)"}
	}
	if len(t.Desc) > 500 {
		return &ValidationError{Field: "description", Message: "task description too long (max 500 characters)"}
	}
	if t.ProjectID == "" {
		return &ValidationError{Field: "project_id", Message: "project ID cannot be empty"}
	}
	// Validate priority is within valid range
	if t.Priority < Low || t.Priority > High {
		return &ValidationError{Field: "priority", Message: "invalid priority value"}
	}
	// Validate status is within valid range
	if t.Status < NotStarted || t.Status > Done {
		return &ValidationError{Field: "status", Message: "invalid status value"}
	}
	return nil
}

func (t Task) Title() string         { return t.Name }
func (t Task) Description() string   { return t.Desc }
func (t Task) GetPriority() Priority { return t.Priority }
func (t Task) FilterValue() string   { return t.Name }
