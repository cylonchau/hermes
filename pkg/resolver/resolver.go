package resolver

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"net/netip"
	"strings"

	"github.com/coredns/coredns/request"
	"github.com/cylonchau/hermes/pkg/dao/rdb"
	"github.com/cylonchau/hermes/pkg/model"
	"github.com/miekg/dns"
	"gorm.io/gorm"
)

// GeoIPProvider defines the interface for GeoIP lookups
type GeoIPProvider interface {
	Lookup(ip string) (country, region string, err error)
}

// Resolver DNS 解析核心处理器
type Resolver struct {
	dao   rdb.DNSQueryRepository
	db    *gorm.DB
	geoip GeoIPProvider
}

// NewResolver 创建解析器实例
func NewResolver(dao rdb.DNSQueryRepository, db *gorm.DB, geoip GeoIPProvider) *Resolver {
	return &Resolver{dao: dao, db: db, geoip: geoip}
}

// Resolve 处理 DNS 解析逻辑
func (r *Resolver) Resolve(ctx context.Context, state request.Request) (*dns.Msg, error) {
	qName := state.Name()
	qType := state.QType()
	clientIP := state.IP()

	// 1. 识别视图
	viewID, _ := r.matchView(ctx, clientIP)

	// 2. 查找最匹配的 Zone
	zone, name, err := r.parseQuery(ctx, qName)
	if err != nil {
		return nil, err
	}

	m := new(dns.Msg)
	m.SetReply(state.Req)
	m.Authoritative = true

	// 3. 根据查询类型检索记录
	switch qType {
	case dns.TypeA:
		records, err := r.dao.QueryARecords(ctx, zone, name, viewID)
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
		records, err := r.dao.QueryAAAARecords(ctx, zone, name, viewID)
		if err == nil {
			for _, rec := range records {
				m.Answer = append(m.Answer, &dns.AAAA{
					Hdr:  dns.RR_Header{Name: qName, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: rec.TTL},
					AAAA: net.IP(rec.IP),
				})
			}
		}
	case dns.TypeCNAME:
		records, err := r.dao.QueryCNAMERecords(ctx, zone, name, viewID)
		if err == nil {
			for _, rec := range records {
				m.Answer = append(m.Answer, &dns.CNAME{
					Hdr:    dns.RR_Header{Name: qName, Rrtype: dns.TypeCNAME, Class: dns.ClassINET, Ttl: rec.TTL},
					Target: dns.Fqdn(rec.Target),
				})
			}
		}
	case dns.TypeMX:
		records, err := r.dao.QueryMXRecords(ctx, zone, name, viewID)
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
		records, err := r.dao.QueryTXTRecords(ctx, zone, name, viewID)
		if err == nil {
			for _, rec := range records {
				m.Answer = append(m.Answer, &dns.TXT{
					Hdr: dns.RR_Header{Name: qName, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: rec.TTL},
					Txt: []string{rec.Text},
				})
			}
		}
	case dns.TypeNS:
		records, err := r.dao.QueryNSRecords(ctx, zone, name, viewID)
		if err == nil {
			for _, rec := range records {
				m.Answer = append(m.Answer, &dns.NS{
					Hdr: dns.RR_Header{Name: qName, Rrtype: dns.TypeNS, Class: dns.ClassINET, Ttl: rec.TTL},
					Ns:  dns.Fqdn(rec.NameServer),
				})
			}
		}
	case dns.TypeSOA:
		rec, err := r.dao.QuerySOARecord(ctx, zone, viewID)
		if err == nil && rec != nil {
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
		records, err := r.dao.QuerySRVRecords(ctx, zone, name, viewID)
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

	// 4. 如果未找到记录且是 NXDOMAIN 情况
	if len(m.Answer) == 0 {
		return r.handleNoData(ctx, zone, viewID, m)
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
func (r *Resolver) handleNoData(ctx context.Context, zone string, viewID int64, m *dns.Msg) (*dns.Msg, error) {
	rec, err := r.dao.QuerySOARecord(ctx, zone, viewID)
	if err == nil && rec != nil {
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

// matchView 根据客户端 IP 匹配视图
func (r *Resolver) matchView(ctx context.Context, clientIP string) (int64, error) {
	// 1. 获取所有视图并按优先级排序
	// 注意：这里后续应该增加缓存机制，避免每次查询都扫库
	var views []model.View
	db := r.db
	if db == nil {
		db = model.DB
	}
	if db == nil {
		return 0, fmt.Errorf("database not initialized")
	}
	err := db.WithContext(ctx).Order("priority DESC").Find(&views).Error
	if err != nil {
		return 0, err
	}

	ip, err := netip.ParseAddr(clientIP)
	if err != nil {
		return 0, err
	}

	// 2. 遍历视图进行匹配
	for _, v := range views {
		switch v.Category {
		case "acl":
			// 处理 CIDR 列表（支持逗号分隔）
			cidrs := strings.Split(v.Value, ",")
			for _, cidrStr := range cidrs {
				cidrStr = strings.TrimSpace(cidrStr)
				if cidrStr == "" {
					continue
				}
				prefix, err := netip.ParsePrefix(cidrStr)
				if err == nil {
					if prefix.Contains(ip) {
						return v.ID, nil
					}
				}
			}
		case "geoip":
			if r.geoip == nil {
				continue
			}
			country, region, err := r.geoip.Lookup(clientIP)
			if err != nil {
				continue
			}

			// 支持国家代码 (如 CN) 或 城市/区域代码 (如 CN-GD)
			// 注意：具体的 Value 定义需要与 GeoIP 数据源对齐
			if v.Value != "" && (v.Value == country || v.Value == country+"-"+region || v.Value == region) {
				return v.ID, nil
			}
		}
	}

	return 0, nil
}
