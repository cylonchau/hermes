package rdb

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/cylonchau/hermes/pkg/model"
)

func TestRecordDAO_Mock_CreateSOARecord(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewRecordDAO(db)
	ctx := context.Background()

	baseRecord := &model.Record{ZoneID: 1, Name: "@", Type: "SOA", IsActive: true}
	soaRecord := &model.SOARecord{
		PrimaryNS: "ns1.example.com.",
		MBox:      "admin.example.com.",
		Serial:    2023010101,
		Refresh:   7200,
		Retry:     3600,
		Expire:    1209600,
		MinTTL:    3600,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `record`")).
		WithArgs(baseRecord.ZoneID, baseRecord.Name, baseRecord.Type, baseRecord.TTL, baseRecord.Remark, baseRecord.Tags, baseRecord.Source, baseRecord.IsActive).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `record_soa`")).
		WithArgs(1, soaRecord.PrimaryNS, soaRecord.MBox, soaRecord.Serial, soaRecord.Refresh, soaRecord.Retry, soaRecord.Expire, soaRecord.MinTTL, soaRecord.Remark).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = dao.CreateSOARecord(ctx, baseRecord, soaRecord)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), soaRecord.RecordID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRecordDAO_Mock_GetSOARecordByID(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewRecordDAO(db)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "record_id", "primary_ns", "mail_box"}).
		AddRow(1, 1, "ns1.example.com.", "admin.example.com.")

	recordRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "@")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `record_soa` WHERE record_id = ? ORDER BY `record_soa`.`id` LIMIT ?")).
		WithArgs(uint(1), 1).
		WillReturnRows(rows)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `record` WHERE `record`.`id` = ?")).
		WithArgs(1).
		WillReturnRows(recordRows)

	res, err := dao.GetSOARecordByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, "ns1.example.com.", res.PrimaryNS)
	assert.NoError(t, mock.ExpectationsWereMet())
}
