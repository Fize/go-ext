package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/fize/go-ext/config"
	"github.com/fize/go-ext/log"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// sqlStorage represents the Storage implementation with GORM
type sqlStorage struct {
	db *gorm.DB
}

// NewSQLStorage creates a new Storage instance
func NewSQLStorage(cfg *config.SQLConfig) Storage {
	var db *gorm.DB
	var err error
	if cfg.Type == config.MySQL {
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.User, cfg.Password, cfg.Host, cfg.DB)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("failed to connect database with driver 'mysql': %v", err)
		}
	} else {
		db, err = gorm.Open(sqlite.Open(cfg.DB), &gorm.Config{})
		if err != nil {
			log.Fatalf("failed to connect database with driver 'sqlite': %v", err)
		}
	}
	return &sqlStorage{
		db: db,
	}
}

// Client implements Storage.Client
func (s *sqlStorage) Client() any {
	return s.db
}

// Create implements Storage.Create
func (s *sqlStorage) Create(ctx context.Context, model any) error {
	return s.db.WithContext(ctx).Create(model).Error
}

// Get implements Storage.Get
func (s *sqlStorage) Get(ctx context.Context, id uint64, result any) error {
	return s.db.WithContext(ctx).Model(result).First(result, id).Error
}

// GetBy implements Storage.GetBy
func (s *sqlStorage) GetBy(ctx context.Context, filter map[string]any, result any) error {
	if err := ValidateFilter(filter); err != nil {
		return err
	}
	return s.db.WithContext(ctx).Model(result).Where(filter).First(result).Error
}

// Update implements Storage.Update
func (s *sqlStorage) Update(ctx context.Context, id uint64, data any) error {
	return s.db.WithContext(ctx).Model(data).Where("id = ?", id).Save(data).Error
}

// UpdateBy implements Storage.UpdateBy
func (s *sqlStorage) UpdateBy(ctx context.Context, filter map[string]any, data any) error {
	if err := ValidateFilter(filter); err != nil {
		return err
	}
	result := s.db.WithContext(ctx).Model(data).Where(filter).Save(data)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Delete implements Storage.Delete
// If the record does not exist, it returns nil without error.
func (s *sqlStorage) Delete(ctx context.Context, id uint64, model any) error {
	return s.db.WithContext(ctx).Unscoped().Model(model).Delete("id = ?", id).Error
}

// DeleteBy implements Storage.DeleteBy
func (s *sqlStorage) DeleteBy(ctx context.Context, filter map[string]any, model any) error {
	if len(filter) == 0 {
		return errors.New("filter cannot be empty")
	}
	if err := ValidateFilter(filter); err != nil {
		return err
	}
	return s.db.WithContext(ctx).Unscoped().Model(model).Where(filter).Delete(filter).Error
}

// List implements Storage.List, support association query and preloading.
// If the association key is not empty, it will perform an association query.
// If the preload key is not empty, it will perform preloading.
// If both keys are not empty, it will return use association query.
func (s *sqlStorage) List(ctx context.Context, query *Query, mainModel, assModel any) (int64, error) {
	db := s.db.WithContext(ctx).Model(mainModel)
	// Count total records

	if len(query.Filter) > 0 {
		if err := ValidateFilter(query.Filter); err != nil {
			return 0, err
		}
		db = db.Where(query.Filter)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return 0, err
	}

	// Apply sorting
	if len(query.Sort) > 0 {
		for field, order := range query.Sort {
			db = db.Order(field + " " + order)
		}
	}

	// Apply pagination
	if query.Page > 0 && query.Size > 0 {
		offset := (query.Page - 1) * query.Size
		db = db.Offset(offset).Limit(query.Size)
	}

	//  assciation query or not
	if len(query.AssociationKey) > 0 {
		return total, db.Association(query.AssociationKey).Find(assModel)
	}
	if len(query.Preload) > 0 {
		return total, db.Preload(query.Preload).Find(mainModel).Error
	}
	if query.AllPreload {
		return total, db.Preload(clause.Associations).Find(mainModel).Error
	}

	return total, db.Find(mainModel).Error
}
