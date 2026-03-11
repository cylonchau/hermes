package rdb

import (
	"context"

	"github.com/cylonchau/hermes/pkg/model"
	"gorm.io/gorm"
)

// DNSQueryRepository 专用于 CoreDNS 的 DNS 查询接口
type DNSQueryRepository interface {
	QueryARecords(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.ARecord, error)
	QueryAAAARecords(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.AAAARecord, error)
	QueryMXRecords(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.MXRecord, error)
	QueryTXTRecords(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.TXTRecord, error)
	QuerySOARecord(ctx context.Context, zoneName string, viewID int64) (*model.SOARecord, error)
	QueryNSRecords(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.NSRecord, error)
	QueryCNAMERecords(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.CNAMERecord, error)
	QuerySRVRecords(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.SRVRecord, error)
}

// ========== RecordDAO 实现 DNSQueryRepository 接口 ==========

// QueryARecords CoreDNS专用A记录查询（高性能版本）
func (dao *RecordDAO) QueryARecords(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.ARecord, error) {
	var aRecords []*model.ARecord
	baseQuery := dao.db.WithContext(ctx).
		Model(&model.ARecord{}).
		Select("`record_a`.*, `record`.ttl").
		Joins("JOIN `record` ON `record`.id = `record_a`.record_id").
		Joins("JOIN `zone` ON `zone`.id = `record`.zone_id").
		Where("`zone`.name = ? AND `record`.name = ? AND `zone`.is_active = 1 AND `record`.is_active = 1",
			zoneName, recordName)

	if viewID > 0 {
		err := baseQuery.Session(&gorm.Session{}).Where("`record`.view_id = ?", viewID).Scan(&aRecords).Error
		if err != nil {
			return nil, err
		}
		if len(aRecords) > 0 {
			return aRecords, nil
		}
	}

	// 回退到默认视图
	err := baseQuery.Session(&gorm.Session{}).Where("(`record`.view_id IS NULL OR `record`.view_id = 0)").Scan(&aRecords).Error
	return aRecords, err
}

// QueryAAAARecords CoreDNS专用AAAA记录查询
func (dao *RecordDAO) QueryAAAARecords(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.AAAARecord, error) {
	var aaaaRecords []*model.AAAARecord
	baseQuery := dao.db.WithContext(ctx).
		Model(&model.AAAARecord{}).
		Select("`record_aaaa`.*, `record`.ttl").
		Joins("JOIN `record` ON `record`.id = `record_aaaa`.record_id").
		Joins("JOIN `zone` ON `zone`.id = `record`.zone_id").
		Where("`zone`.name = ? AND `record`.name = ? AND `zone`.is_active = 1 AND `record`.is_active = 1",
			zoneName, recordName)

	if viewID > 0 {
		err := baseQuery.Session(&gorm.Session{}).Where("`record`.view_id = ?", viewID).Scan(&aaaaRecords).Error
		if err != nil {
			return nil, err
		}
		if len(aaaaRecords) > 0 {
			return aaaaRecords, nil
		}
	}

	// 回退到默认视图
	err := baseQuery.Session(&gorm.Session{}).Where("(`record`.view_id IS NULL OR `record`.view_id = 0)").Scan(&aaaaRecords).Error
	return aaaaRecords, err
}

// QueryMXRecords CoreDNS专用MX记录查询
func (dao *RecordDAO) QueryMXRecords(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.MXRecord, error) {
	var mxRecords []*model.MXRecord
	baseQuery := dao.db.WithContext(ctx).
		Model(&model.MXRecord{}).
		Select("`record_mx`.*, `record`.ttl").
		Joins("JOIN `record` ON `record`.id = `record_mx`.record_id").
		Joins("JOIN `zone` ON `zone`.id = `record`.zone_id").
		Where("`zone`.name = ? AND `record`.name = ? AND `zone`.is_active = 1 AND `record`.is_active = 1",
			zoneName, recordName)

	if viewID > 0 {
		err := baseQuery.Session(&gorm.Session{}).Where("`record`.view_id = ?", viewID).Order("`record_mx`.priority ASC").Scan(&mxRecords).Error
		if err != nil {
			return nil, err
		}
		if len(mxRecords) > 0 {
			return mxRecords, nil
		}
	}

	// 回退到默认视图
	err := baseQuery.Session(&gorm.Session{}).Where("(`record`.view_id IS NULL OR `record`.view_id = 0)").Order("`record_mx`.priority ASC").Scan(&mxRecords).Error
	return mxRecords, err
}

// QueryTXTRecords CoreDNS专用TXT记录查询
func (dao *RecordDAO) QueryTXTRecords(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.TXTRecord, error) {
	var txtRecords []*model.TXTRecord
	baseQuery := dao.db.WithContext(ctx).
		Model(&model.TXTRecord{}).
		Select("`record_txt`.*, `record`.ttl").
		Joins("JOIN `record` ON `record`.id = `record_txt`.record_id").
		Joins("JOIN `zone` ON `zone`.id = `record`.zone_id").
		Where("`zone`.name = ? AND `record`.name = ? AND `zone`.is_active = 1 AND `record`.is_active = 1",
			zoneName, recordName)

	if viewID > 0 {
		err := baseQuery.Session(&gorm.Session{}).Where("`record`.view_id = ?", viewID).Scan(&txtRecords).Error
		if err != nil {
			return nil, err
		}
		if len(txtRecords) > 0 {
			return txtRecords, nil
		}
	}

	// 回退到默认视图
	err := baseQuery.Session(&gorm.Session{}).Where("(`record`.view_id IS NULL OR `record`.view_id = 0)").Scan(&txtRecords).Error
	return txtRecords, err
}

