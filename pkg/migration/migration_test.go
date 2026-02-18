package migration

import (
	"context"
	"testing"

	"github.com/cylonchau/hermes/pkg/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockStore is a mock implementation of store.Store
type MockStore struct {
	mock.Mock
}

func (m *MockStore) Initialize(config store.DatabaseConfig) error {
	return m.Called(config).Error(0)
}

func (m *MockStore) GetDB() *gorm.DB {
	return m.Called().Get(0).(*gorm.DB)
}

func (m *MockStore) Close() error {
	return m.Called().Error(0)
}

func (m *MockStore) HealthCheck() error {
	return m.Called().Error(0)
}

func (m *MockStore) AutoMigrate(models ...interface{}) error {
	// Match the variadic arguments as a slice
	return m.Called(models).Error(0)
}

func (m *MockStore) GetDatabaseType() store.DBType {
	return m.Called().Get(0).(store.DBType)
}

func (m *MockStore) IsInitialized() bool {
	return m.Called().Bool(0)
}

func (m *MockStore) MonitorConnectionPool(ctx context.Context) {
	m.Called(ctx)
}

func TestMigrate_Mock(t *testing.T) {
	mockStore := new(MockStore)

	// Inject the mock store
	store.ResetInstance(mockStore)
	defer store.ResetInstance(nil)

	// Expect AutoMigrate to be called with any slice of interfaces.
	mockStore.On("AutoMigrate", mock.Anything).Return(nil)

	// Perform Migration
	err := Migrate("mysql")
	assert.NoError(t, err)

	// Verify expectations
	mockStore.AssertExpectations(t)
}

func TestUpgrade_Mock(t *testing.T) {
	mockStore := new(MockStore)
	store.ResetInstance(mockStore)
	defer store.ResetInstance(nil)

	mockStore.On("AutoMigrate", mock.Anything).Return(nil)

	err := Upgrade("postgres")
	assert.NoError(t, err)
	mockStore.AssertExpectations(t)
}
