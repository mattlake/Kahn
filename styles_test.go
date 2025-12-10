package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultStyle(t *testing.T) {
	style := defaultStyle

	assert.NotNil(t, style, "defaultStyle should not be nil")

	// Test that style has expected properties
	rendered := style.Render("test content")
	assert.NotEmpty(t, rendered, "Style should render content")
	assert.Contains(t, rendered, "test content", "Rendered content should contain original text")
}

func TestFocusedStyle(t *testing.T) {
	style := focusedStyle

	assert.NotNil(t, style, "focusedStyle should not be nil")

	// Test that style has expected properties
	rendered := style.Render("test content")
	assert.NotEmpty(t, rendered, "Style should render content")
	assert.Contains(t, rendered, "test content", "Rendered content should contain original text")
}

func TestStyles_Difference(t *testing.T) {
	defaultRendered := defaultStyle.Render("test")
	focusedRendered := focusedStyle.Render("test")

	// Styles should produce different output (focused has border)
	assert.NotEqual(t, defaultRendered, focusedRendered, "Default and focused styles should produce different output")

	// Both should contain the original content
	assert.Contains(t, defaultRendered, "test", "Default style should contain original content")
	assert.Contains(t, focusedRendered, "test", "Focused style should contain original content")
}

func TestStyles_BorderProperties(t *testing.T) {
	// Test that focused style has border properties
	focusedRendered := focusedStyle.Render("test")
	defaultRendered := defaultStyle.Render("test")

	// Focused style should be longer due to border characters
	assert.Greater(t, len(focusedRendered), len(defaultRendered), "Focused style should be longer due to border")
}

func TestStyles_MarginAndPadding(t *testing.T) {
	// Test that styles apply margin and padding
	shortContent := "x"
	defaultRendered := defaultStyle.Render(shortContent)

	// Should be longer than original content due to margin and padding
	assert.Greater(t, len(defaultRendered), len(shortContent), "Style should add margin and padding")
}

func TestStyles_EmptyContent(t *testing.T) {
	// Test rendering empty content
	defaultRendered := defaultStyle.Render("")
	focusedRendered := focusedStyle.Render("")

	// Should still produce some output due to margin, padding, and borders
	assert.NotEmpty(t, defaultRendered, "Default style should render empty content")
	assert.NotEmpty(t, focusedRendered, "Focused style should render empty content")
}

func TestStyles_NewlineHandling(t *testing.T) {
	content := "line1\nline2"

	defaultRendered := defaultStyle.Render(content)
	focusedRendered := focusedStyle.Render(content)

	// Should handle newlines properly
	assert.Contains(t, defaultRendered, "line1", "Default style should handle first line")
	assert.Contains(t, defaultRendered, "line2", "Default style should handle second line")
	assert.Contains(t, focusedRendered, "line1", "Focused style should handle first line")
	assert.Contains(t, focusedRendered, "line2", "Focused style should handle second line")
}
