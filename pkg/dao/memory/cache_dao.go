package memory

import (
	"fmt"
	"sync"

	"github.com/coocood/freecache"
)

// CacheDAO represents the L1 memory cache (backed by FreeCache with Native TTL).
type CacheDAO struct {
	cache     *freecache.Cache
	serialMap sync.Map   // Structure: map[string]int64 (Zone Name -> Serial Version)
	zoneKeys  sync.Map   // Secondary Index: map[string][]string (Zone Name -> List of FreeCache Keys)
	lock      sync.Mutex // Protects zoneKeys Concurrent slice appending and iteration
}

// NewCacheDAO initializes CacheDAO. maxBytes is the maximum allocated memory in bytes (e.g., 32 * 1024 * 1024).
func NewCacheDAO(maxBytes int) *CacheDAO {
	if maxBytes <= 0 {
		maxBytes = 32 * 1024 * 1024 // Default 32MB
	}
	return &CacheDAO{
		cache: freecache.NewCache(maxBytes),
	}
}

// Get fetches cached DNS record.
func (d *CacheDAO) Get(zone string, qType uint16, qName string, viewID int64) ([]byte, bool) {
	serial := d.getSerial(zone)
	// Build logical key: [Zone]_[Serial]_[QType]_[ViewID]_[QName]
	key := fmt.Sprintf("%s_%d_%d_%d_%s", zone, serial, qType, viewID, qName)

	val, err := d.cache.Get([]byte(key))
	if err != nil {
		return nil, false // err is usually freecache.ErrNotFound
	}
	return val, true
}

// Set writes cached DNS record (ttlSeconds is the expiration time).
func (d *CacheDAO) Set(zone string, qType uint16, qName string, viewID int64, value []byte, ttlSeconds int) {
	if ttlSeconds <= 0 {
		ttlSeconds = 60 // Default fallback 60 seconds
	}

	serial := d.getSerial(zone)
	key := fmt.Sprintf("%s_%d_%d_%d_%s", zone, serial, qType, viewID, qName)

	// Direct Native TTL Set
	d.cache.Set([]byte(key), value, ttlSeconds)

	// Record to secondary index for physical sweeping invalidation
	d.lock.Lock()
	var list []string
	if val, ok := d.zoneKeys.Load(zone); ok {
		list = val.([]string)
	}
	d.zoneKeys.Store(zone, append(list, key))
	d.lock.Unlock()
}

// UpdateZoneSerial updates the version (Serial) of a zone and triggers physical sweeping evacuation.
func (d *CacheDAO) UpdateZoneSerial(zone string, serial int64) {
	d.serialMap.Store(zone, serial)

	// Actively sweep and clear all orphaned keys belonging to this zone!
	d.lock.Lock()
	if val, ok := d.zoneKeys.Load(zone); ok {
		keys := val.([]string)
		for _, k := range keys {
			d.cache.Del([]byte(k)) // Physical removal to free ring buffer space
		}
		d.zoneKeys.Delete(zone) // Clear index chain for this zone
	}
	d.lock.Unlock()
}

// getSerial is an internal helper that fetches current version of a Zone. Returns 0 if absent.
func (d *CacheDAO) getSerial(zone string) int64 {
	if val, ok := d.serialMap.Load(zone); ok {
		return val.(int64)
	}

	return 0
}
