package styles

import (
	"testing"

	"github.com/charmbracelet/bubbles/list"
	"github.com/stretchr/testify/assert"
	"kahn/internal/domain"
)

func TestGetListTitleStyles(t *testing.T) {
	styles := GetListTitleStyles()

	// Test that all styles are not nil
	assert.NotNil(t, styles.NotStarted, "NotStarted style should not be nil")
	assert.NotNil(t, styles.InProgress, "InProgress style should not be nil")
	assert.NotNil(t, styles.Done, "Done style should not be nil")
}

func TestGetListTitleStyle(t *testing.T) {
	// Test each status returns correct style
	notStartedStyle := GetListTitleStyle(domain.NotStarted)
	inProgressStyle := GetListTitleStyle(domain.InProgress)
	doneStyle := GetListTitleStyle(domain.Done)

	assert.NotEmpty(t, notStartedStyle.Render("test"), "NotStarted style should render")
	assert.NotEmpty(t, inProgressStyle.Render("test"), "InProgress style should render")
	assert.NotEmpty(t, doneStyle.Render("test"), "Done style should render")

	// Test default case
	defaultStyle := GetListTitleStyle(domain.Status(999))
	assert.Equal(t, notStartedStyle.String(), defaultStyle.String(), "Default should be NotStarted style")
}

func TestApplyListTitleStyles(t *testing.T) {
	// Create test lists
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), 20, 10)
	taskLists := []list.Model{defaultList, defaultList, defaultList}

	// Apply styles
	ApplyListTitleStyles(taskLists)

	// Verify styles were applied (non-empty title styles)
	assert.NotEmpty(t, taskLists[domain.NotStarted].Styles.Title.Render("test"), "NotStarted title style should be applied")
	assert.NotEmpty(t, taskLists[domain.InProgress].Styles.Title.Render("test"), "InProgress title style should be applied")
	assert.NotEmpty(t, taskLists[domain.Done].Styles.Title.Render("test"), "Done title style should be applied")
}

func TestApplyListTitleStyles_InsufficientLists(t *testing.T) {
	// Test with insufficient lists (should not panic)
	taskLists := []list.Model{}

	// Should not panic
	assert.NotPanics(t, func() {
		ApplyListTitleStyles(taskLists)
	}, "Should not panic with insufficient lists")
}
