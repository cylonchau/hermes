package rdb

import (
	"context"

	"gorm.io/gorm"

	"github.com/cylonchau/hermes/pkg/model"
)

// ========== SRVRecordRepository 接口实现 ==========

// CreateSRVRecord 创建SRV记录
func (dao *RecordDAO) CreateSRVRecord(ctx context.Context, record *model.Record, srvRecord *model.SRVRecord) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(record).Error; err != nil {
			return err
		}
		srvRecord.RecordID = record.ID
		return tx.Create(srvRecord).Error
	})
}

// GetSRVRecords 获取SRV记录
func (dao *RecordDAO) GetSRVRecords(ctx context.Context, zoneName, recordName string) ([]*model.SRVRecord, error) {
	var srvRecords []*model.SRVRecord
	err := dao.db.WithContext(ctx).
		Joins("JOIN record ON record.id = record_srv.record_id").
		Joins("JOIN zone ON zone.id = record.zone_id").
		Where("zone.name = ? AND zone.is_active = ? AND record.name = ? AND record.is_active = ?",
			zoneName, true, recordName, true).
		Preload("Record").
		Order("record_srv.priority ASC, record_srv.weight DESC").
		Find(&srvRecords).Error
	return srvRecords, err
}

// GetSRVRecordByID 根据记录ID获取SRV记录
func (dao *RecordDAO) GetSRVRecordByID(ctx context.Context, recordID uint) (*model.SRVRecord, error) {
	var srvRecord model.SRVRecord
	err := dao.db.WithContext(ctx).Where("record_id = ?", recordID).Preload("Record").First(&srvRecord).Error
	if err != nil {
		return nil, err
	}
	return &srvRecord, nil
}

// UpdateSRVRecord 更新SRV记录
func (dao *RecordDAO) UpdateSRVRecord(ctx context.Context, record *model.Record, srvRecord *model.SRVRecord) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(record).Error; err != nil {
			return err
		}
		return tx.Save(srvRecord).Error
	})
}

// DeleteSRVRecord 删除SRV记录
func (dao *RecordDAO) DeleteSRVRecord(ctx context.Context, recordID uint) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("record_id = ?", recordID).Delete(&model.SRVRecord{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Record{}, recordID).Error
	})
}
