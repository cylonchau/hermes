package rdb

import (
	"context"

	"gorm.io/gorm"

	"github.com/cylonchau/hermes/pkg/model"
)

// CreateNSRecord 创建NS记录
func (dao *RecordDAO) CreateNSRecord(ctx context.Context, record *model.Record, nsRecord *model.NSRecord) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(record).Error; err != nil {
			return err
		}
		nsRecord.RecordID = record.ID
		return tx.Create(nsRecord).Error
	})
}

// GetNSRecords 获取NS记录
func (dao *RecordDAO) GetNSRecords(ctx context.Context, zoneName, recordName string) ([]*model.NSRecord, error) {
	var nsRecords []*model.NSRecord
	err := dao.db.WithContext(ctx).
		Joins("JOIN dns_records ON dns_records.id = dns_ns_records.record_id").
		Joins("JOIN dns_zones ON dns_zones.id = dns_records.zone_id").
		Where("dns_zones.name = ? AND dns_zones.is_active = ? AND dns_records.name = ? AND dns_records.is_active = ?",
			zoneName, true, recordName, true).
		Preload("Record").
		Find(&nsRecords).Error
	return nsRecords, err
}

// GetNSRecordByID 根据记录ID获取NS记录
func (dao *RecordDAO) GetNSRecordByID(ctx context.Context, recordID uint) (*model.NSRecord, error) {
	var nsRecord model.NSRecord
	err := dao.db.WithContext(ctx).
		Where("record_id = ?", recordID).
		Preload("Record").
		First(&nsRecord).Error
	if err != nil {
		return nil, err
	}
	return &nsRecord, nil
}

// UpdateNSRecord 更新NS记录
func (dao *RecordDAO) UpdateNSRecord(ctx context.Context, record *model.Record, nsRecord *model.NSRecord) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(record).Error; err != nil {
			return err
		}
		return tx.Save(nsRecord).Error
	})
}

// DeleteNSRecord 删除NS记录
func (dao *RecordDAO) DeleteNSRecord(ctx context.Context, recordID uint) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("record_id = ?", recordID).Delete(&model.NSRecord{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Record{}, recordID).Error
	})
}
