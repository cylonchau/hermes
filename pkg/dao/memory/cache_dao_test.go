package memory

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCacheDAO_Get_Set(t *testing.T) {
	dao := NewCacheDAO(1024 * 1024) // 1MB 
	zone := "test.com"

	// 1. Set normal cache (10s)
	dao.Set(zone, 1, "www.test.com", 0, []byte("1.1.1.1"), 10)

	// 2. Get
	res, ok := dao.Get(zone, 1, "www.test.com", 0)
	assert.True(t, ok)
	assert.Equal(t, []byte("1.1.1.1"), res)
}

func TestCacheDAO_TTLExpire(t *testing.T) {
	dao := NewCacheDAO(1024 * 1024)
	zone := "test.com"

	// 1. Set short TTL cache (simulate expiration by waiting 1s, or verify with lower-level hooks)
	// Write 1s cache
	dao.Set(zone, 1, "www.test.com", 0, []byte("1.1.1.1"), 1)

	// Wait 1.1s to allow expiration in FreeCache
	time.Sleep(1100 * time.Millisecond)

	// 2. Get, expect to miss/fail
	res, ok := dao.Get(zone, 1, "www.test.com", 0)
	assert.False(t, ok)
	assert.Nil(t, res)
}

func TestCacheDAO_UpdateZoneSerial(t *testing.T) {
	dao := NewCacheDAO(1024 * 1024)
	zone := "test.com"

	// 1. Write initial cache (Serial = 0)
	dao.Set(zone, 1, "www.test.com", 0, []byte("1.1.1.1"), 10)
	
	res1, ok1 := dao.Get(zone, 1, "www.test.com", 0)
	assert.True(t, ok1)
	assert.Equal(t, []byte("1.1.1.1"), res1)

	// Prepare the next Serial version: since `serialMap` is joined with current serial during Set,
	// so we update serial first to trigger sweeping, and verify if it can still be fetched.
	
	// 2. Simulate platform update: upgrade Serial to 1
	dao.UpdateZoneSerial(zone, 1)

	// 3. Query same QName again, expect Miss (because the original Key `test.com_0_1_0_www.test.com` will be actively deleted)
	res2, ok2 := dao.Get(zone, 1, "www.test.com", 0)
	assert.False(t, ok2)
	assert.Nil(t, res2)
}

func TestCacheDAO_ClearAndReload(t *testing.T) {
	dao := NewCacheDAO(1024 * 1024)
	zone := "test.com"

	// Simulate writing multiple caches under same Zone
	for i := 1; i <= 5; i++ {
		name := fmt.Sprintf("www%d.test.com", i)
		dao.Set(zone, 1, name, 0, []byte(fmt.Sprintf("1.1.1.%d", i)), 10)
	}

	// Verify all hit
	for i := 1; i <= 5; i++ {
		name := fmt.Sprintf("www%d.test.com", i)
		res, ok := dao.Get(zone, 1, name, 0)
		assert.True(t, ok)
		assert.Equal(t, []byte(fmt.Sprintf("1.1.1.%d", i)), res)
	}

	// Update Zone Serial version
	dao.UpdateZoneSerial(zone, 2)

	// Verify all evicted/empty
	for i := 1; i <= 5; i++ {
		name := fmt.Sprintf("www%d.test.com", i)
		_, ok := dao.Get(zone, 1, name, 0)
		assert.False(t, ok, "Expected name %s to be evicted", name)
	}
}
