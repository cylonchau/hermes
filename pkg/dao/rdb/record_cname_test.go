package rdb

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/cylonchau/hermes/pkg/model"
)

func TestRecordDAO_Mock_CreateCNAMERecord(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewRecordDAO(db)
	ctx := context.Background()

	baseRecord := &model.Record{ZoneID: 1, Name: "cname", Type: "CNAME", IsActive: true}
	cnameRecord := &model.CNAMERecord{Target: "example.com."}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `record`")).
		WithArgs(baseRecord.ZoneID, baseRecord.Name, baseRecord.Type, baseRecord.TTL, baseRecord.Remark, baseRecord.Tags, baseRecord.Source, baseRecord.IsActive).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `record_cname`")).
		WithArgs(1, cnameRecord.Target, cnameRecord.Remark).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = dao.CreateCNAMERecord(ctx, baseRecord, cnameRecord)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), cnameRecord.RecordID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRecordDAO_Mock_GetCNAMERecordByID(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewRecordDAO(db)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "record_id", "target"}).
		AddRow(1, 1, "example.com.")

	recordRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "cname")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `record_cname` WHERE record_id = ? ORDER BY `record_cname`.`id` LIMIT ?")).
		WithArgs(uint(1), 1).
		WillReturnRows(rows)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `record` WHERE `record`.`id` = ?")).
		WithArgs(1).
		WillReturnRows(recordRows)

	res, err := dao.GetCNAMERecordByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, "example.com.", res.Target)
	assert.NoError(t, mock.ExpectationsWereMet())
}
