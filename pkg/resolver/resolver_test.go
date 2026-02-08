package resolver

import (
	"context"
	"testing"

	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"

	"github.com/cylonchau/hermes/pkg/model"
)

// MockDNSQueryRepository is a mock implementation of DNSQueryRepository
type MockDNSQueryRepository struct {
	QueryARecordsFn     func(ctx context.Context, zoneName, recordName string) ([]*model.ARecord, error)
	QueryAAAARecordsFn  func(ctx context.Context, zoneName, recordName string) ([]*model.AAAARecord, error)
	QueryMXRecordsFn    func(ctx context.Context, zoneName, recordName string) ([]*model.MXRecord, error)
	QueryTXTRecordsFn   func(ctx context.Context, zoneName, recordName string) ([]*model.TXTRecord, error)
	QuerySOARecordFn    func(ctx context.Context, zoneName string) (*model.SOARecord, error)
	QueryNSRecordsFn    func(ctx context.Context, zoneName, recordName string) ([]*model.NSRecord, error)
	QueryCNAMERecordsFn func(ctx context.Context, zoneName, recordName string) ([]*model.CNAMERecord, error)
	QuerySRVRecordsFn   func(ctx context.Context, zoneName, recordName string) ([]*model.SRVRecord, error)
}

func (m *MockDNSQueryRepository) QueryARecords(ctx context.Context, zoneName, recordName string) ([]*model.ARecord, error) {
	if m.QueryARecordsFn != nil {
		return m.QueryARecordsFn(ctx, zoneName, recordName)
	}
	return nil, nil
}
func (m *MockDNSQueryRepository) QueryAAAARecords(ctx context.Context, zoneName, recordName string) ([]*model.AAAARecord, error) {
	if m.QueryAAAARecordsFn != nil {
		return m.QueryAAAARecordsFn(ctx, zoneName, recordName)
	}
	return nil, nil
}
func (m *MockDNSQueryRepository) QueryMXRecords(ctx context.Context, zoneName, recordName string) ([]*model.MXRecord, error) {
	if m.QueryMXRecordsFn != nil {
		return m.QueryMXRecordsFn(ctx, zoneName, recordName)
	}
	return nil, nil
}
func (m *MockDNSQueryRepository) QueryTXTRecords(ctx context.Context, zoneName, recordName string) ([]*model.TXTRecord, error) {
	if m.QueryTXTRecordsFn != nil {
		return m.QueryTXTRecordsFn(ctx, zoneName, recordName)
	}
	return nil, nil
}
func (m *MockDNSQueryRepository) QuerySOARecord(ctx context.Context, zoneName string) (*model.SOARecord, error) {
	if m.QuerySOARecordFn != nil {
		return m.QuerySOARecordFn(ctx, zoneName)
	}
	return nil, nil
}
func (m *MockDNSQueryRepository) QueryNSRecords(ctx context.Context, zoneName, recordName string) ([]*model.NSRecord, error) {
	if m.QueryNSRecordsFn != nil {
		return m.QueryNSRecordsFn(ctx, zoneName, recordName)
	}
	return nil, nil
}
func (m *MockDNSQueryRepository) QueryCNAMERecords(ctx context.Context, zoneName, recordName string) ([]*model.CNAMERecord, error) {
	if m.QueryCNAMERecordsFn != nil {
		return m.QueryCNAMERecordsFn(ctx, zoneName, recordName)
	}
	return nil, nil
}
func (m *MockDNSQueryRepository) QuerySRVRecords(ctx context.Context, zoneName, recordName string) ([]*model.SRVRecord, error) {
	if m.QuerySRVRecordsFn != nil {
		return m.QuerySRVRecordsFn(ctx, zoneName, recordName)
	}
	return nil, nil
}

func TestResolver_Resolve(t *testing.T) {
	mockRepo := &MockDNSQueryRepository{}
	r := NewResolver(mockRepo)
	ctx := context.Background()

	t.Run("Resolve A Record", func(t *testing.T) {
		mockRepo.QueryARecordsFn = func(ctx context.Context, zoneName, recordName string) ([]*model.ARecord, error) {
			if zoneName == "test.com" && recordName == "www" {
				return []*model.ARecord{
					{IP: 0x01020304, TTL: 3600},
				}, nil
			}
			return nil, nil
		}

		req := new(dns.Msg)
		req.SetQuestion("www.test.com.", dns.TypeA)
		state := request.Request{W: nil, Req: req}

		msg, err := r.Resolve(ctx, state)
		if err != nil {
			t.Fatalf("Resolve failed: %v", err)
		}

		if len(msg.Answer) != 1 {
			t.Fatalf("expected 1 answer, got %d", len(msg.Answer))
		}
		if a, ok := msg.Answer[0].(*dns.A); ok {
			if a.A.String() != "1.2.3.4" {
				t.Errorf("expected 1.2.3.4, got %s", a.A.String())
			}
		} else {
			t.Errorf("expected *dns.A, got %T", msg.Answer[0])
		}
	})

	t.Run("Resolve NXDOMAIN (Return SOA)", func(t *testing.T) {
		mockRepo.QueryARecordsFn = func(ctx context.Context, zoneName, recordName string) ([]*model.ARecord, error) {
			return nil, nil
		}
		mockRepo.QuerySOARecordFn = func(ctx context.Context, zoneName string) (*model.SOARecord, error) {
			if zoneName == "test.com" {
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
		state := request.Request{W: nil, Req: req}

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
}
