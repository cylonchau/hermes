package model

import "gorm.io/gorm"

type ARecord struct {
	gorm.Model
	RecordID uint   `gorm:"unique;not null" json:"record_id"`
	IP       string `gorm:"not null;index" json:"ip"` // IPv4地址
	Remark   string `gorm:"type:text" json:"remark"`  // 备注

	// 关联关系
	Record Record `gorm:"foreignKey:record_id;constraint:OnDelete:CASCADE" json:"record,omitempty"`
}
