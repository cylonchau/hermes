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

// Resolver is the core DNS resolving processor
type Resolver struct {
	dao   rdb.DNSQueryRepository
	db    *gorm.DB
	geoip GeoIPProvider
}

// NewResolver creates a resolver instance
func NewResolver(dao rdb.DNSQueryRepository, db *gorm.DB, geoip GeoIPProvider) *Resolver {
	return &Resolver{dao: dao, db: db, geoip: geoip}
}

// Resolve handles DNS resolution logic
func (r *Resolver) Resolve(ctx context.Context, state request.Request) (*dns.Msg, error) {
	qName := state.Name()
	qType := state.QType()
	clientIP := state.IP()

	// 1. Identify view
	viewID, _ := r.matchView(ctx, clientIP)

	// 2. Find most matching Zone
	zone, name, err := r.parseQuery(ctx, qName)
	if err != nil {
		return nil, err
	}

	m := new(dns.Msg)
	m.SetReply(state.Req)
	m.Authoritative = true

	// 3. Retrieve records based on query type
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

	// 4. If record is not found or is NXDOMAIN state
	if len(m.Answer) == 0 {
		return r.handleNoData(ctx, zone, viewID, m)
	}

	return m, nil
}

// parseQuery parses query domain name, splitting into Zone and Record Name
func (r *Resolver) parseQuery(ctx context.Context, qName string) (zone, name string, err error) {
	// Currently uses a simple dot split approach
	labels := dns.SplitDomainName(qName)
	if len(labels) < 2 {
		return "", "", fmt.Errorf("invalid domain name: %s", qName)
	}

	// Use FQDN format (with trailing dot) to match database storage spec
	zone = strings.Join(labels[len(labels)-2:], ".") + "."
	name = qName

	return zone, name, nil
}

// handleNoData handles NO DATA states appending SOA to Authority section
func (r *Resolver) handleNoData(ctx context.Context, zone string, viewID int64, m *dns.Msg) (*dns.Msg, error) {
	rec, err := r.dao.QuerySOARecord(ctx, zone, viewID)
	if err == nil && rec != nil && rec.ID > 0 {
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
	} else {
		// Explicitly mark: completely absent within server authority
		m.Rcode = dns.RcodeNameError
	}
	return m, nil
}



// matchView matches View based on client IP
func (r *Resolver) matchView(ctx context.Context, clientIP string) (int64, error) {
	// 1. Fetch all views and sort by priority
	// Note: cache support should be added later to avoid full-scans
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

	// 2. Traverse views for match
	for _, v := range views {
		switch v.Category {
		case "acl":
			// Handle CIDR lists (supports comma separated)
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

			// Supports country code (e.g., CN) or city/region code (e.g., CN-GD)
			// Note: Value definitions must align with GeoIP datasource
			if v.Value != "" && (v.Value == country || v.Value == country+"-"+region || v.Value == region) {
				return v.ID, nil
			}
		}
	}

	return 0, nil
}
