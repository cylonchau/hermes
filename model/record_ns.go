package model

import "gorm.io/gorm"

type NSRecord struct {
	gorm.Model
	RecordID   uint   `gorm:"unique;not null" json:"record_id"`
	NameServer string `gorm:"not null;index" json:"name_server"` // 名称服务器 (a.iana-servers.net.)
	Remark     string `gorm:"type:text" json:"remark"`           // 备注
	IsGlue     bool   `gorm:"default:false" json:"is_glue"`      // 是否为胶水记录

	// 关联关系
	Record Record `gorm:"foreignKey:record_id;constraint:OnDelete:CASCADE" json:"record,omitempty"`
}
