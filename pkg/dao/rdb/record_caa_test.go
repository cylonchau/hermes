package rdb

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/cylonchau/hermes/pkg/model"
)

func TestRecordDAO_Mock_CreateCAARecord(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewRecordDAO(db)
	ctx := context.Background()

	baseRecord := &model.Record{ZoneID: 1, Name: "@", Type: "CAA", IsActive: true}
	caaRecord := &model.CAARecord{Flag: 0, Tag: "issue", Value: "letsencrypt.org"}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `record`")).
		WithArgs(baseRecord.ZoneID, baseRecord.Name, baseRecord.Type, baseRecord.TTL, baseRecord.Remark, baseRecord.Tags, baseRecord.Source, baseRecord.IsActive).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `record_caa`")).
		WithArgs(1, caaRecord.Flag, caaRecord.Tag, caaRecord.Value).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = dao.CreateCAARecord(ctx, baseRecord, caaRecord)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), caaRecord.RecordID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRecordDAO_Mock_GetCAARecordByID(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewRecordDAO(db)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "record_id", "flag", "tag", "value"}).
		AddRow(1, 1, 0, "issue", "letsencrypt.org")

	recordRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "@")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `record_caa` WHERE record_id = ? ORDER BY `record_caa`.`id` LIMIT ?")).
		WithArgs(uint(1), 1).
		WillReturnRows(rows)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `record` WHERE `record`.`id` = ?")).
		WithArgs(1).
		WillReturnRows(recordRows)

	res, err := dao.GetCAARecordByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, uint8(0), res.Flag)
	assert.Equal(t, "issue", res.Tag)
	assert.NoError(t, mock.ExpectationsWereMet())
}
