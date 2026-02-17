package resolver

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"strings"

	"github.com/coredns/coredns/request"
	"github.com/cylonchau/hermes/pkg/dao/rdb"
	"github.com/miekg/dns"
)

// Resolver DNS 解析核心处理器
type Resolver struct {
	dao rdb.DNSQueryRepository
}

// NewResolver 创建解析器实例
func NewResolver(dao rdb.DNSQueryRepository) *Resolver {
	return &Resolver{dao: dao}
}

// Resolve 处理 DNS 解析逻辑
func (r *Resolver) Resolve(ctx context.Context, state request.Request) (*dns.Msg, error) {
	qName := state.Name()
	qType := state.QType()

	// 1. 查找最匹配的 Zone
	zone, name, err := r.parseQuery(ctx, qName)
	if err != nil {
		return nil, err
	}

	m := new(dns.Msg)
	m.SetReply(state.Req)
	m.Authoritative = true

	// 2. 根据查询类型检索记录
	switch qType {
	case dns.TypeA:
		records, err := r.dao.QueryARecords(ctx, zone, name)
		if err == nil {
			for _, rec := range records {
				ip := make(net.IP, 4)
				binary.BigEndian.PutUint32(ip, uint32(rec.IP))
				m.Answer = append(m.Answer, &dns.A{
					Hdr: dns.RR_Header{Name: qName, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: rec.TTL},
					A:   ip,
				})
			}
		}
	case dns.TypeAAAA:
		records, err := r.dao.QueryAAAARecords(ctx, zone, name)
		if err == nil {
			for _, rec := range records {
				m.Answer = append(m.Answer, &dns.AAAA{
					Hdr:  dns.RR_Header{Name: qName, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: rec.TTL},
					AAAA: net.IP(rec.IP),
				})
			}
		}
	case dns.TypeCNAME:
		records, err := r.dao.QueryCNAMERecords(ctx, zone, name)
		if err == nil {
			for _, rec := range records {
				m.Answer = append(m.Answer, &dns.CNAME{
					Hdr:    dns.RR_Header{Name: qName, Rrtype: dns.TypeCNAME, Class: dns.ClassINET, Ttl: rec.TTL},
					Target: dns.Fqdn(rec.Target),
				})
			}
		}
	case dns.TypeMX:
		records, err := r.dao.QueryMXRecords(ctx, zone, name)
		if err == nil {
			for _, rec := range records {
				m.Answer = append(m.Answer, &dns.MX{
					Hdr:        dns.RR_Header{Name: qName, Rrtype: dns.TypeMX, Class: dns.ClassINET, Ttl: rec.TTL},
					Preference: uint16(rec.Priority),
					Mx:         dns.Fqdn(rec.Host),
				})
			}
		}
	case dns.TypeTXT:
		records, err := r.dao.QueryTXTRecords(ctx, zone, name)
		if err == nil {
			for _, rec := range records {
				m.Answer = append(m.Answer, &dns.TXT{
					Hdr: dns.RR_Header{Name: qName, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: rec.TTL},
					Txt: []string{rec.Text},
				})
			}
		}
	case dns.TypeNS:
		records, err := r.dao.QueryNSRecords(ctx, zone, name)
		if err == nil {
			for _, rec := range records {
				m.Answer = append(m.Answer, &dns.NS{
					Hdr: dns.RR_Header{Name: qName, Rrtype: dns.TypeNS, Class: dns.ClassINET, Ttl: rec.TTL},
					Ns:  dns.Fqdn(rec.NameServer),
				})
			}
		}
	case dns.TypeSOA:
		rec, err := r.dao.QuerySOARecord(ctx, zone)
		if err == nil {
			m.Answer = append(m.Answer, &dns.SOA{
				Hdr:     dns.RR_Header{Name: dns.Fqdn(zone), Rrtype: dns.TypeSOA, Class: dns.ClassINET, Ttl: rec.TTL},
				Ns:      dns.Fqdn(rec.PrimaryNS),
				Mbox:    dns.Fqdn(rec.MBox),
				Serial:  rec.Serial,
				Refresh: rec.Refresh,
				Retry:   rec.Retry,
				Expire:  rec.Expire,
				Minttl:  rec.MinTTL,
			})
		}
	case dns.TypeSRV:
		records, err := r.dao.QuerySRVRecords(ctx, zone, name)
		if err == nil {
			for _, rec := range records {
				m.Answer = append(m.Answer, &dns.SRV{
					Hdr:      dns.RR_Header{Name: qName, Rrtype: dns.TypeSRV, Class: dns.ClassINET, Ttl: rec.TTL},
					Priority: uint16(rec.Priority),
					Weight:   uint16(rec.Weight),
					Port:     uint16(rec.Port),
					Target:   dns.Fqdn(rec.Target),
				})
			}
		}
	}

	// 3. 如果未找到记录且是 NXDOMAIN 情况
	if len(m.Answer) == 0 {
		return r.handleNoData(ctx, zone, m)
	}

	return m, nil
}

// parseQuery 解析查询域名，将其拆分为 Zone 和 Record Name
func (r *Resolver) parseQuery(ctx context.Context, qName string) (zone, name string, err error) {
	// 暂时使用简单的点号拆分方案
	labels := dns.SplitDomainName(qName)
	if len(labels) < 2 {
		return "", "", fmt.Errorf("invalid domain name: %s", qName)
	}

	// 使用 FQDN 格式（带末尾点号）以匹配数据库中的存储规范
	zone = strings.Join(labels[len(labels)-2:], ".") + "."
	name = qName

	return zone, name, nil
}

// handleNoData 处理无数据情况，补充 SOA 到 Authority 段
func (r *Resolver) handleNoData(ctx context.Context, zone string, m *dns.Msg) (*dns.Msg, error) {
	rec, err := r.dao.QuerySOARecord(ctx, zone)
	if err == nil {
		m.Ns = append(m.Ns, &dns.SOA{
			Hdr:     dns.RR_Header{Name: dns.Fqdn(zone), Rrtype: dns.TypeSOA, Class: dns.ClassINET, Ttl: rec.TTL},
			Ns:      dns.Fqdn(rec.PrimaryNS),
			Mbox:    dns.Fqdn(rec.MBox),
			Serial:  rec.Serial,
			Refresh: rec.Refresh,
			Retry:   rec.Retry,
			Expire:  rec.Expire,
			Minttl:  rec.MinTTL,
		})
	}
	return m, nil
}
