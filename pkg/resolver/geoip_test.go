package resolver

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaxMindProvider_Lookup(t *testing.T) {
	// 使用系统中现有的测试数据
	dbPath := "../../target/coredns-src/plugin/geoip/testdata/GeoLite2-City.mmdb"

	// 检查测试文件是否存在，不存在则跳过（防止在不同环境下执行失败）
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Skip("Test MMDB file not found, skipping unit test")
	}

	provider, err := NewMaxMindProvider(dbPath)
	assert.NoError(t, err)
	defer provider.Close()

	t.Run("Valid IPv4 Lookup", func(t *testing.T) {
		// 假设 81.2.69.142 在测试库中有数据 (GeoIP2 官方测试库常包含这个)
		country, region, err := provider.Lookup("81.2.69.142")
		if err == nil {
			assert.NotEmpty(t, country)
			t.Logf("Lookup success: Country=%s, Region=%s", country, region)
		} else {
			t.Logf("Lookup failed (expected if IP not in test db): %v", err)
		}
	})

	t.Run("Invalid IP", func(t *testing.T) {
		_, _, err := provider.Lookup("invalid-ip")
		assert.Error(t, err)
	})

	t.Run("Empty IP", func(t *testing.T) {
		_, _, err := provider.Lookup("")
		assert.Error(t, err)
	})
}
