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
		WithArgs(baseRecord.ZoneID, baseRecord.Name, baseRecord.Type, baseRecord.TTL, baseRecord.Remark, baseRecord.Tags, baseRecord.Source, baseRecord.IsActive).
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
