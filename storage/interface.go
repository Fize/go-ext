package storage

import (
	"context"
)

// Query defines common query parameters
type Query struct {
	// support for filter condition, e.g: {"name": "test"}
	Filter map[string]any
	// support for pagination
	Page int
	// support for pagination size
	Size int
	// support sort condition, e.g: {"created_at": "desc"}
	Sort map[string]string
	// support for preloading
	Preload string
	// support for preload all associations
	AllPreload bool
	// support for association query
	AssociationKey string
}

// Storage defines basic database operations
type Storage interface {
	// returns the database client
	Client() any
	// Create creates a new record
	Create(ctx context.Context, model any) error
	// Get retrieves a single record by ID
	Get(ctx context.Context, id uint64, result any) error
	// GetBy retrieves a single record by custom conditions
	GetBy(ctx context.Context, filter map[string]any, result any) error
	// Update updates a record by ID
	Update(ctx context.Context, id uint64, data any) error
	// UpdateBy updates records that match the filter
	UpdateBy(ctx context.Context, filter map[string]any, data any) error
	// Delete deletes a record by ID
	Delete(ctx context.Context, id uint64, model any) error
	// DeleteBy deletes records that match the filter
	DeleteBy(ctx context.Context, filter map[string]any, model any) error
	// List retrieves multiple records with pagination, support association query
	List(ctx context.Context, query *Query, mainModel, assModel any) (total int64, err error)
}
