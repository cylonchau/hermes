package model

import "gorm.io/gorm"

// CNAME记录表
type CNAMERecord struct {
	gorm.Model
	RecordID uint   `gorm:"unique;not null" json:"record_id"`
	Target   string `gorm:"not null;index" json:"target"` // 目标域名
	Remark   string `gorm:"type:text" json:"remark"`      // CNAME记录备注

	// 关联关系
	Record Record `gorm:"foreignKey:record_id;constraint:OnDelete:CASCADE" json:"record,omitempty"`
}
