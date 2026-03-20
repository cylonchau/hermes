package rdb

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/cylonchau/hermes/pkg/model"
)

func TestRecordDAO_Mock_CreateAAAARecord(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewRecordDAO(db)
	ctx := context.Background()

	baseRecord := &model.Record{ZoneID: 1, Name: "aaaa", Type: "AAAA", IsActive: true}
	aaaaRecord := &model.AAAARecord{IP: []byte{0x20, 0x01, 0x0d, 0xb8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `record`")).
		WithArgs(baseRecord.ZoneID, baseRecord.Name, baseRecord.Type, baseRecord.TTL, baseRecord.Remark, baseRecord.Tags, baseRecord.Source, baseRecord.IsActive, baseRecord.ViewID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `record_aaaa`")).
		WithArgs(1, aaaaRecord.IP, aaaaRecord.Remark).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = dao.CreateAAAARecord(ctx, baseRecord, aaaaRecord)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), aaaaRecord.RecordID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRecordDAO_Mock_GetAAAARecordByID(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewRecordDAO(db)
	ctx := context.Background()

	ip := []byte{0x20, 0x01, 0x0d, 0xb8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	rows := sqlmock.NewRows([]string{"id", "record_id", "ip"}).
		AddRow(1, 1, ip)

	recordRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "aaaa")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `record_aaaa` WHERE record_id = ? ORDER BY `record_aaaa`.`id` LIMIT ?")).
		WithArgs(uint(1), 1).
		WillReturnRows(rows)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `record` WHERE `record`.`id` = ?")).
		WithArgs(1).
		WillReturnRows(recordRows)

	res, err := dao.GetAAAARecordByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, ip, res.IP)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRecordDAO_Mock_ListAAAARecords(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewRecordDAO(db)
	ctx := context.Background()

	viewID := int64(1)

	aRecordRows := sqlmock.NewRows([]string{"id", "record_id", "ttl"}).AddRow(1, 10, 600)
	recordRows := sqlmock.NewRows([]string{"id", "name", "type", "zone_id", "view_id"}).AddRow(10, "www", "AAAA", 5, 1)
	viewRows := sqlmock.NewRows([]string{"id", "name", "category", "value"}).AddRow(1, "LOCAL", "acl", "127.0.0.1")
	zoneRows := sqlmock.NewRows([]string{"id", "name"}).AddRow(5, "test.com.")

	mock.ExpectQuery("^SELECT .*? FROM `record_aaaa` JOIN record ON record.id = record_aaaa.record_id WHERE record.view_id = \\?$").
		WithArgs(viewID).
		WillReturnRows(aRecordRows)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `record` WHERE `record`.`id` = ?")).
		WithArgs(int64(10)).WillReturnRows(recordRows)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `view` WHERE `view`.`id` = ?")).
		WithArgs(int64(1)).WillReturnRows(viewRows)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `zone` WHERE `zone`.`id` = ?")).
		WithArgs(int64(5)).WillReturnRows(zoneRows)

	// 4. Execute
	res, err := dao.ListAAAARecords(ctx, &viewID)

	// 5. Verify
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "www", res[0].Record.Name)
	assert.Equal(t, "LOCAL", res[0].Record.View.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}
