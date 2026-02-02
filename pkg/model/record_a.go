package model

import (
	"time"
)

type ARecord struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	RecordID  int64     `gorm:"unique;not null" json:"record_id"` // 关联的record_id
	IP        string    `gorm:"not null;index" json:"ip"`         // IPv4地址
	Remark    string    `gorm:"type:text" json:"remark"`          // 备注
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联关系
	Record Record `gorm:"foreignKey:RecordID;constraint:OnDelete:CASCADE" json:"record,omitempty"`
}

func (ARecord) TableName() string {
	return "record_a"
}
