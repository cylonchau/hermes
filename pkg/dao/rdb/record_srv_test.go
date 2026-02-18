package rdb

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/cylonchau/hermes/pkg/model"
)

func TestRecordDAO_Mock_CreateSRVRecord(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewRecordDAO(db)
	ctx := context.Background()

	baseRecord := &model.Record{ZoneID: 1, Name: "_sip._tcp", Type: "SRV", IsActive: true}
	srvRecord := &model.SRVRecord{Priority: 10, Weight: 60, Port: 5060, Target: "sipserver.example.com."}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `record`")).
		WithArgs(baseRecord.ZoneID, baseRecord.Name, baseRecord.Type, baseRecord.TTL, baseRecord.Remark, baseRecord.Tags, baseRecord.Source, baseRecord.IsActive).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `record_srv`")).
		WithArgs(1, srvRecord.Priority, srvRecord.Weight, srvRecord.Port, srvRecord.Target, srvRecord.Remark, srvRecord.Service, srvRecord.Protocol).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = dao.CreateSRVRecord(ctx, baseRecord, srvRecord)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), srvRecord.RecordID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRecordDAO_Mock_GetSRVRecordByID(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewRecordDAO(db)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "record_id", "target", "port"}).
		AddRow(1, 1, "sipserver.example.com.", 5060)

	recordRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "_sip._tcp")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `record_srv` WHERE record_id = ? ORDER BY `record_srv`.`id` LIMIT ?")).
		WithArgs(uint(1), 1).
		WillReturnRows(rows)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `record` WHERE `record`.`id` = ?")).
		WithArgs(1).
		WillReturnRows(recordRows)

	res, err := dao.GetSRVRecordByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, uint16(5060), res.Port)
	assert.NoError(t, mock.ExpectationsWereMet())
}
