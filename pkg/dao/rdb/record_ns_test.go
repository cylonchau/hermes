package rdb

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/cylonchau/hermes/pkg/model"
)

func TestRecordDAO_Mock_CreateNSRecord(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewRecordDAO(db)
	ctx := context.Background()

	baseRecord := &model.Record{ZoneID: 1, Name: "@", Type: "NS", IsActive: true}
	nsRecord := &model.NSRecord{NameServer: "ns1.example.com."}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `record`")).
		WithArgs(baseRecord.ZoneID, baseRecord.Name, baseRecord.Type, baseRecord.TTL, baseRecord.Remark, baseRecord.Tags, baseRecord.Source, baseRecord.IsActive).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `record_ns`")).
		WithArgs(1, nsRecord.NameServer, nsRecord.Remark, nsRecord.IsGlue).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = dao.CreateNSRecord(ctx, baseRecord, nsRecord)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), nsRecord.RecordID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRecordDAO_Mock_GetNSRecordByID(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewRecordDAO(db)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "record_id", "name_server"}).
		AddRow(1, 1, "ns1.example.com.")

	recordRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "@")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `record_ns` WHERE record_id = ? ORDER BY `record_ns`.`id` LIMIT ?")).
		WithArgs(uint(1), 1).
		WillReturnRows(rows)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `record` WHERE `record`.`id` = ?")).
		WithArgs(1).
		WillReturnRows(recordRows)

	res, err := dao.GetNSRecordByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, "ns1.example.com.", res.NameServer)
	assert.NoError(t, mock.ExpectationsWereMet())
}