// QuerySOARecord CoreDNS专用SOA记录查询
func (dao *RecordDAO) QuerySOARecord(ctx context.Context, zoneName string, viewID int64) (*model.SOARecord, error) {
	var soaRecord model.SOARecord
	baseQuery := dao.db.WithContext(ctx).
		Model(&model.SOARecord{}).
		Select("`record_soa`.*, `record`.ttl").
		Joins("JOIN `record` ON `record`.id = `record_soa`.record_id").
		Joins("JOIN `zone` ON `zone`.id = `record`.zone_id").
		Where("`zone`.name = ? AND `record`.name IN (?, '@') AND `zone`.is_active = 1 AND `record`.is_active = 1",
			zoneName, zoneName)

	if viewID > 0 {
		err := baseQuery.Session(&gorm.Session{}).Where("`record`.view_id = ?", viewID).Order("`record_soa`.id ASC").Scan(&soaRecord).Error
		if err == nil {
			return &soaRecord, nil
		}
	}

	// 回退到默认视图
	err := baseQuery.Session(&gorm.Session{}).Where("(`record`.view_id IS NULL OR `record`.view_id = 0)").Order("`record_soa`.id ASC").Scan(&soaRecord).Error
	if err != nil {
		return nil, err
	}
	return &soaRecord, nil
}

// QueryNSRecords CoreDNS专用NS记录查询
func (dao *RecordDAO) QueryNSRecords(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.NSRecord, error) {
	var nsRecords []*model.NSRecord
	baseQuery := dao.db.WithContext(ctx).
		Model(&model.NSRecord{}).
		Select("`record_ns`.*, `record`.ttl").
		Joins("JOIN `record` ON `record`.id = `record_ns`.record_id").
		Joins("JOIN `zone` ON `zone`.id = `record`.zone_id").
		Where("`zone`.name = ? AND `record`.name = ? AND `zone`.is_active = 1 AND `record`.is_active = 1",
			zoneName, recordName)

	if viewID > 0 {
		err := baseQuery.Session(&gorm.Session{}).Where("`record`.view_id = ?", viewID).Scan(&nsRecords).Error
		if err != nil {
			return nil, err
		}
		if len(nsRecords) > 0 {
			return nsRecords, nil
		}
	}

	// 回退到默认视图
	err := baseQuery.Session(&gorm.Session{}).Where("(`record`.view_id IS NULL OR `record`.view_id = 0)").Scan(&nsRecords).Error
	return nsRecords, err
}

// QueryCNAMERecords CoreDNS专用CNAME记录查询
func (dao *RecordDAO) QueryCNAMERecords(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.CNAMERecord, error) {
	var cnameRecords []*model.CNAMERecord
	baseQuery := dao.db.WithContext(ctx).
		Model(&model.CNAMERecord{}).
		Select("`record_cname`.*, `record`.ttl").
		Joins("JOIN `record` ON `record`.id = `record_cname`.record_id").
		Joins("JOIN `zone` ON `zone`.id = `record`.zone_id").
		Where("`zone`.name = ? AND `record`.name = ? AND `zone`.is_active = 1 AND `record`.is_active = 1",
			zoneName, recordName)

	if viewID > 0 {
		err := baseQuery.Session(&gorm.Session{}).Where("`record`.view_id = ?", viewID).Scan(&cnameRecords).Error
		if err != nil {
			return nil, err
		}
		if len(cnameRecords) > 0 {
			return cnameRecords, nil
		}
	}

	// 回退到默认视图
	err := baseQuery.Session(&gorm.Session{}).Where("(`record`.view_id IS NULL OR `record`.view_id = 0)").Scan(&cnameRecords).Error
	return cnameRecords, err
}

// QuerySRVRecords CoreDNS专用SRV记录查询
func (dao *RecordDAO) QuerySRVRecords(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.SRVRecord, error) {
	var srvRecords []*model.SRVRecord
	baseQuery := dao.db.WithContext(ctx).
		Model(&model.SRVRecord{}).
		Select("`record_srv`.*, `record`.ttl").
		Joins("JOIN `record` ON `record`.id = `record_srv`.record_id").
		Joins("JOIN `zone` ON `zone`.id = `record`.zone_id").
		Where("`zone`.name = ? AND `record`.name = ? AND `zone`.is_active = 1 AND `record`.is_active = 1",
			zoneName, recordName)

	if viewID > 0 {
		err := baseQuery.Session(&gorm.Session{}).Where("`record`.view_id = ?", viewID).Order("`record_srv`.priority ASC, `record_srv`.weight DESC").Scan(&srvRecords).Error
		if err != nil {
			return nil, err
		}
		if len(srvRecords) > 0 {
			return srvRecords, nil
		}
	}

	// 回退到默认视图
	err := baseQuery.Session(&gorm.Session{}).Where("(`record`.view_id IS NULL OR `record`.view_id = 0)").Order("`record_srv`.priority ASC, `record_srv`.weight DESC").Scan(&srvRecords).Error
	return srvRecords, err
}
