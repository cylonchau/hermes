package rdb

import (
	"context"

	"gorm.io/gorm"

	"github.com/cylonchau/hermes/pkg/model"
)

// ========== SOARecordRepository 接口实现 ==========

// CreateSOARecord 创建SOA记录
func (dao *RecordDAO) CreateSOARecord(ctx context.Context, record *model.Record, soaRecord *model.SOARecord) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(record).Error; err != nil {
			return err
		}
		soaRecord.RecordID = record.ID
		return tx.Create(soaRecord).Error
	})
}

// GetSOARecord 获取SOA记录
func (dao *RecordDAO) GetSOARecord(ctx context.Context, zoneName string) (*model.SOARecord, error) {
	var soaRecord model.SOARecord
	err := dao.db.WithContext(ctx).
		Joins("JOIN dns_records ON dns_records.id = dns_soa_records.record_id").
		Joins("JOIN dns_zones ON dns_zones.id = dns_records.zone_id").
		Where("dns_zones.name = ? AND dns_zones.is_active = ? AND dns_records.name IN (?, '@') AND dns_records.is_active = ?",
			zoneName, true, zoneName, true).
		Preload("Record").
		First(&soaRecord).Error
	if err != nil {
		return nil, err
	}
	return &soaRecord, nil
}

// GetSOARecordByID 根据记录ID获取SOA记录
func (dao *RecordDAO) GetSOARecordByID(ctx context.Context, recordID uint) (*model.SOARecord, error) {
	var soaRecord model.SOARecord
	err := dao.db.WithContext(ctx).
		Where("record_id = ?", recordID).
		Preload("Record").
		First(&soaRecord).Error
	if err != nil {
		return nil, err
	}
	return &soaRecord, nil
}

// UpdateSOARecord 更新SOA记录
func (dao *RecordDAO) UpdateSOARecord(ctx context.Context, record *model.Record, soaRecord *model.SOARecord) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(record).Error; err != nil {
			return err
		}
		return tx.Save(soaRecord).Error
	})
}

// DeleteSOARecord 删除SOA记录
func (dao *RecordDAO) DeleteSOARecord(ctx context.Context, recordID uint) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("record_id = ?", recordID).Delete(&model.SOARecord{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Record{}, recordID).Error
	})
}
