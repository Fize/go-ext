package storage

import (
	"fmt"
	"strings"
)

// ValidateFilter validates the filter parameters to prevent SQL injection
func ValidateFilter(filter map[string]any) error {
	for key := range filter {
		if !isValidColumnName(key) {
			return fmt.Errorf("invalid column name: %s", key)
		}
	}
	return nil
}

// isValidColumnName checks if the column name is valid to prevent SQL injection
func isValidColumnName(name string) bool {
	// TODO: Add more validation rules as needed
	return !strings.ContainsAny(name, "'\";--")
}
