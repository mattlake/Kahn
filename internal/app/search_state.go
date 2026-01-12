package app

// SearchState manages the state of the search/filter feature including
// the active status, current query string, and match count.
type SearchState struct {
	active     bool
	query      string
	matchCount int
}

// NewSearchState creates a new SearchState with search mode inactive
func NewSearchState() *SearchState {
	return &SearchState{
		active:     false,
		query:      "",
		matchCount: 0,
	}
}

// IsActive returns whether search mode is currently active
func (ss *SearchState) IsActive() bool {
	return ss.active
}

// GetQuery returns the current search query string
func (ss *SearchState) GetQuery() string {
	return ss.query
}

// GetMatchCount returns the number of tasks matching the current query
func (ss *SearchState) GetMatchCount() int {
	return ss.matchCount
}

// Activate enters search mode and resets query and match count to empty/zero
func (ss *SearchState) Activate() {
	ss.active = true
	ss.query = ""
	ss.matchCount = 0
}

// SetQuery updates the search query string
func (ss *SearchState) SetQuery(query string) {
	ss.query = query
}

// UpdateMatchCount updates the count of tasks matching the current query
func (ss *SearchState) UpdateMatchCount(count int) {
	ss.matchCount = count
}

// Clear exits search mode and resets all state to defaults
func (ss *SearchState) Clear() {
	ss.active = false
	ss.query = ""
	ss.matchCount = 0
}

// AppendChar adds a single character to the end of the query string
func (ss *SearchState) AppendChar(char string) {
	ss.query += char
}

// Backspace removes the last character from the query string if it's not empty
func (ss *SearchState) Backspace() {
	if len(ss.query) > 0 {
		ss.query = ss.query[:len(ss.query)-1]
	}
}
