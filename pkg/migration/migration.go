package migration

import (
	"fmt"

	"github.com/cylonchau/hermes/pkg/model"
	"github.com/cylonchau/hermes/pkg/store"
)

// Migrate 执行数据库初始化和表创建
func Migrate(driver string) error {
	fmt.Printf("Starting migration for driver: %s\n", driver)
	s := store.GetInstance()
	return s.AutoMigrate(model.Models...)
}

// Upgrade 执行数据库架构升级
func Upgrade(driver string) error {
	fmt.Printf("Starting upgrade for driver: %s\n", driver)
	// TODO: 实现具体的升级逻辑
	return Migrate(driver)
}
