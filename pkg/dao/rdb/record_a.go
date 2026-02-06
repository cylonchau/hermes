package rdb

import (
	"context"

	"gorm.io/gorm"

	"github.com/cylonchau/hermes/pkg/model"
)

// ========== ARecordRepository 接口实现 ==========

// CreateARecord 创建A记录
func (dao *RecordDAO) CreateARecord(ctx context.Context, record *model.Record, ARecord *model.ARecord) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(record).Error; err != nil {
			return err
		}
		ARecord.RecordID = record.ID
		return tx.Create(ARecord).Error
	})
}

// GetARecords 获取A记录
func (dao *RecordDAO) GetARecords(ctx context.Context, zoneName, recordName string) ([]*model.ARecord, error) {
	var ARecords []*model.ARecord
	err := dao.db.WithContext(ctx).
		Joins("JOIN record ON record.id = record_a.record_id").
		Joins("JOIN zone ON zone.id = record.zone_id").
		Where("zone.name = ? AND zone.is_active = ? AND record.name = ? AND record.is_active = ?",
			zoneName, true, recordName, true).
		Preload("Record").
		Order("record_a.id ASC").
		Find(&ARecords).Error
	return ARecords, err
}

// GetARecordByID 根据记录ID获取A记录
func (dao *RecordDAO) GetARecordByID(ctx context.Context, recordID uint) (*model.ARecord, error) {
	var ARecord model.ARecord
	err := dao.db.WithContext(ctx).Where("record_id = ?", recordID).Preload("Record").First(&ARecord).Error
	if err != nil {
		return nil, err
	}
	return &ARecord, nil
}

// UpdateARecord 更新A记录
func (dao *RecordDAO) UpdateARecord(ctx context.Context, record *model.Record, ARecord *model.ARecord) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(record).Error; err != nil {
			return err
		}
		return tx.Save(ARecord).Error
	})
}

// DeleteARecord 删除A记录
func (dao *RecordDAO) DeleteARecord(ctx context.Context, recordID uint) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("record_id = ?", recordID).Delete(&model.ARecord{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Record{}, recordID).Error
	})
}
