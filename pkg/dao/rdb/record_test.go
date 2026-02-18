package rdb

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/cylonchau/hermes/pkg/model"
)

func TestRecordDAO_Mock_CreateRecord(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewRecordDAO(db)
	ctx := context.Background()

	record := &model.Record{
		ZoneID:   1,
		Name:     "www",
		Type:     "A",
		TTL:      600,
		IsActive: true,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `record`")).
		WithArgs(record.ZoneID, record.Name, record.Type, record.TTL, record.Remark, record.Tags, record.Source, record.IsActive).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = dao.CreateRecord(ctx, record)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRecordDAO_Mock_GetRecordByID(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewRecordDAO(db)
	ctx := context.Background()

	// Row for Record
	recordRows := sqlmock.NewRows([]string{"id", "zone_id", "name", "type", "is_active"}).
		AddRow(1, 1, "www", "A", true)

	// Row for Zone (Preload)
	zoneRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "example.com")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `record` WHERE id = ? AND is_active = ? ORDER BY `record`.`id` LIMIT ?")).
		WithArgs(1, true, 1).
		WillReturnRows(recordRows)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `zone` WHERE `zone`.`id` = ?")).
		WithArgs(1).
		WillReturnRows(zoneRows)

	record, err := dao.GetRecordByID(ctx, 1)
	assert.NoError(t, err)
	assert.NotNil(t, record)
	assert.Equal(t, "www", record.Name)
	assert.NotNil(t, record.Zone)
	assert.Equal(t, "example.com", record.Zone.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRecordDAO_Mock_GetRecordsByZone(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewRecordDAO(db)
	ctx := context.Background()

	recordRows := sqlmock.NewRows([]string{"id", "zone_id", "name"}).
		AddRow(1, 1, "www").
		AddRow(2, 1, "mail")

	zoneRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "example.com")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `record` WHERE zone_id = ? AND is_active = ?")).
		WithArgs(1, true).
		WillReturnRows(recordRows)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `zone` WHERE `zone`.`id` = ?")).
		WithArgs(1).
		WillReturnRows(zoneRows)

	records, err := dao.GetRecordsByZone(ctx, 1)
	assert.NoError(t, err)
	assert.Len(t, records, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRecordDAO_Mock_SoftDeleteRecord(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewRecordDAO(db)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `record` SET `is_active`=? WHERE id = ?")).
		WithArgs(false, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = dao.SoftDeleteRecord(ctx, 1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
