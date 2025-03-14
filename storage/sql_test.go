package storage

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// TestModel represents a model for testing purposes
type TestModel struct {
	ID            uint64 `gorm:"primaryKey"`
	Name          string
	RelatedModels []RelatedTestModel `gorm:"foreignKey:ModelID"`
}

// Related model for testing preload
type RelatedTestModel struct {
	ID        uint64 `gorm:"primaryKey"`
	Name      string
	ModelID   uint64
	TestModel TestModel `gorm:"foreignKey:ModelID"`
}

// setupMockDB creates a new mock database connection for testing
func setupMockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}
	dialector := mysql.New(mysql.Config{
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	return db, mock, nil
}

// TestCreate verifies the Create method functionality
func TestCreate(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)

	store := &sqlStorage{db: db}
	ctx := context.Background()

	model := &TestModel{Name: "test"}
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `test_models` (`name`) VALUES (?)")).
		WithArgs(model.Name).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = store.Create(ctx, model)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestGet verifies the Get method functionality
func TestGet(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)

	store := &sqlStorage{db: db}
	ctx := context.Background()

	model := &TestModel{ID: 1, Name: "test"}
	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(model.ID, model.Name)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `test_models` WHERE `test_models`.`id` = ? ORDER BY `test_models`.`id` LIMIT ?")).
		WithArgs(model.ID, 1).
		WillReturnRows(rows)

	var result TestModel
	err = store.Get(ctx, model.ID, &result)
	assert.NoError(t, err)
	assert.Equal(t, model.Name, result.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestGetBy verifies the GetBy method functionality
func TestGetBy(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)

	store := &sqlStorage{db: db}
	ctx := context.Background()

	model := &TestModel{ID: 1, Name: "test"}
	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(model.ID, model.Name)
	filter := map[string]interface{}{"name": "test"}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `test_models` WHERE `name` = ? ORDER BY `test_models`.`id` LIMIT ?")).
		WithArgs(filter["name"], 1).
		WillReturnRows(rows)

	var result TestModel
	err = store.GetBy(ctx, filter, &result)
	assert.NoError(t, err)
	assert.Equal(t, model.Name, result.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdate(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)

	store := &sqlStorage{db: db}
	ctx := context.Background()

	id := uint64(1)
	updateData := &TestModel{Name: "updated"}

	// 修正 SQL 匹配模式以适应实际的 GORM 查询
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `test_models` SET `name`=? WHERE id = ?")).
		WithArgs(updateData.Name, id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = store.Update(ctx, id, updateData)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateBy(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)

	store := &sqlStorage{db: db}
	ctx := context.Background()

	filter := map[string]interface{}{"name": "test"}
	updateData := &TestModel{Name: "updated"}

	// Add transaction expectations
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `test_models` SET `name`=? WHERE `name` = ?")).
		WithArgs(updateData.Name, filter["name"]).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = store.UpdateBy(ctx, filter, updateData)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDelete(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)

	store := &sqlStorage{db: db}
	ctx := context.Background()

	id := uint64(1)

	// Add transaction expectations and fix SQL query format
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `test_models` WHERE `test_models`.`id` = ?")).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = store.Delete(ctx, id, &TestModel{})
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteBy(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)

	store := &sqlStorage{db: db}
	ctx := context.Background()

	filter := map[string]interface{}{"name": "test"}
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `test_models` WHERE `name` = ?")).
		WithArgs(filter["name"]).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// 指定模型类型
	err = store.DeleteBy(ctx, filter, &TestModel{})
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestList verifies the List method functionality
func TestList(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)

	store := &sqlStorage{db: db}
	ctx := context.Background()

	// Count query
	countRows := sqlmock.NewRows([]string{"count(*)"}).AddRow(2)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `test_models`")).
		WillReturnRows(countRows)

	// Data query
	dataRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "test1").
		AddRow(2, "test2")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `test_models` LIMIT ?")).
		WithArgs(10). // 只需要 LIMIT 参数
		WillReturnRows(dataRows)

	query := Query{
		Page: 1,
		Size: 10,
	}

	var results []TestModel
	total, err := store.List(ctx, &query, &results, nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Equal(t, 2, len(results))
	if len(results) > 0 {
		assert.Equal(t, "test1", results[0].Name)
		assert.Equal(t, "test2", results[1].Name)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestListWithFilter verifies the List method with filter functionality
func TestListWithFilter(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)

	store := &sqlStorage{db: db}
	ctx := context.Background()

	// First expect base count query (without filter)
	countRows := sqlmock.NewRows([]string{"count(*)"}).AddRow(1)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `test_models`")).
		WillReturnRows(countRows)

	// Then expect the filtered data query
	dataRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "test1")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `test_models` WHERE `name` = ? ORDER BY id asc LIMIT ?")).
		WithArgs("test1", 10).
		WillReturnRows(dataRows)

	query := Query{
		Filter: map[string]interface{}{"name": "test1"},
		Page:   1,
		Size:   10,
		Sort:   map[string]string{"id": "asc"},
	}

	var results []TestModel
	total, err := store.List(ctx, &query, &results, nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	if len(results) > 0 {
		assert.Equal(t, "test1", results[0].Name)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestListWithPreload verifies the List method with preload functionality
func TestListWithPreload(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)

	store := &sqlStorage{db: db}
	ctx := context.Background()

	// Count query
	countRows := sqlmock.NewRows([]string{"count(*)"}).AddRow(2)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `test_models`")).
		WillReturnRows(countRows)

	// Main data query
	dataRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "test1").
		AddRow(2, "test2")
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `test_models`")).
		WillReturnRows(dataRows)

	// Preload query
	relatedRows := sqlmock.NewRows([]string{"id", "name", "model_id"}).
		AddRow(1, "related1", 1).
		AddRow(2, "related2", 2)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `related_test_models` WHERE `related_test_models`.`model_id` IN (?,?)")).
		WithArgs(1, 2).
		WillReturnRows(relatedRows)

	query := Query{
		Page:    1,
		Size:    10,
		Preload: "RelatedModels", // 修改为正确的关联名称
	}

	var results []TestModel
	total, err := store.List(ctx, &query, &results, nil)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, results, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}
