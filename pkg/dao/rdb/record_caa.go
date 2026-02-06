package rdb

import (
	"context"

	"gorm.io/gorm"

	"github.com/cylonchau/hermes/pkg/model"
)

// CreateCAARecord 创建CAA记录
func (dao *RecordDAO) CreateCAARecord(ctx context.Context, record *model.Record, caaRecord *model.CAARecord) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(record).Error; err != nil {
			return err
		}
		caaRecord.RecordID = record.ID
		return tx.Create(caaRecord).Error
	})
}

// GetCAARecords 获取CAA记录
func (dao *RecordDAO) GetCAARecords(ctx context.Context, zoneName, recordName string) ([]*model.CAARecord, error) {
	var caaRecords []*model.CAARecord
	err := dao.db.WithContext(ctx).
		Joins("JOIN record ON record.id = record_caa.record_id").
		Joins("JOIN zone ON zone.id = record.zone_id").
		Where("zone.name = ? AND zone.is_active = ? AND record.name = ? AND record.is_active = ?",
			zoneName, true, recordName, true).
		Preload("Record").
		Order("record_caa.id ASC").
		Find(&caaRecords).Error
	return caaRecords, err
}

// GetCAARecordByID 根据记录ID获取CAA记录
func (dao *RecordDAO) GetCAARecordByID(ctx context.Context, recordID uint) (*model.CAARecord, error) {
	var caaRecord model.CAARecord
	err := dao.db.WithContext(ctx).
		Where("record_id = ?", recordID).
		Preload("Record").
		First(&caaRecord).Error
	if err != nil {
		return nil, err
	}
	return &caaRecord, nil
}

// UpdateCAARecord 更新CAA记录
func (dao *RecordDAO) UpdateCAARecord(ctx context.Context, record *model.Record, caaRecord *model.CAARecord) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(record).Error; err != nil {
			return err
		}
		return tx.Save(caaRecord).Error
	})
}

// DeleteCAARecord 删除CAA记录
func (dao *RecordDAO) DeleteCAARecord(ctx context.Context, recordID uint) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("record_id = ?", recordID).Delete(&model.CAARecord{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Record{}, recordID).Error
	})
}
