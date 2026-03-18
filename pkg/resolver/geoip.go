package resolver

import (
	"fmt"
	"net"

	"github.com/oschwald/geoip2-golang"
)

type MaxMindProvider struct {
	db *geoip2.Reader
}

func NewMaxMindProvider(dbPath string) (*MaxMindProvider, error) {
	db, err := geoip2.Open(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open maxmind db: %w", err)
	}
	return &MaxMindProvider{db: db}, nil
}

func (p *MaxMindProvider) Lookup(ipStr string) (countryCode string, regionCode string, err error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "", "", fmt.Errorf("invalid IP address: %s", ipStr)
	}

	record, err := p.db.City(ip)
	if err != nil {
		return "", "", err
	}

	countryCode = record.Country.IsoCode
	if len(record.Subdivisions) > 0 {
		regionCode = record.Subdivisions[0].IsoCode
	}

	return countryCode, regionCode, nil
}

func (p *MaxMindProvider) Close() error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}
