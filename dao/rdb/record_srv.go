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
		Joins("JOIN dns_records ON dns_records.id = dns_srv_records.record_id").
		Joins("JOIN dns_zones ON dns_zones.id = dns_records.zone_id").
		Where("dns_zones.name = ? AND dns_zones.is_active = ? AND dns_records.name = ? AND dns_records.is_active = ?",
			zoneName, true, recordName, true).
		Preload("Record").
		Order("dns_srv_records.priority ASC, dns_srv_records.weight DESC").
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
