package model

import "gorm.io/gorm"

// SOA记录表 (Start of Authority)
type SOARecord struct {
	gorm.Model
	RecordID  uint   `gorm:"unique;not null" json:"record_id"`
	PrimaryNS string `gorm:"not null" json:"primary_ns"` // 主名称服务器 (sns.dns.icann.org.)
	MBox      string `gorm:"not null" json:"mail_box"`   // 管理员邮箱 (noc.dns.icann.org.)
	Serial    uint32 `gorm:"not null" json:"serial"`     // 序列号 (2017042745)
	Refresh   uint32 `gorm:"not null" json:"refresh"`    // 刷新间隔 (7200)
	Retry     uint32 `gorm:"not null" json:"retry"`      // 重试间隔 (3600)
	Expire    uint32 `gorm:"not null" json:"expire"`     // 过期时间 (1209600)
	MinTTL    uint32 `gorm:"not null" json:"minttl"`     // 最小TTL (3600)
	Remark    string `gorm:"type:text" json:"remark"`    // SOA记录备注

	// 关联关系
	Record Record `gorm:"foreignKey:record_id;constraint:OnDelete:CASCADE" json:"record,omitempty"`
}
