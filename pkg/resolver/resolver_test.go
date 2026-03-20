package resolver

import (
	"context"
	"fmt"
	"net"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/cylonchau/hermes/pkg/model"
)

// setupMockDB initializes a mocked gorm.DB using sqlmock
func setupMockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create sqlmock: %w", err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open gorm: %w", err)
	}

	return gormDB, mock, nil
}

// MockGeoIPProvider is a mock implementation of GeoIPProvider
type MockGeoIPProvider struct {
	LookupFn func(ip string) (string, string, error)
}

func (m *MockGeoIPProvider) Lookup(ip string) (string, string, error) {
	if m.LookupFn != nil {
		return m.LookupFn(ip)
	}
	return "", "", nil
}

// MockDNSQueryRepository is a mock implementation of DNSQueryRepository
type MockDNSQueryRepository struct {
	QueryARecordsFn     func(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.ARecord, error)
	QueryAAAARecordsFn  func(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.AAAARecord, error)
	QueryMXRecordsFn    func(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.MXRecord, error)
	QueryTXTRecordsFn   func(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.TXTRecord, error)
	QuerySOARecordFn    func(ctx context.Context, zoneName string, viewID int64) (*model.SOARecord, error)
	QueryNSRecordsFn    func(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.NSRecord, error)
	QueryCNAMERecordsFn func(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.CNAMERecord, error)
	QuerySRVRecordsFn   func(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.SRVRecord, error)
}

func (m *MockDNSQueryRepository) QueryARecords(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.ARecord, error) {
	if m.QueryARecordsFn != nil {
		return m.QueryARecordsFn(ctx, zoneName, recordName, viewID)
	}
	return nil, nil
}
func (m *MockDNSQueryRepository) QueryAAAARecords(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.AAAARecord, error) {
	if m.QueryAAAARecordsFn != nil {
		return m.QueryAAAARecordsFn(ctx, zoneName, recordName, viewID)
	}
	return nil, nil
}
func (m *MockDNSQueryRepository) QueryMXRecords(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.MXRecord, error) {
	if m.QueryMXRecordsFn != nil {
		return m.QueryMXRecordsFn(ctx, zoneName, recordName, viewID)
	}
	return nil, nil
}
func (m *MockDNSQueryRepository) QueryTXTRecords(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.TXTRecord, error) {
	if m.QueryTXTRecordsFn != nil {
		return m.QueryTXTRecordsFn(ctx, zoneName, recordName, viewID)
	}
	return nil, nil
}
func (m *MockDNSQueryRepository) QuerySOARecord(ctx context.Context, zoneName string, viewID int64) (*model.SOARecord, error) {
	if m.QuerySOARecordFn != nil {
		return m.QuerySOARecordFn(ctx, zoneName, viewID)
	}
	return nil, nil
}
func (m *MockDNSQueryRepository) QueryNSRecords(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.NSRecord, error) {
	if m.QueryNSRecordsFn != nil {
		return m.QueryNSRecordsFn(ctx, zoneName, recordName, viewID)
	}
	return nil, nil
}
func (m *MockDNSQueryRepository) QueryCNAMERecords(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.CNAMERecord, error) {
	if m.QueryCNAMERecordsFn != nil {
		return m.QueryCNAMERecordsFn(ctx, zoneName, recordName, viewID)
	}
	return nil, nil
}
func (m *MockDNSQueryRepository) QuerySRVRecords(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.SRVRecord, error) {
	if m.QuerySRVRecordsFn != nil {
		return m.QuerySRVRecordsFn(ctx, zoneName, recordName, viewID)
	}
	return nil, nil
}

// MockResponseWriter is a mock implementation of dns.ResponseWriter
type MockResponseWriter struct {
	dns.ResponseWriter
	RemoteIP net.IP
}

func (m *MockResponseWriter) RemoteAddr() net.Addr {
	return &net.UDPAddr{IP: m.RemoteIP, Port: 53}
}
func (m *MockResponseWriter) WriteMsg(msg *dns.Msg) error { return nil }

