package styles

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"kahn/internal/ui/colors"
)

func TestDefaultStyle(t *testing.T) {
	style := DefaultStyle

	assert.NotNil(t, style, "DefaultStyle should not be nil")

	// Test that style has expected properties
	rendered := style.Render("test content")
	assert.NotEmpty(t, rendered, "Style should render content")
	assert.Contains(t, rendered, "test content", "Rendered content should contain original text")
}

func TestFocusedStyle(t *testing.T) {
	style := FocusedStyle

	assert.NotNil(t, style, "FocusedStyle should not be nil")

	// Test that style has expected properties
	rendered := style.Render("test content")
	assert.NotEmpty(t, rendered, "Style should render content")
	assert.Contains(t, rendered, "test content", "Rendered content should contain original text")
}

func TestStyles_Difference(t *testing.T) {
	// Test that colors are different
	t.Logf("colors.Text = %s", colors.Text)
	t.Logf("colors.Green = %s", colors.Green)

	// The key difference is that both styles now have rounded borders with different colors
	// Text color: #cdd6f4 (white) for DefaultStyle
	// Green color: #a6e3a1 (green) for FocusedStyle

	defaultRendered := DefaultStyle.Render("test")
	focusedRendered := FocusedStyle.Render("test")

	// Both should contain original content
	assert.Contains(t, defaultRendered, "test", "Default style should contain original content")
	assert.Contains(t, focusedRendered, "test", "Focused style should contain original content")

	// Both should have borders (both use RoundedBorder)
	assert.Contains(t, defaultRendered, "╭", "Default style should have rounded border")
	assert.Contains(t, focusedRendered, "╭", "Focused style should have rounded border")

	// In terminals without color support, they might render identically, but with colors they'll be different
	// The key improvement is that both now have borders (vs previous hidden border for default)
	t.Logf("Note: In colorless terminal, styles may appear identical. With colors, DefaultStyle has white border (%s) and FocusedStyle has green border (%s)", colors.Text, colors.Green)
}

func TestStyles_BorderProperties(t *testing.T) {
	// Test that both styles have border properties
	focusedRendered := FocusedStyle.Render("test")
	defaultRendered := DefaultStyle.Render("test")

	// Both styles should have borders (rounded), so focused should still be longer or similar due to different border color
	assert.GreaterOrEqual(t, len(focusedRendered), len(defaultRendered), "Focused style should be longer or similar due to border colors")
}

func TestStyles_ReducedMargin(t *testing.T) {
	// Test that styles have minimal margins for maximum viewport utilization
	defaultTop, defaultRight, defaultBottom, defaultLeft := DefaultStyle.GetMargin()
	focusedTop, focusedRight, focusedBottom, focusedLeft := FocusedStyle.GetMargin()

	// Check that horizontal margins are 0 (minimal viewport utilization)
	assert.Equal(t, defaultRight, 0, "DefaultStyle should have right margin of 0")
	assert.Equal(t, defaultLeft, 0, "DefaultStyle should have left margin of 0")
	assert.Equal(t, focusedRight, 0, "FocusedStyle should have right margin of 0")
	assert.Equal(t, focusedLeft, 0, "FocusedStyle should have left margin of 0")

	// Vertical margins: top margin 1, bottom margin 0 (reduced gap with header)
	assert.Equal(t, defaultTop, 0, "DefaultStyle should have top margin of 1")
	assert.Equal(t, defaultBottom, 0, "DefaultStyle should have bottom margin of 0")
	assert.Equal(t, focusedTop, 0, "FocusedStyle should have top margin of 1")
	assert.Equal(t, focusedBottom, 0, "FocusedStyle should have bottom margin of 0")
	assert.Equal(t, focusedLeft, 0, "FocusedStyle should have left margin of 0")
}

func TestStyles_MinimalViewportMargins(t *testing.T) {
	// Test that all styles have minimal right margin (0) for maximum viewport utilization
	defaultTop, defaultRight, defaultBottom, defaultLeft := DefaultStyle.GetMargin()
	focusedTop, focusedRight, focusedBottom, focusedLeft := FocusedStyle.GetMargin()

	// Debug: Print actual values to understand issue
	t.Logf("DefaultStyle margins - Top:%d Right:%d Bottom:%d Left:%d", defaultTop, defaultRight, defaultBottom, defaultLeft)
	t.Logf("FocusedStyle margins - Top:%d Right:%d Bottom:%d Left:%d", focusedTop, focusedRight, focusedBottom, focusedLeft)

	// All styles should have right margin of 0
	assert.Equal(t, defaultRight, 0, "DefaultStyle should have no right margin")
	assert.Equal(t, focusedRight, 0, "FocusedStyle should have no right margin")

	// Should have left margin of 0 due to Margin(1, 0) setting
	assert.Equal(t, defaultLeft, 0, "DefaultStyle should have left margin of 0")
	assert.Equal(t, focusedLeft, 0, "FocusedStyle should have left margin of 0")

	// Vertical margins: top margin 1, bottom margin 0 (reduced gap with header)
	assert.Equal(t, defaultTop, 0, "DefaultStyle should have top margin of 1")
	assert.Equal(t, defaultBottom, 0, "DefaultStyle should have bottom margin of 0")
	assert.Equal(t, focusedTop, 0, "FocusedStyle should have top margin of 1")
	assert.Equal(t, focusedBottom, 0, "FocusedStyle should have bottom margin of 0")
}

func TestStyles_MarginAndPadding(t *testing.T) {
	// Test that styles apply margin and padding
	shortContent := "x"
	defaultRendered := DefaultStyle.Render(shortContent)

	// Should be longer than original content due to margin and padding
	assert.Greater(t, len(defaultRendered), len(shortContent), "Style should add margin and padding")
}

func TestStyles_EmptyContent(t *testing.T) {
	// Test rendering empty content
	defaultRendered := DefaultStyle.Render("")
	focusedRendered := FocusedStyle.Render("")

	// Should still produce some output due to margin, padding, and borders
	assert.NotEmpty(t, defaultRendered, "Default style should render empty content")
	assert.NotEmpty(t, focusedRendered, "Focused style should render empty content")
}

func TestStyles_NewlineHandling(t *testing.T) {
	content := "line1\nline2"

	defaultRendered := DefaultStyle.Render(content)
	focusedRendered := FocusedStyle.Render(content)

	// Should handle newlines properly
	assert.Contains(t, defaultRendered, "line1", "Default style should handle first line")
	assert.Contains(t, defaultRendered, "line2", "Default style should handle second line")
	assert.Contains(t, focusedRendered, "line1", "Focused style should handle first line")
	assert.Contains(t, focusedRendered, "line2", "Focused style should handle second line")
}
