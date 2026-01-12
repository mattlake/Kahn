package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchState_DefaultInactive(t *testing.T) {
	ss := NewSearchState()
	assert.False(t, ss.IsActive())
	assert.Equal(t, "", ss.GetQuery())
	assert.Equal(t, 0, ss.GetMatchCount())
}

func TestSearchState_Activate(t *testing.T) {
	ss := NewSearchState()
	ss.Activate()
	assert.True(t, ss.IsActive())
	assert.Equal(t, "", ss.GetQuery())
}

func TestSearchState_Activate_ClearsExistingQuery(t *testing.T) {
	ss := NewSearchState()
	ss.SetQuery("old query")
	ss.UpdateMatchCount(5)

	ss.Activate()

	assert.True(t, ss.IsActive())
	assert.Equal(t, "", ss.GetQuery())
	assert.Equal(t, 0, ss.GetMatchCount())
}

func TestSearchState_SetQuery(t *testing.T) {
	ss := NewSearchState()
	ss.Activate()
	ss.SetQuery("test")

	assert.Equal(t, "test", ss.GetQuery())
}

func TestSearchState_AppendChar(t *testing.T) {
	ss := NewSearchState()
	ss.Activate()
	ss.AppendChar("a")
	ss.AppendChar("p")
	ss.AppendChar("i")

	assert.Equal(t, "api", ss.GetQuery())
}

func TestSearchState_Backspace(t *testing.T) {
	ss := NewSearchState()
	ss.Activate()
	ss.SetQuery("api")
	ss.Backspace()

	assert.Equal(t, "ap", ss.GetQuery())
}

func TestSearchState_Backspace_EmptyQuery(t *testing.T) {
	ss := NewSearchState()
	ss.Activate()
	ss.Backspace() // Should not panic

	assert.Equal(t, "", ss.GetQuery())
}

func TestSearchState_Backspace_MultipleChars(t *testing.T) {
	ss := NewSearchState()
	ss.Activate()
	ss.SetQuery("test")

	ss.Backspace()
	assert.Equal(t, "tes", ss.GetQuery())

	ss.Backspace()
	assert.Equal(t, "te", ss.GetQuery())

	ss.Backspace()
	assert.Equal(t, "t", ss.GetQuery())

	ss.Backspace()
	assert.Equal(t, "", ss.GetQuery())

	ss.Backspace() // Should not panic
	assert.Equal(t, "", ss.GetQuery())
}

func TestSearchState_Clear(t *testing.T) {
	ss := NewSearchState()
	ss.Activate()
	ss.SetQuery("test")
	ss.UpdateMatchCount(5)

	ss.Clear()

	assert.False(t, ss.IsActive())
	assert.Equal(t, "", ss.GetQuery())
	assert.Equal(t, 0, ss.GetMatchCount())
}

func TestSearchState_UpdateMatchCount(t *testing.T) {
	ss := NewSearchState()
	ss.Activate()
	ss.UpdateMatchCount(10)

	assert.Equal(t, 10, ss.GetMatchCount())
}

func TestSearchState_FullWorkflow(t *testing.T) {
	ss := NewSearchState()

	// Start inactive
	assert.False(t, ss.IsActive())

	// Activate search
	ss.Activate()
	assert.True(t, ss.IsActive())
	assert.Equal(t, "", ss.GetQuery())

	// Build query character by character
	ss.AppendChar("a")
	ss.AppendChar("p")
	ss.AppendChar("i")
	assert.Equal(t, "api", ss.GetQuery())

	// Update match count
	ss.UpdateMatchCount(3)
	assert.Equal(t, 3, ss.GetMatchCount())

	// Fix typo with backspace
	ss.Backspace()
	assert.Equal(t, "ap", ss.GetQuery())

	// Continue typing
	ss.AppendChar("p")
	assert.Equal(t, "app", ss.GetQuery())

	// Clear and exit
	ss.Clear()
	assert.False(t, ss.IsActive())
	assert.Equal(t, "", ss.GetQuery())
	assert.Equal(t, 0, ss.GetMatchCount())
}
