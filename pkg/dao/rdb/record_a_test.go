package rdb

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/cylonchau/hermes/pkg/model"
)

func TestRecordDAO_Mock_CreateARecord(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewRecordDAO(db)
	ctx := context.Background()

	baseRecord := &model.Record{ZoneID: 1, Name: "a", Type: "A", IsActive: true}
	aRecord := &model.ARecord{IP: 16843009} // 1.1.1.1

	mock.ExpectBegin()
	// Create base record
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `record`")).
		WithArgs(baseRecord.ZoneID, baseRecord.Name, baseRecord.Type, baseRecord.TTL, baseRecord.Remark, baseRecord.Tags, baseRecord.Source, baseRecord.IsActive).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Create A record (RecordID is populated from baseRecord.ID)
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `record_a`")).
		WithArgs(1, aRecord.IP, aRecord.Remark).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = dao.CreateARecord(ctx, baseRecord, aRecord)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), aRecord.RecordID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRecordDAO_Mock_GetARecordByID(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewRecordDAO(db)
	ctx := context.Background()

	aRows := sqlmock.NewRows([]string{"id", "record_id", "ip"}).
		AddRow(1, 1, 16843009)

	recordRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "a")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `record_a` WHERE record_id = ? ORDER BY `record_a`.`id` LIMIT ?")).
		WithArgs(1, 1).
		WillReturnRows(aRows)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `record` WHERE `record`.`id` = ?")).
		WithArgs(1).
		WillReturnRows(recordRows)

	res, err := dao.GetARecordByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, int32(16843009), res.IP)
	assert.Equal(t, "a", res.Record.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRecordDAO_Mock_DeleteARecord(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewRecordDAO(db)
	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `record_a` WHERE record_id = ?")).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `record` WHERE `record`.`id` = ?")).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = dao.DeleteARecord(ctx, 1)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
