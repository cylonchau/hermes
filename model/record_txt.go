package model

import "gorm.io/gorm"

// TXT记录表
type TXTRecord struct {
	gorm.Model
	RecordID uint   `gorm:"unique;not null" json:"record_id"`
	Text     string `gorm:"type:text;not null" json:"text"` // 内容
	Remark   string `gorm:"type:text" json:"remark"`        // 备注
	Purpose  string `gorm:"size:100" json:"purpose"`        // 用途说明 (SPF, DKIM, 验证等)

	// 关联关系
	Record Record `gorm:"foreignKey:RecordID;constraint:OnDelete:CASCADE" json:"record,omitempty"`
}
