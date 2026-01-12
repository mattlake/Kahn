package domain

import "strings"

// SearchTasks filters tasks whose Name contains the query using case-insensitive substring matching.
// Returns all tasks if query is empty.
func SearchTasks(tasks []Task, query string) []Task {
	if query == "" {
		return tasks
	}

	lowerQuery := strings.ToLower(query)
	filtered := make([]Task, 0, len(tasks))

	for _, task := range tasks {
		if strings.Contains(strings.ToLower(task.Name), lowerQuery) {
			filtered = append(filtered, task)
		}
	}

	return filtered
}

// CountSearchMatches returns the number of tasks whose Name contains the query.
// Returns total count if query is empty.
func CountSearchMatches(tasks []Task, query string) int {
	if query == "" {
		return len(tasks)
	}

	lowerQuery := strings.ToLower(query)
	count := 0

	for _, task := range tasks {
		if strings.Contains(strings.ToLower(task.Name), lowerQuery) {
			count++
		}
	}

	return count
}
