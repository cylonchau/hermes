package rdb

import (
	"context"

	"gorm.io/gorm"

	"github.com/cylonchau/hermes/pkg/model"
)

// CreateTXTRecord 创建TXT记录
func (dao *RecordDAO) CreateTXTRecord(ctx context.Context, record *model.Record, txtRecord *model.TXTRecord) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(record).Error; err != nil {
			return err
		}
		txtRecord.RecordID = record.ID
		return tx.Create(txtRecord).Error
	})
}

// GetTXTRecords 获取TXT记录
func (dao *RecordDAO) GetTXTRecords(ctx context.Context, zoneName, recordName string) ([]*model.TXTRecord, error) {
	var txtRecords []*model.TXTRecord
	err := dao.db.WithContext(ctx).
		Joins("JOIN record ON record.id = record_txt.record_id").
		Joins("JOIN zone ON zone.id = record.zone_id").
		Where("zone.name = ? AND zone.is_active = ? AND record.name = ? AND record.is_active = ?",
			zoneName, true, recordName, true).
		Preload("Record").
		Order("record_txt.id ASC").
		Find(&txtRecords).Error
	return txtRecords, err
}

// GetTXTRecordByID 根据记录ID获取TXT记录
func (dao *RecordDAO) GetTXTRecordByID(ctx context.Context, recordID uint) (*model.TXTRecord, error) {
	var txtRecord model.TXTRecord
	err := dao.db.WithContext(ctx).
		Where("record_id = ?", recordID).
		Preload("Record").
		First(&txtRecord).Error
	if err != nil {
		return nil, err
	}
	return &txtRecord, nil
}

// UpdateTXTRecord 更新TXT记录
func (dao *RecordDAO) UpdateTXTRecord(ctx context.Context, record *model.Record, txtRecord *model.TXTRecord) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(record).Error; err != nil {
			return err
		}
		return tx.Save(txtRecord).Error
	})
}

// DeleteTXTRecord 删除TXT记录
func (dao *RecordDAO) DeleteTXTRecord(ctx context.Context, recordID uint) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("record_id = ?", recordID).Delete(&model.TXTRecord{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Record{}, recordID).Error
	})
}
