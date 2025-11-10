package rdb

import (
	"context"

	"gorm.io/gorm"

	"github.com/cylonchau/hermes/pkg/model"
)

// CreateCNAMERecord 创建CNAME记录
func (dao *RecordDAO) CreateCNAMERecord(ctx context.Context, record *model.Record, cnameRecord *model.CNAMERecord) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(record).Error; err != nil {
			return err
		}
		cnameRecord.RecordID = record.ID
		return tx.Create(cnameRecord).Error
	})
}

// GetCNAMERecords 获取CNAME记录
func (dao *RecordDAO) GetCNAMERecords(ctx context.Context, zoneName, recordName string) ([]*model.CNAMERecord, error) {
	var cnameRecords []*model.CNAMERecord
	err := dao.db.WithContext(ctx).
		Joins("JOIN dns_records ON dns_records.id = dns_cname_records.record_id").
		Joins("JOIN dns_zones ON dns_zones.id = dns_records.zone_id").
		Where("dns_zones.name = ? AND dns_zones.is_active = ? AND dns_records.name = ? AND dns_records.is_active = ?",
			zoneName, true, recordName, true).
		Preload("Record").
		Find(&cnameRecords).Error
	return cnameRecords, err
}

// GetCNAMERecordByID 根据记录ID获取CNAME记录
func (dao *RecordDAO) GetCNAMERecordByID(ctx context.Context, recordID uint) (*model.CNAMERecord, error) {
	var cnameRecord model.CNAMERecord
	err := dao.db.WithContext(ctx).
		Where("record_id = ?", recordID).
		Preload("Record").
		First(&cnameRecord).Error
	if err != nil {
		return nil, err
	}
	return &cnameRecord, nil
}

// UpdateCNAMERecord 更新CNAME记录
func (dao *RecordDAO) UpdateCNAMERecord(ctx context.Context, record *model.Record, cnameRecord *model.CNAMERecord) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(record).Error; err != nil {
			return err
		}
		return tx.Save(cnameRecord).Error
	})
}

// DeleteCNAMERecord 删除CNAME记录
func (dao *RecordDAO) DeleteCNAMERecord(ctx context.Context, recordID uint) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("record_id = ?", recordID).Delete(&model.CNAMERecord{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Record{}, recordID).Error
	})
}
