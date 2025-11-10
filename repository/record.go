package repository

import (
	"context"

	"github.com/cylonchau/hermes/pkg/model"
)

// RecordRepository 基础记录仓储接口
type RecordRepository interface {
	CreateRecord(ctx context.Context, record *model.Record) error
	GetRecordByID(ctx context.Context, recordID uint) (*model.Record, error)
	GetRecordsByZone(ctx context.Context, zoneID uint) ([]*model.Record, error)
	GetRecordsByName(ctx context.Context, zoneID uint, recordName string) ([]*model.Record, error)
	UpdateRecord(ctx context.Context, record *model.Record) error
	DeleteRecord(ctx context.Context, recordID uint) error
	SoftDeleteRecord(ctx context.Context, recordID uint) error
	CountRecordsByZone(ctx context.Context, zoneID uint) (int64, error)
	BatchCreateRecords(ctx context.Context, records []*model.Record) error
	BatchDeleteRecords(ctx context.Context, recordIDs []uint) error
}

// ARecordRepository A记录专用仓储接口
type ARecordRepository interface {
	CreateARecord(ctx context.Context, record *model.Record, aRecord *model.ARecord) error
	GetARecords(ctx context.Context, zoneName, recordName string) ([]*model.ARecord, error)
	GetARecordByID(ctx context.Context, recordID uint) (*model.ARecord, error)
	UpdateARecord(ctx context.Context, record *model.Record, aRecord *model.ARecord) error
	DeleteARecord(ctx context.Context, recordID uint) error
}

// AAAARecordRepository AAAA记录专用仓储接口
type AAAARecordRepository interface {
	CreateAAAARecord(ctx context.Context, record *model.Record, aaaaRecord *model.AAAARecord) error
	GetAAAARecords(ctx context.Context, zoneName, recordName string) ([]*model.AAAARecord, error)
	GetAAAARecordByID(ctx context.Context, recordID uint) (*model.AAAARecord, error)
	UpdateAAAARecord(ctx context.Context, record *model.Record, aaaaRecord *model.AAAARecord) error
	DeleteAAAARecord(ctx context.Context, recordID uint) error
}

// MXRecordRepository MX记录专用仓储接口
type MXRecordRepository interface {
	CreateMXRecord(ctx context.Context, record *model.Record, mxRecord *model.MXRecord) error
	GetMXRecords(ctx context.Context, zoneName, recordName string) ([]*model.MXRecord, error)
	GetMXRecordByID(ctx context.Context, recordID uint) (*model.MXRecord, error)
	UpdateMXRecord(ctx context.Context, record *model.Record, mxRecord *model.MXRecord) error
	DeleteMXRecord(ctx context.Context, recordID uint) error
}

// TXTRecordRepository TXT记录专用仓储接口
type TXTRecordRepository interface {
	CreateTXTRecord(ctx context.Context, record *model.Record, txtRecord *model.TXTRecord) error
	GetTXTRecords(ctx context.Context, zoneName, recordName string) ([]*model.TXTRecord, error)
	GetTXTRecordByID(ctx context.Context, recordID uint) (*model.TXTRecord, error)
	UpdateTXTRecord(ctx context.Context, record *model.Record, txtRecord *model.TXTRecord) error
	DeleteTXTRecord(ctx context.Context, recordID uint) error
}

// SOARecordRepository SOA记录专用仓储接口
type SOARecordRepository interface {
	CreateSOARecord(ctx context.Context, record *model.Record, soaRecord *model.SOARecord) error
	GetSOARecord(ctx context.Context, zoneName string) (*model.SOARecord, error)
	GetSOARecordByID(ctx context.Context, recordID uint) (*model.SOARecord, error)
	UpdateSOARecord(ctx context.Context, record *model.Record, soaRecord *model.SOARecord) error
	DeleteSOARecord(ctx context.Context, recordID uint) error
}

// NSRecordRepository NS记录专用仓储接口
type NSRecordRepository interface {
	CreateNSRecord(ctx context.Context, record *model.Record, nsRecord *model.NSRecord) error
	GetNSRecords(ctx context.Context, zoneName, recordName string) ([]*model.NSRecord, error)
	GetNSRecordByID(ctx context.Context, recordID uint) (*model.NSRecord, error)
	UpdateNSRecord(ctx context.Context, record *model.Record, nsRecord *model.NSRecord) error
	DeleteNSRecord(ctx context.Context, recordID uint) error
}

// CNAMERecordRepository CNAME记录专用仓储接口
type CNAMERecordRepository interface {
	CreateCNAMERecord(ctx context.Context, record *model.Record, cnameRecord *model.CNAMERecord) error
	GetCNAMERecords(ctx context.Context, zoneName, recordName string) ([]*model.CNAMERecord, error)
	GetCNAMERecordByID(ctx context.Context, recordID uint) (*model.CNAMERecord, error)
	UpdateCNAMERecord(ctx context.Context, record *model.Record, cnameRecord *model.CNAMERecord) error
	DeleteCNAMERecord(ctx context.Context, recordID uint) error
}

// SRVRecordRepository SRV记录专用仓储接口
type SRVRecordRepository interface {
	CreateSRVRecord(ctx context.Context, record *model.Record, srvRecord *model.SRVRecord) error
	GetSRVRecords(ctx context.Context, zoneName, recordName string) ([]*model.SRVRecord, error)
	GetSRVRecordByID(ctx context.Context, recordID uint) (*model.SRVRecord, error)
	UpdateSRVRecord(ctx context.Context, record *model.Record, srvRecord *model.SRVRecord) error
	DeleteSRVRecord(ctx context.Context, recordID uint) error
}

// DNSQueryRepository CoreDNS专用查询接口（只读，高性能）
type DNSQueryRepository interface {
	QueryARecords(ctx context.Context, zoneName, recordName string) ([]*model.ARecord, error)
	QueryAAAARecords(ctx context.Context, zoneName, recordName string) ([]*model.AAAARecord, error)
	QueryMXRecords(ctx context.Context, zoneName, recordName string) ([]*model.MXRecord, error)
	QueryTXTRecords(ctx context.Context, zoneName, recordName string) ([]*model.TXTRecord, error)
	QuerySOARecord(ctx context.Context, zoneName string) (*model.SOARecord, error)
	QueryNSRecords(ctx context.Context, zoneName, recordName string) ([]*model.NSRecord, error)
	QueryCNAMERecords(ctx context.Context, zoneName, recordName string) ([]*model.CNAMERecord, error)
	QuerySRVRecords(ctx context.Context, zoneName, recordName string) ([]*model.SRVRecord, error)
}
