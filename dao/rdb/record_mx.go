package rdb

import (
	"context"

	"gorm.io/gorm"

	"github.com/cylonchau/hermes/pkg/model"
)

// CreateMXRecord 创建MX记录
func (dao *RecordDAO) CreateMXRecord(ctx context.Context, record *model.Record, mxRecord *model.MXRecord) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(record).Error; err != nil {
			return err
		}
		mxRecord.RecordID = record.ID
		return tx.Create(mxRecord).Error
	})
}

// GetMXRecords 获取MX记录
func (dao *RecordDAO) GetMXRecords(ctx context.Context, zoneName, recordName string) ([]*model.MXRecord, error) {
	var mxRecords []*model.MXRecord
	err := dao.db.WithContext(ctx).
		Joins("JOIN dns_records ON dns_records.id = dns_mx_records.record_id").
		Joins("JOIN dns_zones ON dns_zones.id = dns_records.zone_id").
		Where("dns_zones.name = ? AND dns_zones.is_active = ? AND dns_records.name = ? AND dns_records.is_active = ?",
			zoneName, true, recordName, true).
		Preload("Record").
		Order("dns_mx_records.priority ASC").
		Find(&mxRecords).Error
	return mxRecords, err
}

// GetMXRecordByID 根据记录ID获取MX记录
func (dao *RecordDAO) GetMXRecordByID(ctx context.Context, recordID uint) (*model.MXRecord, error) {
	var mxRecord model.MXRecord
	err := dao.db.WithContext(ctx).
		Where("record_id = ?", recordID).
		Preload("Record").
		First(&mxRecord).Error
	if err != nil {
		return nil, err
	}
	return &mxRecord, nil
}

// UpdateMXRecord 更新MX记录
func (dao *RecordDAO) UpdateMXRecord(ctx context.Context, record *model.Record, mxRecord *model.MXRecord) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(record).Error; err != nil {
			return err
		}
		return tx.Save(mxRecord).Error
	})
}

// DeleteMXRecord 删除MX记录
func (dao *RecordDAO) DeleteMXRecord(ctx context.Context, recordID uint) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("record_id = ?", recordID).Delete(&model.MXRecord{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Record{}, recordID).Error
	})
}
