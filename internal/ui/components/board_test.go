package components

import (
	"testing"

	"github.com/charmbracelet/bubbles/list"
	"github.com/stretchr/testify/assert"
	"kahn/internal/domain"
)

func TestBoardComponent_RenderProjectHeader(t *testing.T) {
	board := &BoardComponent{}

	// Test with valid project
	project := &domain.Project{
		ID:          "test_proj_1",
		Name:        "Test Project",
		Description: "A test project",
		Color:       "#ff6b6b",
	}

	result := board.RenderProjectHeader(project, 80)

	assert.NotEmpty(t, result, "RenderProjectHeader should not return empty string")
	assert.Contains(t, result, "Test Project", "Should contain project name")
	assert.Contains(t, result, "Project:", "Should contain project label")
}

func TestBoardComponent_RenderProjectHeader_NilProject(t *testing.T) {
	board := &BoardComponent{}

	result := board.RenderProjectHeader(nil, 80)

	assert.Empty(t, result, "RenderProjectHeader with nil project should return empty string")
}

func TestBoardComponent_RenderNoProjectsBoard(t *testing.T) {
	board := &BoardComponent{}

	result := board.RenderNoProjectsBoard(80, 24)

	assert.NotEmpty(t, result, "RenderNoProjectsBoard should not return empty string")
	assert.Contains(t, result, "No Projects", "Should contain 'No Projects' title")
	assert.Contains(t, result, "Create your first project", "Should contain instructions")
}

func TestBoardComponent_RenderTaskDeleteConfirm(t *testing.T) {
	board := &BoardComponent{}

	task := &domain.Task{
		ID:        "test_task_1",
		Name:      "Test Task",
		Desc:      "A test task",
		ProjectID: "test_proj_1",
		Status:    domain.NotStarted,
	}

	result := board.RenderTaskDeleteConfirm(task, 80, 24)

	assert.NotEmpty(t, result, "RenderTaskDeleteConfirm should not return empty string")
	assert.Contains(t, result, "Delete Task", "Should contain deletion title")
	assert.Contains(t, result, "Test Task", "Should contain task name")
	assert.Contains(t, result, "Yes, Delete", "Should contain confirmation option")
}

func TestBoardComponent_RenderTaskDeleteConfirm_NilTask(t *testing.T) {
	board := &BoardComponent{}

	result := board.RenderTaskDeleteConfirm(nil, 80, 24)

	assert.Empty(t, result, "RenderTaskDeleteConfirm with nil task should return empty string")
}

func TestBoardComponent_RenderBoard(t *testing.T) {
	board := &BoardComponent{}

	// Create test project
	project := &domain.Project{
		ID:          "test_proj_1",
		Name:        "Test Project",
		Description: "A test project",
		Color:       "#ff6b6b",
	}

	// Create test task lists
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), 20, 10)
	var taskLists [3]list.Model
	taskLists[domain.NotStarted] = defaultList
	taskLists[domain.InProgress] = defaultList
	taskLists[domain.Done] = defaultList

	result := board.RenderBoard(project, taskLists, domain.NotStarted, 80)

	assert.NotEmpty(t, result, "RenderBoard should not return empty string")
	assert.Contains(t, result, "Test Project", "Should contain project name")
}

func TestBoardComponent_RenderBoard_NilProject(t *testing.T) {
	board := &BoardComponent{}

	// Create test task lists
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), 20, 10)
	var taskLists [3]list.Model
	taskLists[domain.NotStarted] = defaultList
	taskLists[domain.InProgress] = defaultList
	taskLists[domain.Done] = defaultList

	result := board.RenderBoard(nil, taskLists, domain.NotStarted, 80)

	assert.Empty(t, result, "RenderBoard with nil project should return empty string")
}

func TestNewBoard(t *testing.T) {
	board := NewBoard()

	assert.NotNil(t, board, "NewBoard should return a non-nil Board")
	assert.NotNil(t, board.renderer, "Board should have a non-nil renderer")
}

func TestBoard_GetRenderer(t *testing.T) {
	board := NewBoard()

	renderer := board.GetRenderer()

	assert.NotNil(t, renderer, "GetRenderer should return a non-nil renderer")
	assert.Implements(t, (*BoardRenderer)(nil), renderer, "Renderer should implement BoardRenderer interface")
}
