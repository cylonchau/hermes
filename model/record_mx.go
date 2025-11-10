package model

import "gorm.io/gorm"

// MX记录表
type MXRecord struct {
	gorm.Model
	RecordID uint   `gorm:"unique;not null" json:"record_id"`
	Host     string `gorm:"not null;index" json:"host"`     // 邮件服务器域名
	Priority uint16 `gorm:"not null;index" json:"priority"` // 优先级
	Remark   string `gorm:"type:text" json:"remark"`        // 备注
	Provider string `gorm:"size:100" json:"provider"`       // 邮件服务提供商

	// 关联关系
	Record Record `gorm:"foreignKey:record_id;constraint:OnDelete:CASCADE" json:"record,omitempty"`
}
