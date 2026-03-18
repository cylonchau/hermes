package rdb

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/cylonchau/hermes/pkg/model"
)

func TestViewDAO_Mock_Create(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewViewDAO(db)
	ctx := context.Background()

	view := &model.View{Name: "default", Category: "acl", Value: "1.1.1.1", Priority: 10}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `view`")).
		WithArgs(view.Name, view.Category, view.Value, view.Priority, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = dao.Create(ctx, view)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestViewDAO_Mock_GetByID(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewViewDAO(db)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "name", "category", "value", "priority", "created_at", "updated_at"}).
		AddRow(1, "default", "acl", "", 0, time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `view` WHERE `view`.`id` = ? ORDER BY `view`.`id` LIMIT ?")).
		WithArgs(1, 1).
		WillReturnRows(rows)

	res, err := dao.GetByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, "default", res.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestViewDAO_Mock_GetByName(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewViewDAO(db)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "name", "category", "value", "priority"}).
		AddRow(1, "default", "acl", "1.1.1.1", 0)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `view` WHERE name = ? ORDER BY `view`.`id` LIMIT ?")).
		WithArgs("default", 1).
		WillReturnRows(rows)

	res, err := dao.GetByName(ctx, "default")
	assert.NoError(t, err)
	assert.Equal(t, "default", res.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestViewDAO_Mock_GetAll(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewViewDAO(db)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "name", "category"}).
		AddRow(1, "default", "acl").
		AddRow(2, "another", "geoip")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `view`")).
		WillReturnRows(rows)

	res, err := dao.GetAll(ctx)
	assert.NoError(t, err)
	assert.Len(t, res, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestViewDAO_Mock_Update(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewViewDAO(db)
	ctx := context.Background()

	view := &model.View{ID: 1, Name: "default", Category: "geoip"}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `view`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = dao.Update(ctx, view)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestViewDAO_Mock_Delete(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewViewDAO(db)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `view` WHERE `view`.`id` = ?")).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = dao.Delete(ctx, 1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
