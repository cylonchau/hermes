package rdb

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/cylonchau/hermes/pkg/model"
)

func TestRecordDAO_Mock_CreateTXTRecord(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewRecordDAO(db)
	ctx := context.Background()

	baseRecord := &model.Record{ZoneID: 1, Name: "v=spf1", Type: "TXT", IsActive: true}
	txtRecord := &model.TXTRecord{Text: "v=spf1 include:_spf.example.com ~all"}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `record`")).
		WithArgs(baseRecord.ZoneID, baseRecord.Name, baseRecord.Type, baseRecord.TTL, baseRecord.Remark, baseRecord.Tags, baseRecord.Source, baseRecord.IsActive).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `record_txt`")).
		WithArgs(1, txtRecord.Text, txtRecord.Remark, txtRecord.Purpose).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = dao.CreateTXTRecord(ctx, baseRecord, txtRecord)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), txtRecord.RecordID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRecordDAO_Mock_GetTXTRecordByID(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewRecordDAO(db)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "record_id", "text"}).
		AddRow(1, 1, "v=spf1 include:_spf.example.com ~all")

	recordRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "v=spf1")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `record_txt` WHERE record_id = ? ORDER BY `record_txt`.`id` LIMIT ?")).
		WithArgs(uint(1), 1).
		WillReturnRows(rows)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `record` WHERE `record`.`id` = ?")).
		WithArgs(1).
		WillReturnRows(recordRows)

	res, err := dao.GetTXTRecordByID(ctx, 1)
	assert.NoError(t, err)
	assert.Contains(t, res.Text, "include:_spf.example.com")
	assert.NoError(t, mock.ExpectationsWereMet())
}
