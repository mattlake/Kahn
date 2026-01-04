package domain

import "time"

func createTestTask(name, description, projectID string, status Status) *Task {
	now := time.Now()
	task := &Task{
		ID:        "test_task_" + now.Format("20060102150405.000000000"),
		Name:      name,
		Desc:      description,
		ProjectID: projectID,
		Status:    status,
		CreatedAt: now,
		UpdatedAt: now,
		Priority:  Medium,
	}
	return task
}

func createTestProject(name, description, color string) *Project {
	project := &Project{
		ID:          "test_proj_" + time.Now().Format("20060102150405.000000000"),
		Name:        name,
		Description: description,
		Color:       color,
		Tasks:       []Task{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return project
}

func createTestTaskWithPriority(name, description, projectID string, status Status, priority Priority, createdAt, updatedAt time.Time) *Task {
	task := &Task{
		ID:        "test_task_" + createdAt.Format("20060102150405.000000000"),
		Name:      name,
		Desc:      description,
		ProjectID: projectID,
		Status:    status,
		Type:      RegularTask,
		Priority:  priority,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
	return task
}

func createProjectWithTasksOfVaryingPriorities() *Project {
	project := createTestProject("Sorting Test Project", "Testing task sorting", "blue")

	// Create base time for consistent ordering
	baseTime := time.Now()

	// Create tasks with different priorities and timestamps
	highPriorityOld := createTestTaskWithPriority("High Priority Old", "High priority oldest", project.ID, NotStarted, High, baseTime, baseTime)
	highPriorityNew := createTestTaskWithPriority("High Priority New", "High priority newest", project.ID, NotStarted, High, baseTime.Add(1*time.Hour), baseTime.Add(1*time.Hour))

	mediumPriorityOld := createTestTaskWithPriority("Medium Priority Old", "Medium priority oldest", project.ID, NotStarted, Medium, baseTime.Add(30*time.Minute), baseTime.Add(30*time.Minute))

	lowPriorityOld := createTestTaskWithPriority("Low Priority Old", "Low priority oldest", project.ID, NotStarted, Low, baseTime.Add(2*time.Hour), baseTime.Add(2*time.Hour))
	lowPriorityNew := createTestTaskWithPriority("Low Priority New", "Low priority newest", project.ID, NotStarted, Low, baseTime.Add(3*time.Hour), baseTime.Add(3*time.Hour))

	// Add tasks in random order
	project.AddTask(*lowPriorityNew)
	project.AddTask(*highPriorityNew)
	project.AddTask(*mediumPriorityOld)
	project.AddTask(*lowPriorityOld)
	project.AddTask(*highPriorityOld)

	return project
}

type TaskUpdate struct {
	TaskID    string
	NewStatus Status
	Delay     time.Duration
}

func updateTaskStatusesWithTimestamps(project *Project, updates []TaskUpdate) {
	for _, update := range updates {
		if update.Delay > 0 {
			time.Sleep(update.Delay)
		}
		project.UpdateTaskStatus(update.TaskID, update.NewStatus)
	}
}

func getTaskNamesFromTasks(tasks []Task) []string {
	names := make([]string, len(tasks))
	for i, task := range tasks {
		names[i] = task.Name
	}
	return names
}
