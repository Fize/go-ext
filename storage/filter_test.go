package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateFilter(t *testing.T) {
	validFilter := map[string]any{
		"name": "test",
		"age":  30,
	}

	invalidFilter := map[string]any{
		"name; DROP TABLE users; --": "test",
	}

	err := ValidateFilter(validFilter)
	assert.NoError(t, err)

	err = ValidateFilter(invalidFilter)
	assert.Error(t, err)
	assert.Equal(t, "invalid column name: name; DROP TABLE users; --", err.Error())
}

func TestIsValidColumnName(t *testing.T) {
	validNames := []string{"name", "age", "created_at"}
	invalidNames := []string{"name;", "age--", "created_at'"}

	for _, name := range validNames {
		assert.True(t, isValidColumnName(name))
	}

	for _, name := range invalidNames {
		assert.False(t, isValidColumnName(name))
	}
}
