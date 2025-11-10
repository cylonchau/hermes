package model

import (
	"time"

	"gorm.io/gorm"
)

type Zone struct {
	gorm.Model
	Name        string    `gorm:"unique;not null;index" json:"name"` // example.org.
	Serial      uint32    `json:"serial"`
	Description string    `gorm:"type:text" json:"description"`
	Remark      string    `gorm:"type:text" json:"remark"` // 备注信息
	Contact     string    `gorm:"size:255" json:"contact"` // 联系人
	Email       string    `gorm:"size:255 s" json:"email"` // 联系邮箱
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedBy   string    `gorm:"size:100" json:"created_by"` // 创建人
	UpdatedBy   string    `gorm:"size:100" json:"updated_by"` // 更新人
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 关联关系
	Records []Record `gorm:"foreignKey:ZoneID;constraint:OnDelete:CASCADE" json:"records,omitempty"`
}
