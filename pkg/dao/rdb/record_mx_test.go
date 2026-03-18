package rdb

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/cylonchau/hermes/pkg/model"
)

func TestRecordDAO_Mock_CreateMXRecord(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewRecordDAO(db)
	ctx := context.Background()

	baseRecord := &model.Record{ZoneID: 1, Name: "@", Type: "MX", IsActive: true}
	mxRecord := &model.MXRecord{Host: "mail.example.com", Priority: 10}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `record`")).
		WithArgs(baseRecord.ZoneID, baseRecord.Name, baseRecord.Type, baseRecord.TTL, baseRecord.Remark, baseRecord.Tags, baseRecord.Source, baseRecord.IsActive, baseRecord.ViewID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `record_mx`")).
		WithArgs(1, mxRecord.Host, mxRecord.Priority, mxRecord.Remark, mxRecord.Provider).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = dao.CreateMXRecord(ctx, baseRecord, mxRecord)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), mxRecord.RecordID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRecordDAO_Mock_GetMXRecordByID(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewRecordDAO(db)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "record_id", "host", "priority"}).
		AddRow(1, 1, "mail.example.com", 10)

	recordRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "@")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `record_mx` WHERE record_id = ? ORDER BY `record_mx`.`id` LIMIT ?")).
		WithArgs(uint(1), 1).
		WillReturnRows(rows)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `record` WHERE `record`.`id` = ?")).
		WithArgs(1).
		WillReturnRows(recordRows)

	res, err := dao.GetMXRecordByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, uint16(10), res.Priority)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRecordDAO_Mock_ListMXRecords(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewRecordDAO(db)
	ctx := context.Background()

	viewID := int64(1)

	aRecordRows := sqlmock.NewRows([]string{"id", "record_id", "ttl"}).AddRow(1, 10, 600)
	recordRows := sqlmock.NewRows([]string{"id", "name", "type", "zone_id", "view_id"}).AddRow(10, "www", "MX", 5, 1)
	viewRows := sqlmock.NewRows([]string{"id", "name", "category", "value"}).AddRow(1, "LOCAL", "acl", "127.0.0.1")
	zoneRows := sqlmock.NewRows([]string{"id", "name"}).AddRow(5, "test.com.")

	mock.ExpectQuery("^SELECT .*? FROM `record_mx` JOIN record ON record.id = record_mx.record_id WHERE record.view_id = \\?$").
		WithArgs(viewID).
		WillReturnRows(aRecordRows)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `record` WHERE `record`.`id` = ?")).
		WithArgs(int64(10)).WillReturnRows(recordRows)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `view` WHERE `view`.`id` = ?")).
		WithArgs(int64(1)).WillReturnRows(viewRows)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `zone` WHERE `zone`.`id` = ?")).
		WithArgs(int64(5)).WillReturnRows(zoneRows)

	res, err := dao.ListMXRecords(ctx, &viewID)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "www", res[0].Record.Name)
	assert.Equal(t, "LOCAL", res[0].Record.View.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}
