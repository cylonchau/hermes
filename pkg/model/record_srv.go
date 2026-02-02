package model

import (
	"time"
)

type SRVRecord struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	RecordID  int64     `gorm:"unique;not null" json:"record_id"`
	Priority  uint16    `gorm:"not null;index" json:"priority"` // 优先级
	Weight    uint16    `gorm:"not null" json:"weight"`         // 权重
	Port      uint16    `gorm:"not null;index" json:"port"`     // 端口
	Target    string    `gorm:"not null;index" json:"target"`   // 目标主机
	Remark    string    `gorm:"type:text" json:"remark"`        // 备注
	Service   string    `gorm:"size:100" json:"service"`        // 服务类型
	Protocol  string    `gorm:"size:50" json:"protocol"`        // 协议类型
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联关系
	Record Record `gorm:"foreignKey:RecordID;constraint:OnDelete:CASCADE" json:"record,omitempty"`
}

func (SRVRecord) TableName() string {
	return "record_srv"
}