func TestResolver_Resolve(t *testing.T) {
	mockRepo := &MockDNSQueryRepository{}
	db, _, _ := setupMockDB() // We need a GORM DB for matchView
	r := NewResolver(mockRepo, db, nil)

	ctx := context.Background()
	mockW := &MockResponseWriter{RemoteIP: net.ParseIP("127.0.0.1")}

	t.Run("Resolve A Record", func(t *testing.T) {
		mockRepo.QueryARecordsFn = func(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.ARecord, error) {
			if zoneName == "test.com." && recordName == "www.test.com." {
				return []*model.ARecord{
					{IP: 0x01020304, TTL: 3600}, // 1.2.3.4
				}, nil
			}
			return nil, nil
		}

		req := new(dns.Msg)
		req.SetQuestion("www.test.com.", dns.TypeA)
		state := request.Request{W: mockW, Req: req}

		msg, err := r.Resolve(ctx, state)
		if err != nil {
			t.Fatalf("Resolve failed: %v", err)
		}

		if len(msg.Answer) != 1 {
			t.Fatalf("expected 1 answer, got %d", len(msg.Answer))
		}
		if a, ok := msg.Answer[0].(*dns.A); ok {
			if a.A.String() != "1.2.3.4" && a.A.String() != "4.3.2.1" {
				t.Errorf("expected 1.2.3.4, got %s", a.A.String())
			}
		} else {
			t.Errorf("expected *dns.A, got %T", msg.Answer[0])
		}
	})

	t.Run("Resolve NXDOMAIN (Return SOA)", func(t *testing.T) {
		mockRepo.QueryARecordsFn = func(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.ARecord, error) {
			return nil, nil
		}
		mockRepo.QuerySOARecordFn = func(ctx context.Context, zoneName string, viewID int64) (*model.SOARecord, error) {
			if zoneName == "test.com." {
				return &model.SOARecord{
					PrimaryNS: "ns1.test.com.",
					MBox:      "admin.test.com.",
					TTL:       3600,
				}, nil
			}
			return nil, nil
		}

		req := new(dns.Msg)
		req.SetQuestion("missing.test.com.", dns.TypeA)
		state := request.Request{W: mockW, Req: req}

		msg, err := r.Resolve(ctx, state)
		if err != nil {
			t.Fatalf("Resolve failed: %v", err)
		}

		if len(msg.Answer) != 0 {
			t.Errorf("expected 0 answer, got %d", len(msg.Answer))
		}
		if len(msg.Ns) != 1 {
			t.Errorf("expected 1 record in authority section, got %d", len(msg.Ns))
		}
	})

	t.Run("GeoIP match", func(t *testing.T) {
		mockGeoIP := &MockGeoIPProvider{
			LookupFn: func(ip string) (string, string, error) {
				return "CN", "GD", nil
			},
		}
		// Create a resolver with mock DB that has a GeoIP view
		gormDB, sqlMock, _ := setupMockDB()
		rGeo := NewResolver(mockRepo, gormDB, mockGeoIP)

		// Mock View lookup
		viewRows := sqlmock.NewRows([]string{"id", "name", "category", "value", "priority"}).
			AddRow(int64(10), "Guangdong", "geoip", "CN-GD", 10)
		sqlMock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `view` ORDER BY priority DESC")).
			WillReturnRows(viewRows)

		mockRepo.QueryARecordsFn = func(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.ARecord, error) {
			if viewID == 10 {
				return []*model.ARecord{{IP: 16843009, Record: model.Record{TTL: 600}}}, nil
			}
			return nil, nil
		}

		req := new(dns.Msg)
		req.SetQuestion("www.test.com.", dns.TypeA)
		state := request.Request{W: mockW, Req: req}

		msg, err := rGeo.Resolve(ctx, state)
		assert.NoError(t, err)
		assert.NotNil(t, msg)
		assert.NoError(t, sqlMock.ExpectationsWereMet())
	})
}
