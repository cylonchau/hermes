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
		WithArgs(baseRecord.ZoneID, baseRecord.Name, baseRecord.Type, baseRecord.TTL, baseRecord.Remark, baseRecord.Tags, baseRecord.Source, baseRecord.IsActive, baseRecord.ViewID).
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

	recordRows := sqlmock.NewRows([]string{"id", "name", "view_id"}).
		AddRow(1, "a", 0)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `record_a` WHERE record_id = ? ORDER BY `record_a`.`id` LIMIT ?")).
		WithArgs(1, 1).
		WillReturnRows(aRows)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `record` WHERE `record`.`id` = ?")).
		WithArgs(1).
		WillReturnRows(recordRows)

	res, err := dao.GetARecordByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, uint32(16843009), res.IP)
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

func TestRecordDAO_Mock_ListARecords(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewRecordDAO(db)
	ctx := context.Background()

	viewID := int64(1)

	aRecordRows := sqlmock.NewRows([]string{"id", "record_id", "ip", "ttl"}).AddRow(1, 10, 16843009, 600)
	// record.id = 10, zone.id = 5, a_record.id = 1
	recordRows := sqlmock.NewRows([]string{"id", "name", "type", "zone_id", "view_id"}).AddRow(10, "www", "A", 5, 1)
	viewRows := sqlmock.NewRows([]string{"id", "name", "category", "value"}).AddRow(1, "LOCAL", "acl", "127.0.0.1")
	zoneRows := sqlmock.NewRows([]string{"id", "name"}).AddRow(5, "test.com.")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `record_a`.`id`,`record_a`.`record_id`,`record_a`.`ip`,`record_a`.`remark`,`record_a`.`ttl` FROM `record_a` JOIN record ON record.id = record_a.record_id WHERE record.view_id = ?")).
		WithArgs(viewID).
		WillReturnRows(aRecordRows)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `record` WHERE `record`.`id` = ?")).
		WithArgs(int64(10)).WillReturnRows(recordRows)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `view` WHERE `view`.`id` = ?")).
		WithArgs(int64(1)).WillReturnRows(viewRows)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `zone` WHERE `zone`.`id` = ?")).
		WithArgs(int64(5)).WillReturnRows(zoneRows)

	res, err := dao.ListARecords(ctx, &viewID)

	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "www", res[0].Record.Name)
	assert.Equal(t, "LOCAL", res[0].Record.View.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}
