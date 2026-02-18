package rdb

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestRecordDAO_Mock_QueryARecords(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewRecordDAO(db)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "record_id", "ip", "ttl"}).
		AddRow(1, 1, 16843009, 600)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `record_a`.*, `record`.ttl FROM `record_a` JOIN `record` ON `record`.id = `record_a`.record_id JOIN `zone` ON `zone`.id = `record`.zone_id WHERE `zone`.name = ? AND `record`.name = ? AND `zone`.is_active = 1 AND `record`.is_active = 1")).
		WithArgs("example.com", "www").
		WillReturnRows(rows)

	res, err := dao.QueryARecords(ctx, "example.com", "www")
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, uint32(600), res[0].TTL)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRecordDAO_Mock_QuerySOARecord(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewRecordDAO(db)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "record_id", "primary_ns", "ttl"}).
		AddRow(1, 1, "ns1.example.com.", 3600)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `record_soa`.*, `record`.ttl FROM `record_soa` JOIN `record` ON `record`.id = `record_soa`.record_id JOIN `zone` ON `zone`.id = `record`.zone_id WHERE `zone`.name = ? AND `record`.name IN (?, '@') AND `zone`.is_active = 1 AND `record`.is_active = 1 ORDER BY `record_soa`.id ASC")).
		WithArgs("example.com", "example.com").
		WillReturnRows(rows)

	res, err := dao.QuerySOARecord(ctx, "example.com")
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "ns1.example.com.", res.PrimaryNS)
	assert.NoError(t, mock.ExpectationsWereMet())
}
