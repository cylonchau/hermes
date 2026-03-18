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

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `record_a`.*, `record`.ttl FROM `record_a` JOIN `record` ON `record`.id = `record_a`.record_id JOIN `zone` ON `zone`.id = `record`.zone_id WHERE (`zone`.name = ? AND `record`.name = ? AND `zone`.is_active = 1 AND `record`.is_active = 1) AND ((`record`.view_id IS NULL OR `record`.view_id = 0))")).
		WithArgs("example.com", "www").
		WillReturnRows(rows)

	res, err := dao.QueryARecords(ctx, "example.com", "www", 0)
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

	mock.ExpectQuery(regexp.QuoteMeta("SELECT `record_soa`.*, `record`.ttl FROM `record_soa` JOIN `record` ON `record`.id = `record_soa`.record_id JOIN `zone` ON `zone`.id = `record`.zone_id WHERE (`zone`.name = ? AND `record`.name IN (?, '@') AND `zone`.is_active = 1 AND `record`.is_active = 1) AND ((`record`.view_id IS NULL OR `record`.view_id = 0)) ORDER BY `record_soa`.id ASC")).
		WithArgs("example.com", "example.com").
		WillReturnRows(rows)

	res, err := dao.QuerySOARecord(ctx, "example.com", 0)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "ns1.example.com.", res.PrimaryNS)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRecordDAO_Mock_QueryARecords_Fallback(t *testing.T) {
	db, mock, err := setupMockDB()
	assert.NoError(t, err)
	dao := NewRecordDAO(db)
	ctx := context.Background()

	// 1. 模拟特定 View 无记录
	mock.ExpectQuery(regexp.QuoteMeta("SELECT `record_a`.*, `record`.ttl FROM `record_a` JOIN `record` ON `record`.id = `record_a`.record_id JOIN `zone` ON `zone`.id = `record`.zone_id WHERE (`zone`.name = ? AND `record`.name = ? AND `zone`.is_active = 1 AND `record`.is_active = 1) AND `record`.view_id = ?")).
		WithArgs("example.com", "www", int64(10)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "record_id", "ip", "ttl"})) // 返回空

	// 2. 模拟回退到默认视图成功
	rows := sqlmock.NewRows([]string{"id", "record_id", "ip", "ttl"}).
		AddRow(1, 1, 16843009, 600)
	mock.ExpectQuery(regexp.QuoteMeta("SELECT `record_a`.*, `record`.ttl FROM `record_a` JOIN `record` ON `record`.id = `record_a`.record_id JOIN `zone` ON `zone`.id = `record`.zone_id WHERE (`zone`.name = ? AND `record`.name = ? AND `zone`.is_active = 1 AND `record`.is_active = 1) AND ((`record`.view_id IS NULL OR `record`.view_id = 0))")).
		WithArgs("example.com", "www").
		WillReturnRows(rows)

	res, err := dao.QueryARecords(ctx, "example.com", "www", 10)
	assert.NoError(t, err)
	assert.Len(t, res, 1) // 应该能拿到默认视图的记录
	assert.NoError(t, mock.ExpectationsWereMet())
}
