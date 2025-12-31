package domain

import (
	"fmt"
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

func (t Task) Title() string         { return t.Name }
func (t Task) Description() string   { return t.Desc }
func (t Task) GetPriority() Priority { return t.Priority }
func (t Task) FilterValue() string   { return t.Name }
