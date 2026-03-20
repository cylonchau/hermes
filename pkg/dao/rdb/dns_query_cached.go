package rdb

import (
	"context"
	"encoding/json"

	"github.com/cylonchau/hermes/pkg/dao/memory"
	"github.com/cylonchau/hermes/pkg/model"
)

// CachedDNSQueryRepository is a DNS Query Proxy with L1 Cache
type CachedDNSQueryRepository struct {
	rdb   DNSQueryRepository
	cache *memory.CacheDAO
}

func NewCachedDNSQueryRepository(rdb DNSQueryRepository, cache *memory.CacheDAO) DNSQueryRepository {
	return &CachedDNSQueryRepository{rdb: rdb, cache: cache}
}

func (c *CachedDNSQueryRepository) QueryARecords(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.ARecord, error) {
	if c.cache != nil {
		if bytes, ok := c.cache.Get(zoneName, 1, recordName, viewID); ok {
			if string(bytes) == "[]" {
				return nil, nil // Negative cache hit, prevents penetration
			}
			var res []*model.ARecord
			if err := json.Unmarshal(bytes, &res); err == nil {
				return res, nil
			}
		}
	}
	res, err := c.rdb.QueryARecords(ctx, zoneName, recordName, viewID)
	if err != nil {
		return nil, err
	}

	if c.cache != nil {
		if len(res) > 0 {
			bytes, _ := json.Marshal(res)
			c.cache.Set(zoneName, 1, recordName, viewID, bytes, int(res[0].TTL))
		} else {
			c.cache.Set(zoneName, 1, recordName, viewID, []byte("[]"), 5) // 5s negative cache
		}
	}
	return res, nil
}

func (c *CachedDNSQueryRepository) QueryAAAARecords(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.AAAARecord, error) {
	if c.cache != nil {
		if bytes, ok := c.cache.Get(zoneName, 28, recordName, viewID); ok {
			if string(bytes) == "[]" {
				return nil, nil
			}
			var res []*model.AAAARecord
			if err := json.Unmarshal(bytes, &res); err == nil {
				return res, nil
			}
		}
	}
	res, err := c.rdb.QueryAAAARecords(ctx, zoneName, recordName, viewID)
	if err != nil {
		return nil, err
	}

	if c.cache != nil {
		if len(res) > 0 {
			bytes, _ := json.Marshal(res)
			c.cache.Set(zoneName, 28, recordName, viewID, bytes, int(res[0].TTL))
		} else {
			c.cache.Set(zoneName, 28, recordName, viewID, []byte("[]"), 5)
		}
	}
	return res, nil
}

func (c *CachedDNSQueryRepository) QueryMXRecords(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.MXRecord, error) {
	if c.cache != nil {
		if bytes, ok := c.cache.Get(zoneName, 15, recordName, viewID); ok {
			if string(bytes) == "[]" {
				return nil, nil
			}
			var res []*model.MXRecord
			if err := json.Unmarshal(bytes, &res); err == nil {
				return res, nil
			}
		}
	}
	res, err := c.rdb.QueryMXRecords(ctx, zoneName, recordName, viewID)
	if err != nil {
		return nil, err
	}

	if c.cache != nil {
		if len(res) > 0 {
			bytes, _ := json.Marshal(res)
			c.cache.Set(zoneName, 15, recordName, viewID, bytes, int(res[0].TTL))
		} else {
			c.cache.Set(zoneName, 15, recordName, viewID, []byte("[]"), 5)
		}
	}
	return res, nil
}

func (c *CachedDNSQueryRepository) QueryTXTRecords(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.TXTRecord, error) {
	if c.cache != nil {
		if bytes, ok := c.cache.Get(zoneName, 16, recordName, viewID); ok {
			if string(bytes) == "[]" {
				return nil, nil
			}
			var res []*model.TXTRecord
			if err := json.Unmarshal(bytes, &res); err == nil {
				return res, nil
			}
		}
	}
	res, err := c.rdb.QueryTXTRecords(ctx, zoneName, recordName, viewID)
	if err != nil {
		return nil, err
	}

	if c.cache != nil {
		if len(res) > 0 {
			bytes, _ := json.Marshal(res)
			c.cache.Set(zoneName, 16, recordName, viewID, bytes, int(res[0].TTL))
		} else {
			c.cache.Set(zoneName, 16, recordName, viewID, []byte("[]"), 5)
		}
	}
	return res, nil
}

