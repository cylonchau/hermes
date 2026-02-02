package model

import (
	"time"
)

// CNAME记录表
type CNAMERecord struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	RecordID  int64     `gorm:"unique;not null" json:"record_id"`
	Target    string    `gorm:"not null;index" json:"target"` // 目标域名
	Remark    string    `gorm:"type:text" json:"remark"`      // CNAME记录备注
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联关系
	Record Record `gorm:"foreignKey:RecordID;constraint:OnDelete:CASCADE" json:"record,omitempty"`
}

func (CNAMERecord) TableName() string {
	return "record_cname"
}
