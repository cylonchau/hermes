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
		Select("dns_a_records.*, dns_records.ttl").
		Joins("JOIN dns_records ON dns_records.id = dns_a_records.record_id").
		Joins("JOIN dns_zones ON dns_zones.id = dns_records.zone_id").
		Where("dns_zones.name = ? AND dns_records.name = ? AND dns_zones.is_active = ? AND dns_records.is_active = ?",
			zoneName, recordName, true, true).
		Find(&aRecords).Error
	return aRecords, err
}

// QueryAAAARecords CoreDNS专用AAAA记录查询
func (dao *RecordDAO) QueryAAAARecords(ctx context.Context, zoneName, recordName string) ([]*model.AAAARecord, error) {
	var aaaaRecords []*model.AAAARecord
	err := dao.db.WithContext(ctx).
		Select("dns_aaaa_records.*, dns_records.ttl").
		Joins("JOIN dns_records ON dns_records.id = dns_aaaa_records.record_id").
		Joins("JOIN dns_zones ON dns_zones.id = dns_records.zone_id").
		Where("dns_zones.name = ? AND dns_records.name = ? AND dns_zones.is_active = ? AND dns_records.is_active = ?",
			zoneName, recordName, true, true).
		Find(&aaaaRecords).Error
	return aaaaRecords, err
}

// QueryMXRecords CoreDNS专用MX记录查询
func (dao *RecordDAO) QueryMXRecords(ctx context.Context, zoneName, recordName string) ([]*model.MXRecord, error) {
	var mxRecords []*model.MXRecord
	err := dao.db.WithContext(ctx).
		Select("dns_mx_records.*, dns_records.ttl").
		Joins("JOIN dns_records ON dns_records.id = dns_mx_records.record_id").
		Joins("JOIN dns_zones ON dns_zones.id = dns_records.zone_id").
		Where("dns_zones.name = ? AND dns_records.name = ? AND dns_zones.is_active = ? AND dns_records.is_active = ?",
			zoneName, recordName, true, true).
		Order("dns_mx_records.priority ASC").
		Find(&mxRecords).Error
	return mxRecords, err
}

// QueryTXTRecords CoreDNS专用TXT记录查询
func (dao *RecordDAO) QueryTXTRecords(ctx context.Context, zoneName, recordName string) ([]*model.TXTRecord, error) {
	var txtRecords []*model.TXTRecord
	err := dao.db.WithContext(ctx).
		Select("dns_txt_records.*, dns_records.ttl").
		Joins("JOIN dns_records ON dns_records.id = dns_txt_records.record_id").
		Joins("JOIN dns_zones ON dns_zones.id = dns_records.zone_id").
		Where("dns_zones.name = ? AND dns_records.name = ? AND dns_zones.is_active = ? AND dns_records.is_active = ?",
			zoneName, recordName, true, true).
		Find(&txtRecords).Error
	return txtRecords, err
}

// QuerySOARecord CoreDNS专用SOA记录查询
func (dao *RecordDAO) QuerySOARecord(ctx context.Context, zoneName string) (*model.SOARecord, error) {
	var soaRecord model.SOARecord
	err := dao.db.WithContext(ctx).
		Select("dns_soa_records.*, dns_records.ttl").
		Joins("JOIN dns_records ON dns_records.id = dns_soa_records.record_id").
		Joins("JOIN dns_zones ON dns_zones.id = dns_records.zone_id").
		Where("dns_zones.name = ? AND dns_records.name IN (?, '@') AND dns_zones.is_active = ? AND dns_records.is_active = ?",
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
		Select("dns_ns_records.*, dns_records.ttl").
		Joins("JOIN dns_records ON dns_records.id = dns_ns_records.record_id").
		Joins("JOIN dns_zones ON dns_zones.id = dns_records.zone_id").
		Where("dns_zones.name = ? AND dns_records.name = ? AND dns_zones.is_active = ? AND dns_records.is_active = ?",
			zoneName, recordName, true, true).
		Find(&nsRecords).Error
	return nsRecords, err
}

// QueryCNAMERecords CoreDNS专用CNAME记录查询
func (dao *RecordDAO) QueryCNAMERecords(ctx context.Context, zoneName, recordName string) ([]*model.CNAMERecord, error) {
	var cnameRecords []*model.CNAMERecord
	err := dao.db.WithContext(ctx).
		Select("dns_cname_records.*, dns_records.ttl").
		Joins("JOIN dns_records ON dns_records.id = dns_cname_records.record_id").
		Joins("JOIN dns_zones ON dns_zones.id = dns_records.zone_id").
		Where("dns_zones.name = ? AND dns_records.name = ? AND dns_zones.is_active = ? AND dns_records.is_active = ?",
			zoneName, recordName, true, true).
		Find(&cnameRecords).Error
	return cnameRecords, err
}

// QuerySRVRecords CoreDNS专用SRV记录查询
func (dao *RecordDAO) QuerySRVRecords(ctx context.Context, zoneName, recordName string) ([]*model.SRVRecord, error) {
	var srvRecords []*model.SRVRecord
	err := dao.db.WithContext(ctx).
		Select("dns_srv_records.*, dns_records.ttl").
		Joins("JOIN dns_records ON dns_records.id = dns_srv_records.record_id").
		Joins("JOIN dns_zones ON dns_zones.id = dns_records.zone_id").
		Where("dns_zones.name = ? AND dns_records.name = ? AND dns_zones.is_active = ? AND dns_records.is_active = ?",
			zoneName, recordName, true, true).
		Order("dns_srv_records.priority ASC, dns_srv_records.weight DESC").
		Find(&srvRecords).Error
	return srvRecords, err
}
