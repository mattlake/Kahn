package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectDAO_Create(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	dao := NewProjectDAO(db)
	project := createTestProject("Test Project", "Test Description", "blue")

	err := dao.Create(project)
	assert.NoError(t, err, "Create should not return error")

	// Verify project was created
	count := countTableRows(t, db, "projects")
	assert.Equal(t, 1, count, "Should have 1 project in database")
}

func TestProjectDAO_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	dao := NewProjectDAO(db)
	project := createTestProject("Test Project", "Test Description", "blue")

	// Insert project
	err := dao.Create(project)
	require.NoError(t, err, "Create should not return error")

	// Get project by ID
	retrieved, err := dao.GetByID(project.ID)
	assert.NoError(t, err, "GetByID should not return error")
	assert.NotNil(t, retrieved, "Retrieved project should not be nil")
	assertProjectEqual(t, project, retrieved)
}

func TestProjectDAO_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	dao := NewProjectDAO(db)

	// Try to get non-existing project
	retrieved, err := dao.GetByID("non_existing_id")
	assert.Error(t, err, "GetByID should return error for non-existing project")
	assert.Nil(t, retrieved, "Retrieved project should be nil")
}

func TestProjectDAO_GetAll(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	dao := NewProjectDAO(db)

	// Create multiple projects
	project1 := createTestProject("Project 1", "Description 1", "blue")
	project2 := createTestProject("Project 2", "Description 2", "red")
	project3 := createTestProject("Project 3", "Description 3", "green")

	err := dao.Create(project1)
	require.NoError(t, err, "Create should not return error")
	err = dao.Create(project2)
	require.NoError(t, err, "Create should not return error")
	err = dao.Create(project3)
	require.NoError(t, err, "Create should not return error")

	// Get all projects
	projects, err := dao.GetAll()
	assert.NoError(t, err, "GetAll should not return error")
	assert.Len(t, projects, 3, "Should have 3 projects")

	// Verify project names
	names := []string{projects[0].Name, projects[1].Name, projects[2].Name}
	assert.Contains(t, names, "Project 1", "Should contain Project 1")
	assert.Contains(t, names, "Project 2", "Should contain Project 2")
	assert.Contains(t, names, "Project 3", "Should contain Project 3")
}

func TestProjectDAO_GetAll_Empty(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	dao := NewProjectDAO(db)

	// Get all projects from empty database
	projects, err := dao.GetAll()
	assert.NoError(t, err, "GetAll should not return error")
	assert.Len(t, projects, 0, "Should have no projects")
}

func TestProjectDAO_Update(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	dao := NewProjectDAO(db)
	project := createTestProject("Original Name", "Original Description", "blue")

	// Insert project
	err := dao.Create(project)
	require.NoError(t, err, "Create should not return error")

	// Update project
	project.Name = "Updated Name"
	project.Description = "Updated Description"
	project.Color = "red"

	err = dao.Update(project)
	assert.NoError(t, err, "Update should not return error")

	// Verify update
	retrieved, err := dao.GetByID(project.ID)
	assert.NoError(t, err, "GetByID should not return error")
	assert.Equal(t, "Updated Name", retrieved.Name, "Name should be updated")
	assert.Equal(t, "Updated Description", retrieved.Description, "Description should be updated")
	assert.Equal(t, "red", retrieved.Color, "Color should be updated")
}

func TestProjectDAO_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	dao := NewProjectDAO(db)
	project := createTestProject("Test Project", "Test Description", "blue")

	// Insert project
	err := dao.Create(project)
	require.NoError(t, err, "Create should not return error")

	// Verify project exists
	count := countTableRows(t, db, "projects")
	assert.Equal(t, 1, count, "Should have 1 project")

	// Delete project
	err = dao.Delete(project.ID)
	assert.NoError(t, err, "Delete should not return error")

	// Verify project is deleted
	count = countTableRows(t, db, "projects")
	assert.Equal(t, 0, count, "Should have no projects")

	// Verify project cannot be retrieved
	retrieved, err := dao.GetByID(project.ID)
	assert.Error(t, err, "GetByID should return error for deleted project")
	assert.Nil(t, retrieved, "Retrieved project should be nil")
}

func TestTaskDAO_Create(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	// Create project first
	project := insertTestProject(t, db, createTestProject("Test Project", "Test Description", "blue"))

	dao := NewTaskDAO(db)
	task := createTestTask("Test Task", "Test Description", project.ID, NotStarted)

	err := dao.Create(task)
	assert.NoError(t, err, "Create should not return error")

	// Verify task was created
	count := countTableRows(t, db, "tasks")
	assert.Equal(t, 1, count, "Should have 1 task in database")
}

func TestTaskDAO_GetByID(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	// Create project first
	project := insertTestProject(t, db, createTestProject("Test Project", "Test Description", "blue"))

	dao := NewTaskDAO(db)
	task := createTestTask("Test Task", "Test Description", project.ID, NotStarted)

	// Insert task
	err := dao.Create(task)
	require.NoError(t, err, "Create should not return error")

	// Get task by ID
	retrieved, err := dao.GetByID(task.ID)
	assert.NoError(t, err, "GetByID should not return error")
	assert.NotNil(t, retrieved, "Retrieved task should not be nil")
	assertTaskEqual(t, task, retrieved)
}