func (c *CachedDNSQueryRepository) QuerySOARecord(ctx context.Context, zoneName string, viewID int64) (*model.SOARecord, error) {
	if c.cache != nil {
		if bytes, ok := c.cache.Get(zoneName, 6, zoneName, viewID); ok {
			if string(bytes) == "[]" {
				return nil, nil
			}
			var res *model.SOARecord
			if err := json.Unmarshal(bytes, &res); err == nil {
				return res, nil
			}
		}
	}
	res, err := c.rdb.QuerySOARecord(ctx, zoneName, viewID)
	if err != nil {
		return nil, err
	}

	if c.cache != nil {
		if res != nil {
			bytes, _ := json.Marshal(res)
			c.cache.Set(zoneName, 6, zoneName, viewID, bytes, int(res.TTL))
		} else {
			c.cache.Set(zoneName, 6, zoneName, viewID, []byte("[]"), 5)
		}
	}
	return res, nil
}

func (c *CachedDNSQueryRepository) QueryNSRecords(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.NSRecord, error) {
	if c.cache != nil {
		if bytes, ok := c.cache.Get(zoneName, 2, recordName, viewID); ok {
			if string(bytes) == "[]" {
				return nil, nil
			}
			var res []*model.NSRecord
			if err := json.Unmarshal(bytes, &res); err == nil {
				return res, nil
			}
		}
	}
	res, err := c.rdb.QueryNSRecords(ctx, zoneName, recordName, viewID)
	if err != nil {
		return nil, err
	}

	if c.cache != nil {
		if len(res) > 0 {
			bytes, _ := json.Marshal(res)
			c.cache.Set(zoneName, 2, recordName, viewID, bytes, int(res[0].TTL))
		} else {
			c.cache.Set(zoneName, 2, recordName, viewID, []byte("[]"), 5)
		}
	}
	return res, nil
}

func (c *CachedDNSQueryRepository) QueryCNAMERecords(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.CNAMERecord, error) {
	if c.cache != nil {
		if bytes, ok := c.cache.Get(zoneName, 5, recordName, viewID); ok {
			if string(bytes) == "[]" {
				return nil, nil
			}
			var res []*model.CNAMERecord
			if err := json.Unmarshal(bytes, &res); err == nil {
				return res, nil
			}
		}
	}
	res, err := c.rdb.QueryCNAMERecords(ctx, zoneName, recordName, viewID)
	if err != nil {
		return nil, err
	}

	if c.cache != nil {
		if len(res) > 0 {
			bytes, _ := json.Marshal(res)
			c.cache.Set(zoneName, 5, recordName, viewID, bytes, int(res[0].TTL))
		} else {
			c.cache.Set(zoneName, 5, recordName, viewID, []byte("[]"), 5)
		}
	}
	return res, nil
}

func (c *CachedDNSQueryRepository) QuerySRVRecords(ctx context.Context, zoneName, recordName string, viewID int64) ([]*model.SRVRecord, error) {
	if c.cache != nil {
		if bytes, ok := c.cache.Get(zoneName, 33, recordName, viewID); ok {
			if string(bytes) == "[]" {
				return nil, nil
			}
			var res []*model.SRVRecord
			if err := json.Unmarshal(bytes, &res); err == nil {
				return res, nil
			}
		}
	}
	res, err := c.rdb.QuerySRVRecords(ctx, zoneName, recordName, viewID)
	if err != nil {
		return nil, err
	}

	if c.cache != nil {
		if len(res) > 0 {
			bytes, _ := json.Marshal(res)
			c.cache.Set(zoneName, 33, recordName, viewID, bytes, int(res[0].TTL))
		} else {
			c.cache.Set(zoneName, 33, recordName, viewID, []byte("[]"), 5)
		}
	}
	return res, nil
}
