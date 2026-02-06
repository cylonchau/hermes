package model

import (
	"fmt"

	"github.com/cylonchau/hermes/pkg/store"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	Models []interface{}
)

// RegisterModel 注册数据库模型
func RegisterModel(m interface{}) {
	Models = append(Models, m)
}

// InitDB 初始化模型层使用的数据库实例
func InitDB(driver string) error {
	s := store.GetInstance()
	if !s.IsInitialized() {
		// 如果尚未初始化，尝试根据驱动类型从配置中查找并初始化
		// 这里假设配置已经加载
		return fmt.Errorf("database store not initialized")
	}

	DB = s.GetDB()
	return nil
}