func TestTaskDAO_GetByProjectID(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	// Create projects
	project1 := insertTestProject(t, db, createTestProject("Project 1", "Description 1", "blue"))
	project2 := insertTestProject(t, db, createTestProject("Project 2", "Description 2", "red"))

	dao := NewTaskDAO(db)

	// Create tasks for project 1
	task1 := createTestTask("Task 1", "Description 1", project1.ID, NotStarted)
	task2 := createTestTask("Task 2", "Description 2", project1.ID, InProgress)

	// Create tasks for project 2
	task3 := createTestTask("Task 3", "Description 3", project2.ID, Done)

	err := dao.Create(task1)
	require.NoError(t, err, "Create should not return error")
	err = dao.Create(task2)
	require.NoError(t, err, "Create should not return error")
	err = dao.Create(task3)
	require.NoError(t, err, "Create should not return error")

	// Get tasks for project 1
	tasks1, err := dao.GetByProjectID(project1.ID)
	assert.NoError(t, err, "GetByProjectID should not return error")
	assert.Len(t, tasks1, 2, "Should have 2 tasks for project 1")

	// Get tasks for project 2
	tasks2, err := dao.GetByProjectID(project2.ID)
	assert.NoError(t, err, "GetByProjectID should not return error")
	assert.Len(t, tasks2, 1, "Should have 1 task for project 2")
}

func TestTaskDAO_GetByStatus(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	// Create project
	project := insertTestProject(t, db, createTestProject("Test Project", "Test Description", "blue"))

	dao := NewTaskDAO(db)

	// Create tasks with different statuses
	task1 := createTestTask("Task 1", "Description 1", project.ID, NotStarted)
	task2 := createTestTask("Task 2", "Description 2", project.ID, NotStarted)
	task3 := createTestTask("Task 3", "Description 3", project.ID, InProgress)
	task4 := createTestTask("Task 4", "Description 4", project.ID, Done)

	err := dao.Create(task1)
	require.NoError(t, err, "Create should not return error")
	err = dao.Create(task2)
	require.NoError(t, err, "Create should not return error")
	err = dao.Create(task3)
	require.NoError(t, err, "Create should not return error")
	err = dao.Create(task4)
	require.NoError(t, err, "Create should not return error")

	// Get tasks by status
	notStartedTasks, err := dao.GetByStatus(project.ID, NotStarted)
	assert.NoError(t, err, "GetByStatus should not return error")
	assert.Len(t, notStartedTasks, 2, "Should have 2 NotStarted tasks")

	inProgressTasks, err := dao.GetByStatus(project.ID, InProgress)
	assert.NoError(t, err, "GetByStatus should not return error")
	assert.Len(t, inProgressTasks, 1, "Should have 1 InProgress task")

	doneTasks, err := dao.GetByStatus(project.ID, Done)
	assert.NoError(t, err, "GetByStatus should not return error")
	assert.Len(t, doneTasks, 1, "Should have 1 Done task")
}

func TestTaskDAO_Update(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	// Create project
	project := insertTestProject(t, db, createTestProject("Test Project", "Test Description", "blue"))

	dao := NewTaskDAO(db)
	task := createTestTask("Original Task", "Original Description", project.ID, NotStarted)

	// Insert task
	err := dao.Create(task)
	require.NoError(t, err, "Create should not return error")

	// Update task
	task.Name = "Updated Task"
	task.Desc = "Updated Description"
	task.Status = InProgress
	task.Priority = High

	err = dao.Update(task)
	assert.NoError(t, err, "Update should not return error")

	// Verify update
	retrieved, err := dao.GetByID(task.ID)
	assert.NoError(t, err, "GetByID should not return error")
	assert.Equal(t, "Updated Task", retrieved.Name, "Name should be updated")
	assert.Equal(t, "Updated Description", retrieved.Desc, "Description should be updated")
	assert.Equal(t, InProgress, retrieved.Status, "Status should be updated")
	assert.Equal(t, High, retrieved.Priority, "Priority should be updated")
}

func TestTaskDAO_UpdateStatus(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	// Create project
	project := insertTestProject(t, db, createTestProject("Test Project", "Test Description", "blue"))

	dao := NewTaskDAO(db)
	task := createTestTask("Test Task", "Test Description", project.ID, NotStarted)

	// Insert task
	err := dao.Create(task)
	require.NoError(t, err, "Create should not return error")

	// Update status only
	err = dao.UpdateStatus(task.ID, InProgress)
	assert.NoError(t, err, "UpdateStatus should not return error")

	// Verify update
	retrieved, err := dao.GetByID(task.ID)
	assert.NoError(t, err, "GetByID should not return error")
	assert.Equal(t, InProgress, retrieved.Status, "Status should be updated")
	assert.Equal(t, "Test Task", retrieved.Name, "Name should remain unchanged")
	assert.Equal(t, "Test Description", retrieved.Desc, "Description should remain unchanged")
}

func TestTaskDAO_Delete(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	// Create project
	project := insertTestProject(t, db, createTestProject("Test Project", "Test Description", "blue"))

	dao := NewTaskDAO(db)
	task := createTestTask("Test Task", "Test Description", project.ID, NotStarted)

	// Insert task
	err := dao.Create(task)
	require.NoError(t, err, "Create should not return error")

	// Verify task exists
	count := countTableRows(t, db, "tasks")
	assert.Equal(t, 1, count, "Should have 1 task")

	// Delete task
	err = dao.Delete(task.ID)
	assert.NoError(t, err, "Delete should not return error")

	// Verify task is deleted
	count = countTableRows(t, db, "tasks")
	assert.Equal(t, 0, count, "Should have no tasks")

	// Verify task cannot be retrieved
	retrieved, err := dao.GetByID(task.ID)
	assert.Error(t, err, "GetByID should return error for deleted task")
	assert.Nil(t, retrieved, "Retrieved task should be nil")
}
