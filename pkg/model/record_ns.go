package model

import (
	"time"
)

type NSRecord struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	RecordID   int64     `gorm:"unique;not null" json:"record_id"`
	NameServer string    `gorm:"not null;index" json:"name_server"` // 名称服务器 (a.iana-servers.net.)
	Remark     string    `gorm:"type:text" json:"remark"`           // 备注
	IsGlue     bool      `gorm:"default:false" json:"is_glue"`      // 是否为胶水记录
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联关系
	Record Record `gorm:"foreignKey:RecordID;constraint:OnDelete:CASCADE" json:"record,omitempty"`
}

func (NSRecord) TableName() string {
	return "record_ns"
}
