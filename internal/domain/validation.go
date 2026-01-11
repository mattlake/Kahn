package domain

import (
	"fmt"
	"strings"
)

// FieldValidator provides common validation utilities for domain entities
type FieldValidator struct{}

// NewFieldValidator creates a new FieldValidator instance
func NewFieldValidator() *FieldValidator {
	return &FieldValidator{}
}

// ValidateNotEmpty checks if a string field is not empty after trimming whitespace
func (v *FieldValidator) ValidateNotEmpty(field, value, entityName string) error {
	if strings.TrimSpace(value) == "" {
		return &ValidationError{
			Field:   field,
			Message: fmt.Sprintf("%s %s cannot be empty", entityName, field),
		}
	}
	return nil
}

// ValidateMaxLength checks if a string field doesn't exceed maximum length
func (v *FieldValidator) ValidateMaxLength(field, value string, maxLen int, entityName string) error {
	if len(value) > maxLen {
		return &ValidationError{
			Field:   field,
			Message: fmt.Sprintf("%s %s too long (max %d characters)", entityName, field, maxLen),
		}
	}
	return nil
}

// ValidateEnum checks if a value is within the valid enum range
func (v *FieldValidator) ValidateEnum(field string, value int, min, max int, entityName string) error {
	if value < min || value > max {
		return &ValidationError{
			Field:   field,
			Message: fmt.Sprintf("invalid %s value for %s", field, entityName),
		}
	}
	return nil
}

// ValidateRequiredID checks if an ID field is not empty
func (v *FieldValidator) ValidateRequiredID(id, entityName string) error {
	if id == "" {
		return &ValidationError{
			Field:   "id",
			Message: fmt.Sprintf("%s ID cannot be empty", entityName),
		}
	}
	return nil
}
