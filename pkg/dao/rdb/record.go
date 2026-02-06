package rdb

import (
	"context"

	"gorm.io/gorm"

	"github.com/cylonchau/hermes/pkg/model"
)

// RecordDAO Record的GORM数据访问层

type RecordDAO struct {
	db *gorm.DB
}

// NewRecordDAO 创建RecordDAO实例，实现所有Record相关接口
func NewRecordDAO(db *gorm.DB) *RecordDAO {
	return &RecordDAO{db: db}
}

// ========== RecordRepository 接口实现 ==========

// CreateRecord 创建基础记录
func (dao *RecordDAO) CreateRecord(ctx context.Context, record *model.Record) error {
	return dao.db.WithContext(ctx).Create(record).Error
}

// GetRecordByID 根据ID获取记录
func (dao *RecordDAO) GetRecordByID(ctx context.Context, recordID int64) (*model.Record, error) {
	var record model.Record
	err := dao.db.WithContext(ctx).
		Where("id = ? AND is_active = ?", recordID, true).
		Preload("Zone").
		First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// GetRecordsByZone 根据Zone获取所有记录
func (dao *RecordDAO) GetRecordsByZone(ctx context.Context, zoneID int64) ([]*model.Record, error) {
	var records []*model.Record
	err := dao.db.WithContext(ctx).
		Where("zone_id = ? AND is_active = ?", zoneID, true).
		Preload("Zone").
		Find(&records).Error
	return records, err
}

// GetRecordsByName 根据名称获取记录
func (dao *RecordDAO) GetRecordsByName(ctx context.Context, zoneID int64, recordName string) ([]*model.Record, error) {
	var records []*model.Record
	err := dao.db.WithContext(ctx).
		Where("zone_id = ? AND name = ? AND is_active = ?", zoneID, recordName, true).
		Find(&records).Error
	return records, err
}

// UpdateRecord 更新基础记录
func (dao *RecordDAO) UpdateRecord(ctx context.Context, record *model.Record) error {
	return dao.db.WithContext(ctx).Save(record).Error
}

// DeleteRecord 物理删除记录
func (dao *RecordDAO) DeleteRecord(ctx context.Context, recordID int64) error {
	return dao.db.WithContext(ctx).Delete(&model.Record{}, recordID).Error
}

// SoftDeleteRecord 软删除记录
func (dao *RecordDAO) SoftDeleteRecord(ctx context.Context, recordID int64) error {
	return dao.db.WithContext(ctx).Model(&model.Record{}).
		Where("id = ?", recordID).
		Update("is_active", false).Error
}

// CountRecordsByZone 统计Zone下的记录数量
func (dao *RecordDAO) CountRecordsByZone(ctx context.Context, zoneID int64) (int64, error) {
	var count int64
	err := dao.db.WithContext(ctx).Model(&model.Record{}).
		Where("zone_id = ? AND is_active = ?", zoneID, true).
		Count(&count).Error
	return count, err
}

// BatchCreateRecords 批量创建记录
func (dao *RecordDAO) BatchCreateRecords(ctx context.Context, records []*model.Record) error {
	return dao.db.WithContext(ctx).CreateInBatches(records, 100).Error
}

// BatchDeleteRecords 批量删除记录
func (dao *RecordDAO) BatchDeleteRecords(ctx context.Context, recordIDs []int64) error {
	return dao.db.WithContext(ctx).Delete(&model.Record{}, recordIDs).Error
}
