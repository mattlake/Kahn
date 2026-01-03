package domain

import (
	"fmt"
	"time"
)

type Project struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Color       string    `json:"color"`
	Tasks       []Task    `json:"tasks"`
}

func NewProject(name, description, color string) *Project {
	now := time.Now()
	return &Project{
		ID:          generateProjectID(),
		Name:        name,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
		Color:       color,
		Tasks:       []Task{},
	}
}

func generateProjectID() string {
	return fmt.Sprintf("proj_%d", time.Now().UnixNano())
}

func (p *Project) AddTask(task Task) {
	task.ProjectID = p.ID
	p.Tasks = append(p.Tasks, task)
	p.UpdatedAt = time.Now()
}

func (p *Project) RemoveTask(taskID string) bool {
	for i, task := range p.Tasks {
		if task.ID == taskID {
			p.Tasks = append(p.Tasks[:i], p.Tasks[i+1:]...)
			p.UpdatedAt = time.Now()
			return true
		}
	}
	return false
}

func (p *Project) GetTasksByStatus(status Status) []Task {
	var tasks []Task
	for _, task := range p.Tasks {
		if task.Status == status {
			tasks = append(tasks, task)
		}
	}
	// Apply sorting based on status to match repository behavior
	return SortTasks(tasks, status)
}

func (p *Project) UpdateTaskStatus(taskID string, newStatus Status) bool {
	for i, task := range p.Tasks {
		if task.ID == taskID {
			p.Tasks[i].Status = newStatus
			p.Tasks[i].UpdatedAt = time.Now()
			p.UpdatedAt = time.Now()
			return true
		}
	}
	return false
}

func (p *Project) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("project name is required")
	}
	if len(p.Name) > 50 {
		return fmt.Errorf("project name too long (max 50 characters)")
	}
	if len(p.Description) > 200 {
		return fmt.Errorf("project description too long (max 200 characters)")
	}
	return nil
}
