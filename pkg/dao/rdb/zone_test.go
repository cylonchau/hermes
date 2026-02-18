package rdb

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/cylonchau/hermes/pkg/model"
)

func TestZoneDAO_Mock_Create(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewZoneDAO(db)
	ctx := context.Background()

	zone := &model.Zone{
		Name:        "example.com",
		Description: "Testing zone",
		Contact:     "admin",
		IsActive:    true,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `zone`")).
		WithArgs(zone.Name, zone.Serial, zone.Description, zone.Remark, zone.Contact, zone.Email, zone.IsActive).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = dao.Create(ctx, zone)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestZoneDAO_Mock_GetByID(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewZoneDAO(db)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "name", "is_active"}).
		AddRow(1, "example.com", true)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `zone` WHERE id = ? AND is_active = ? ORDER BY `zone`.`id` LIMIT ?")).
		WithArgs(1, true, 1). // GORM First adds LIMIT 1
		WillReturnRows(rows)

	zone, err := dao.GetByID(ctx, 1)
	assert.NoError(t, err)
	assert.NotNil(t, zone)
	assert.Equal(t, "example.com", zone.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestZoneDAO_Mock_Update(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewZoneDAO(db)
	ctx := context.Background()

	zone := &model.Zone{
		ID:       1,
		Name:     "example.com",
		IsActive: true,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `zone` SET")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = dao.Update(ctx, zone)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestZoneDAO_Mock_Delete(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewZoneDAO(db)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `zone` WHERE `zone`.`id` = ?")).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = dao.Delete(ctx, 1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestZoneDAO_Mock_Search(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewZoneDAO(db)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "test1.com").
		AddRow(2, "test2.com")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `zone` WHERE is_active = ? AND (name LIKE ? OR description LIKE ? OR contact LIKE ?) LIMIT ?")).
		WithArgs(true, "%test%", "%test%", "%test%", 10).
		WillReturnRows(rows)

	zones, err := dao.Search(ctx, "test", 10, 0)
	assert.NoError(t, err)
	assert.Len(t, zones, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}
