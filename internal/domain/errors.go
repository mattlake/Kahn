package domain

import "fmt"

// Error factory functions for consistent error creation across the application

// NewValidationError creates a ValidationError with the specified field and message
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{Field: field, Message: message}
}

// NewEmptyValidationError creates a ValidationError for empty field validation
func NewEmptyValidationError(field, entityName string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: fmt.Sprintf("%s %s cannot be empty", entityName, field),
	}
}

// NewLengthValidationError creates a ValidationError for field length validation
func NewLengthValidationError(field, entityName string, maxLen int) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: fmt.Sprintf("%s %s too long (max %d characters)", entityName, field, maxLen),
	}
}

// NewEnumValidationError creates a ValidationError for enum field validation
func NewEnumValidationError(field, entityName string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: fmt.Sprintf("invalid %s value for %s", field, entityName),
	}
}

// NewRepositoryError creates a RepositoryError for database operations
func NewRepositoryError(operation, entity, id string, cause error) *RepositoryError {
	return &RepositoryError{
		Operation: operation,
		Entity:    entity,
		ID:        id,
		Cause:     cause,
	}
}

// NewNotFoundError creates a RepositoryError for not found scenarios
func NewNotFoundError(entity, id string, cause error) *RepositoryError {
	return &RepositoryError{
		Operation: "get",
		Entity:    entity,
		ID:        id,
		Cause:     cause,
	}
}
