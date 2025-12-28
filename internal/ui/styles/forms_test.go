package styles

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFormFieldStyles(t *testing.T) {
	styles := GetFormFieldStyles()

	// Test that all styles are not nil
	assert.NotNil(t, styles.Placeholder, "Placeholder style should not be nil")
	assert.NotNil(t, styles.Text, "Text style should not be nil")
	assert.NotNil(t, styles.Cursor, "Cursor style should not be nil")
	assert.NotNil(t, styles.Border, "Border style should not be nil")
	assert.NotNil(t, styles.ErrorBorder, "ErrorBorder style should not be nil")
	assert.NotNil(t, styles.ErrorText, "ErrorText style should not be nil")
}

func TestConfigureNameInput(t *testing.T) {
	fieldStyles := GetFormFieldStyles()
	input := ConfigureNameInput(fieldStyles)

	// Test input properties
	assert.Equal(t, "Name *", input.Placeholder, "Should have correct placeholder")
	assert.Equal(t, 50, input.CharLimit, "Should have correct character limit")
	assert.Equal(t, 40, input.Width, "Should have correct width")
	assert.False(t, input.Focused(), "Should not be focused by default")
}

func TestConfigureDescriptionInput(t *testing.T) {
	fieldStyles := GetFormFieldStyles()
	input := ConfigureDescriptionInput(fieldStyles)

	// Test input properties
	assert.Equal(t, "Description (optional)", input.Placeholder, "Should have correct placeholder")
	assert.Equal(t, 200, input.CharLimit, "Should have correct character limit")
	assert.Equal(t, 40, input.Width, "Should have correct width")
}

func TestConfigureTaskNameInput(t *testing.T) {
	fieldStyles := GetFormFieldStyles()
	input := ConfigureTaskNameInput(fieldStyles)

	// Test input properties
	assert.Equal(t, "Task name *", input.Placeholder, "Should have task-specific placeholder")
	assert.Equal(t, 50, input.CharLimit, "Should have correct character limit")
}

func TestConfigureProjectNameInput(t *testing.T) {
	fieldStyles := GetFormFieldStyles()
	input := ConfigureProjectNameInput(fieldStyles)

	// Test input properties
	assert.Equal(t, "Project name *", input.Placeholder, "Should have project-specific placeholder")
	assert.Equal(t, 50, input.CharLimit, "Should have correct character limit")
}

func TestInputConfigurationsHaveDifferentPlaceholders(t *testing.T) {
	fieldStyles := GetFormFieldStyles()

	taskName := ConfigureTaskNameInput(fieldStyles)
	projectName := ConfigureProjectNameInput(fieldStyles)
	genericName := ConfigureNameInput(fieldStyles)

	assert.NotEqual(t, taskName.Placeholder, projectName.Placeholder, "Task and project placeholders should differ")
	assert.NotEqual(t, genericName.Placeholder, taskName.Placeholder, "Generic and task placeholders should differ")
	assert.NotEqual(t, genericName.Placeholder, projectName.Placeholder, "Generic and project placeholders should differ")
}
