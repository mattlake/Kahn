package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchTasks_EmptyQuery_ReturnsAll(t *testing.T) {
	tasks := []Task{
		{Name: "Task 1"},
		{Name: "Task 2"},
		{Name: "Task 3"},
	}

	result := SearchTasks(tasks, "")

	assert.Equal(t, 3, len(result))
}

func TestSearchTasks_CaseInsensitive(t *testing.T) {
	tasks := []Task{
		{Name: "API Endpoint"},
		{Name: "Database Setup"},
	}

	result := SearchTasks(tasks, "api")

	assert.Equal(t, 1, len(result))
	assert.Equal(t, "API Endpoint", result[0].Name)
}

func TestSearchTasks_CaseInsensitive_UppercaseQuery(t *testing.T) {
	tasks := []Task{
		{Name: "api endpoint"},
		{Name: "database setup"},
	}

	result := SearchTasks(tasks, "API")

	assert.Equal(t, 1, len(result))
	assert.Equal(t, "api endpoint", result[0].Name)
}

func TestSearchTasks_SubstringMatch(t *testing.T) {
	tasks := []Task{
		{Name: "Frontend Development"},
		{Name: "Backend API"},
		{Name: "Database"},
	}

	result := SearchTasks(tasks, "end")

	assert.Equal(t, 2, len(result))
	// Should match "Frontend" and "Backend"
	names := []string{result[0].Name, result[1].Name}
	assert.Contains(t, names, "Frontend Development")
	assert.Contains(t, names, "Backend API")
}

func TestSearchTasks_NoMatches(t *testing.T) {
	tasks := []Task{
		{Name: "Task 1"},
		{Name: "Task 2"},
	}

	result := SearchTasks(tasks, "xyz")

	assert.Equal(t, 0, len(result))
}

func TestSearchTasks_MultipleMatches(t *testing.T) {
	tasks := []Task{
		{Name: "API Task 1"},
		{Name: "API Task 2"},
		{Name: "Database Task"},
		{Name: "API Task 3"},
	}

	result := SearchTasks(tasks, "api")

	assert.Equal(t, 3, len(result))
}

func TestSearchTasks_SpecialCharacters(t *testing.T) {
	tasks := []Task{
		{Name: "Task #1"},
		{Name: "Task @2"},
		{Name: "Task $3"},
	}

	result := SearchTasks(tasks, "#")

	assert.Equal(t, 1, len(result))
	assert.Equal(t, "Task #1", result[0].Name)
}

func TestSearchTasks_EmptyList(t *testing.T) {
	tasks := []Task{}

	result := SearchTasks(tasks, "test")

	assert.Equal(t, 0, len(result))
}

func TestCountSearchMatches_EmptyQuery(t *testing.T) {
	tasks := []Task{
		{Name: "Task 1"},
		{Name: "Task 2"},
		{Name: "Task 3"},
	}

	count := CountSearchMatches(tasks, "")

	assert.Equal(t, 3, count)
}

func TestCountSearchMatches_Accuracy(t *testing.T) {
	tasks := []Task{
		{Name: "API Task 1"},
		{Name: "API Task 2"},
		{Name: "Database Task"},
	}

	count := CountSearchMatches(tasks, "api")

	assert.Equal(t, 2, count)
}

func TestCountSearchMatches_NoMatches(t *testing.T) {
	tasks := []Task{
		{Name: "Task 1"},
		{Name: "Task 2"},
	}

	count := CountSearchMatches(tasks, "xyz")

	assert.Equal(t, 0, count)
}

func TestCountSearchMatches_CaseInsensitive(t *testing.T) {
	tasks := []Task{
		{Name: "API Endpoint"},
		{Name: "api gateway"},
		{Name: "Database"},
	}

	count := CountSearchMatches(tasks, "API")

	assert.Equal(t, 2, count)
}
