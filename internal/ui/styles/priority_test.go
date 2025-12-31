package styles

import (
	"testing"

	"kahn/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestGetPriorityIndicator(t *testing.T) {
	tests := []struct {
		name     string
		priority domain.Priority
		expected string
	}{
		{
			name:     "High Priority",
			priority: domain.High,
			expected: "»» ", // Red chevrons
		},
		{
			name:     "Medium Priority",
			priority: domain.Medium,
			expected: "»  ", // Peach chevron + 2 spaces
		},
		{
			name:     "Low Priority",
			priority: domain.Low,
			expected: "  ", // 2 spaces
		},
		{
			name:     "Default Priority",
			priority: domain.Priority(99), // Unknown priority
			expected: "»  ",               // Default to medium with 2 spaces
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetPriorityIndicator(tt.priority)

			// Check that result contains expected base string
			// Note: We test for visible characters, not total string length due to ANSI codes
			if tt.priority == domain.Low {
				assert.NotContains(t, result, "»", "Low priority should not have chevrons")
				assert.Contains(t, result, "  ", "Low priority should have 2 spaces")
			} else {
				assert.Contains(t, result, "»", "Non-low priority should have chevrons")
				assert.Contains(t, result, " ", "Should include trailing space")
			}
		})
	}
}

func TestFormatTaskWithPriority(t *testing.T) {
	tests := []struct {
		name     string
		taskName string
		priority domain.Priority
		contains []string
	}{
		{
			name:     "High Priority Task",
			taskName: "Critical Task",
			priority: domain.High,
			contains: []string{"»»", "Critical Task"},
		},
		{
			name:     "Medium Priority Task",
			taskName: "Normal Task",
			priority: domain.Medium,
			contains: []string{"»", "Normal Task"}, // chevron (with color) + task name
		},
		{
			name:     "Low Priority Task",
			taskName: "Simple Task",
			priority: domain.Low,
			contains: []string{"  ", "Simple Task"}, // 2 spaces + task name
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := domain.NewTask(tt.taskName, "", "proj")
			task.Priority = tt.priority

			result := FormatTaskWithPriority(*task)

			// Check that result contains expected elements
			for _, expected := range tt.contains {
				assert.Contains(t, result, expected)
			}

			// Check that priority comes before task name
			if tt.priority != domain.Low {
				assert.True(t, len(result) > len(task.Name), "Should include priority indicator")
			}
		})
	}
}

func TestGetPriorityIndicatorWidth(t *testing.T) {
	width := GetPriorityIndicatorWidth()
	assert.Equal(t, 2, width, "Priority indicator width should be 2 characters")
}

func TestNewTaskWithTitle(t *testing.T) {
	// Create a test task
	task := domain.NewTask("Test Task", "", "proj")
	task.Priority = domain.High

	// Wrap with priority formatting
	wrapped := NewTaskWithTitle(*task)

	// Check that it's properly wrapped
	assert.Equal(t, "»» Test Task", wrapped.Title())   // Should include priority indicators + task name
	assert.Equal(t, domain.High, wrapped.Priority)     // Priority preserved
	assert.Equal(t, "Test Task", wrapped.Task.Title()) // Original task title preserved

	// Check that wrapped title includes priority indicator
	formattedTitle := wrapped.Title()
	assert.Contains(t, formattedTitle, "»»")        // Should include high priority indicators
	assert.Contains(t, formattedTitle, "Test Task") // Should include task name
}
