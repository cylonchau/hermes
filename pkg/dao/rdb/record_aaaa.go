package rdb

import (
	"context"

	"gorm.io/gorm"

	"github.com/cylonchau/hermes/pkg/model"
)

// ========== AAAARecordRepository 接口实现 ==========

// CreateAAAARecord 创建AAAA记录
func (dao *RecordDAO) CreateAAAARecord(ctx context.Context, record *model.Record, AAAARecord *model.AAAARecord) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(record).Error; err != nil {
			return err
		}
		AAAARecord.RecordID = record.ID
		return tx.Create(AAAARecord).Error
	})
}

// GetARecords 获取AAAA记录
func (dao *RecordDAO) GetAAAARecords(ctx context.Context, zoneName, recordName string) ([]*model.AAAARecord, error) {
	var AAAARecord []*model.AAAARecord
	err := dao.db.WithContext(ctx).
		Joins("JOIN record ON record.id = record_aaaa.record_id").
		Joins("JOIN zone ON zone.id = record.zone_id").
		Where("zone.name = ? AND zone.is_active = ? AND record.name = ? AND record.is_active = ?",
			zoneName, true, recordName, true).
		Preload("Record").
		Order("record_aaaa.id ASC").
		Find(&AAAARecord).Error
	return AAAARecord, err
}

// GetAAAARecordByID 根据记录ID获取AAAA记录
func (dao *RecordDAO) GetAAAARecordByID(ctx context.Context, recordID uint) (*model.AAAARecord, error) {
	var AAAARecord model.AAAARecord
	err := dao.db.WithContext(ctx).Where("record_id = ?", recordID).Preload("Record").First(&AAAARecord).Error
	if err != nil {
		return nil, err
	}
	return &AAAARecord, nil
}

// UpdateAAAARecord 更新AAAA记录
func (dao *RecordDAO) UpdateAAAARecord(ctx context.Context, record *model.Record, AAAARecord *model.AAAARecord) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(record).Error; err != nil {
			return err
		}
		return tx.Save(AAAARecord).Error
	})
}

// DeleteAAAARecord 删除AAAA记录
func (dao *RecordDAO) DeleteAAAARecord(ctx context.Context, recordID uint) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("record_id = ?", recordID).Delete(&model.AAAARecord{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Record{}, recordID).Error
	})
}
