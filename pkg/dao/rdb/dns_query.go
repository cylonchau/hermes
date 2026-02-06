package rdb

import (
	"context"

	"github.com/cylonchau/hermes/pkg/model"
)

// ========== DNSQueryRepository 接口实现（CoreDNS专用，优化查询性能） ==========

// QueryARecords CoreDNS专用A记录查询（高性能版本）
func (dao *RecordDAO) QueryARecords(ctx context.Context, zoneName, recordName string) ([]*model.ARecord, error) {
	var aRecords []*model.ARecord
	err := dao.db.WithContext(ctx).
		Select("record_a.*, record.ttl").
		Joins("JOIN record ON record.id = record_a.record_id").
		Joins("JOIN zone ON zone.id = record.zone_id").
		Where("zone.name = ? AND record.name = ? AND zone.is_active = ? AND record.is_active = ?",
			zoneName, recordName, true, true).
		Find(&aRecords).Error
	return aRecords, err
}

// QueryAAAARecords CoreDNS专用AAAA记录查询
func (dao *RecordDAO) QueryAAAARecords(ctx context.Context, zoneName, recordName string) ([]*model.AAAARecord, error) {
	var aaaaRecords []*model.AAAARecord
	err := dao.db.WithContext(ctx).
		Select("record_aaaa.*, record.ttl").
		Joins("JOIN record ON record.id = record_aaaa.record_id").
		Joins("JOIN zone ON zone.id = record.zone_id").
		Where("zone.name = ? AND record.name = ? AND zone.is_active = ? AND record.is_active = ?",
			zoneName, recordName, true, true).
		Find(&aaaaRecords).Error
	return aaaaRecords, err
}

// QueryMXRecords CoreDNS专用MX记录查询
func (dao *RecordDAO) QueryMXRecords(ctx context.Context, zoneName, recordName string) ([]*model.MXRecord, error) {
	var mxRecords []*model.MXRecord
	err := dao.db.WithContext(ctx).
		Select("record_mx.*, record.ttl").
		Joins("JOIN record ON record.id = record_mx.record_id").
		Joins("JOIN zone ON zone.id = record.zone_id").
		Where("zone.name = ? AND record.name = ? AND zone.is_active = ? AND record.is_active = ?",
			zoneName, recordName, true, true).
		Order("record_mx.priority ASC").
		Find(&mxRecords).Error
	return mxRecords, err
}

// QueryTXTRecords CoreDNS专用TXT记录查询
func (dao *RecordDAO) QueryTXTRecords(ctx context.Context, zoneName, recordName string) ([]*model.TXTRecord, error) {
	var txtRecords []*model.TXTRecord
	err := dao.db.WithContext(ctx).
		Select("record_txt.*, record.ttl").
		Joins("JOIN record ON record.id = record_txt.record_id").
		Joins("JOIN zone ON zone.id = record.zone_id").
		Where("zone.name = ? AND record.name = ? AND zone.is_active = ? AND record.is_active = ?",
			zoneName, recordName, true, true).
		Find(&txtRecords).Error
	return txtRecords, err
}

// QuerySOARecord CoreDNS专用SOA记录查询
func (dao *RecordDAO) QuerySOARecord(ctx context.Context, zoneName string) (*model.SOARecord, error) {
	var soaRecord model.SOARecord
	err := dao.db.WithContext(ctx).
		Select("record_soa.*, record.ttl").
		Joins("JOIN record ON record.id = record_soa.record_id").
		Joins("JOIN zone ON zone.id = record.zone_id").
		Where("zone.name = ? AND record.name IN (?, '@') AND zone.is_active = ? AND record.is_active = ?",
			zoneName, zoneName, true, true).
		First(&soaRecord).Error
	if err != nil {
		return nil, err
	}
	return &soaRecord, nil
}

// QueryNSRecords CoreDNS专用NS记录查询
func (dao *RecordDAO) QueryNSRecords(ctx context.Context, zoneName, recordName string) ([]*model.NSRecord, error) {
	var nsRecords []*model.NSRecord
	err := dao.db.WithContext(ctx).
		Select("record_ns.*, record.ttl").
		Joins("JOIN record ON record.id = record_ns.record_id").
		Joins("JOIN zone ON zone.id = record.zone_id").
		Where("zone.name = ? AND record.name = ? AND zone.is_active = ? AND record.is_active = ?",
			zoneName, recordName, true, true).
		Find(&nsRecords).Error
	return nsRecords, err
}

// QueryCNAMERecords CoreDNS专用CNAME记录查询
func (dao *RecordDAO) QueryCNAMERecords(ctx context.Context, zoneName, recordName string) ([]*model.CNAMERecord, error) {
	var cnameRecords []*model.CNAMERecord
	err := dao.db.WithContext(ctx).
		Select("record_cname.*, record.ttl").
		Joins("JOIN record ON record.id = record_cname.record_id").
		Joins("JOIN zone ON zone.id = record.zone_id").
		Where("zone.name = ? AND record.name = ? AND zone.is_active = ? AND record.is_active = ?",
			zoneName, recordName, true, true).
		Find(&cnameRecords).Error
	return cnameRecords, err
}

// QuerySRVRecords CoreDNS专用SRV记录查询
func (dao *RecordDAO) QuerySRVRecords(ctx context.Context, zoneName, recordName string) ([]*model.SRVRecord, error) {
	var srvRecords []*model.SRVRecord
	err := dao.db.WithContext(ctx).
		Select("record_srv.*, record.ttl").
		Joins("JOIN record ON record.id = record_srv.record_id").
		Joins("JOIN zone ON zone.id = record.zone_id").
		Where("zone.name = ? AND record.name = ? AND zone.is_active = ? AND record.is_active = ?",
			zoneName, recordName, true, true).
		Order("record_srv.priority ASC, record_srv.weight DESC").
		Find(&srvRecords).Error
	return srvRecords, err
}
