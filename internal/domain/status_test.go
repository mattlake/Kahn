package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatus_ToString(t *testing.T) {
	tests := []struct {
		name     string
		status   Status
		expected string
	}{
		{
			name:     "NotStarted status",
			status:   NotStarted,
			expected: "Not Started",
		},
		{
			name:     "InProgress status",
			status:   InProgress,
			expected: "In Progress",
		},
		{
			name:     "Done status",
			status:   Done,
			expected: "Done",
		},
		{
			name:     "Unknown status (should return Placeholder)",
			status:   Status(999),
			expected: "Placeholder",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.status.ToString()
			assert.Equal(t, tt.expected, result, "Status string representation should match expected value")
		})
	}
}

func TestStatus_Constants(t *testing.T) {
	// Test that status constants have expected integer values
	assert.Equal(t, Status(0), NotStarted, "NotStarted should be 0")
	assert.Equal(t, Status(1), InProgress, "InProgress should be 1")
	assert.Equal(t, Status(2), Done, "Done should be 2")
}

func TestStatus_Ordering(t *testing.T) {
	// Test that statuses are in the expected order for workflow progression
	assert.True(t, NotStarted < InProgress, "NotStarted should come before InProgress")
	assert.True(t, InProgress < Done, "InProgress should come before Done")
	assert.True(t, NotStarted < Done, "NotStarted should come before Done")
}
