package styles

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDialogStyles(t *testing.T) {
	styles := GetDialogStyles()

	// Test that all styles are not nil
	assert.NotNil(t, styles.Title, "Title style should not be nil")
	assert.NotNil(t, styles.Message, "Message style should not be nil")
	assert.NotNil(t, styles.Instruction, "Instruction style should not be nil")
	assert.NotNil(t, styles.Form, "Form style should not be nil")
	assert.NotNil(t, styles.Error, "Error style should not be nil")

	// Test that styles render content
	assert.NotEmpty(t, styles.Title.Render("Test"), "Title style should render")
	assert.NotEmpty(t, styles.Message.Render("Test"), "Message style should render")
	assert.NotEmpty(t, styles.Instruction.Render("Test"), "Instruction style should render")
}

func TestGetProjectItemStyle(t *testing.T) {
	itemStyles := GetProjectItemStyle("#ff6b6b")

	// Test that both styles are not nil
	assert.NotNil(t, itemStyles.Normal, "Normal style should not be nil")
	assert.NotNil(t, itemStyles.Active, "Active style should not be nil")

	// Test that styles render content
	assert.NotEmpty(t, itemStyles.Normal.Render("Test"), "Normal style should render")
	assert.NotEmpty(t, itemStyles.Active.Render("Test"), "Active style should render")
}

func TestGetDeleteConfirmStyles(t *testing.T) {
	styles := GetDeleteConfirmStyles()

	// Test that all styles are not nil
	assert.NotNil(t, styles.Title, "Title style should not be nil")
	assert.NotNil(t, styles.Warning, "Warning style should not be nil")
	assert.NotNil(t, styles.Message, "Message style should not be nil")
	assert.NotNil(t, styles.Form, "Form style should not be nil")

	// Test that styles render content
	assert.NotEmpty(t, styles.Title.Render("Test"), "Title style should render")
	assert.NotEmpty(t, styles.Warning.Render("Test"), "Warning style should render")
	assert.NotEmpty(t, styles.Message.Render("Test"), "Message style should render")
}
