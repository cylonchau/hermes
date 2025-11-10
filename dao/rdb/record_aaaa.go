package rdb

import (
	"context"

	"gorm.io/gorm"

	"github.com/cylonchau/hermes/pkg/model"
)

// ========== ARecordRepository 接口实现 ==========

// CreateAAAARecord 创建A记录
func (dao *RecordDAO) CreateAAAARecord(ctx context.Context, record *model.Record, AAAARecord *model.AAAARecord) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(record).Error; err != nil {
			return err
		}
		AAAARecord.RecordID = record.ID
		return tx.Create(AAAARecord).Error
	})
}

// GetARecords 获取A记录
func (dao *RecordDAO) GetAAAARecords(ctx context.Context, zoneName, recordName string) ([]*model.AAAARecord, error) {
	var AAAARecord []*model.AAAARecord
	err := dao.db.WithContext(ctx).
		Joins("JOIN dns_records ON dns_records.id = dns_aaaa_records.record_id").
		Joins("JOIN dns_zones ON dns_zones.id = dns_records.zone_id").
		Where("dns_zones.name = ? AND dns_zones.is_active = ? AND dns_records.name = ? AND dns_records.is_active = ?",
			zoneName, true, recordName, true).
		Preload("Record").
		Order("dns_A_records.priority ASC, dns_A_records.weight DESC").
		Find(&AAAARecord).Error
	return AAAARecord, err
}

// GetARecordByID 根据记录ID获取A记录
func (dao *RecordDAO) GetAAAARecordByID(ctx context.Context, recordID uint) (*model.AAAARecord, error) {
	var AAAARecord model.AAAARecord
	err := dao.db.WithContext(ctx).Where("record_id = ?", recordID).Preload("Record").First(&AAAARecord).Error
	if err != nil {
		return nil, err
	}
	return &AAAARecord, nil
}

// UpdateARecord 更新A记录
func (dao *RecordDAO) UpdateAAAARecord(ctx context.Context, record *model.Record, AAAARecord *model.AAAARecord) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(record).Error; err != nil {
			return err
		}
		return tx.Save(AAAARecord).Error
	})
}

// DeleteARecord 删除A记录
func (dao *RecordDAO) DeleteAAAARecord(ctx context.Context, recordID uint) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("record_id = ?", recordID).Delete(&model.AAAARecord{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Record{}, recordID).Error
	})
}
