package model

// TXT记录表
type TXTRecord struct {
	ID       int64  `gorm:"type:bigint;primaryKey;autoIncrement;comment:主键id;" json:"id"`
	RecordID int64  `gorm:"type:bigint;unique;not null;comment:关联record表的id;" json:"record_id"`
	Text     string `gorm:"type:text;not null;comment:内容;" json:"text"` // 内容
	Remark   string `gorm:"type:text;comment:备注;" json:"remark"`        // 备注
	Purpose  string `gorm:"size:100;comment:用途说明;" json:"purpose"`      // 用途说明 (SPF, DKIM, 验证等)

	// 关联关系
	Record Record `gorm:"foreignKey:RecordID;constraint:OnDelete:CASCADE" json:"record,omitempty"`
}

func (TXTRecord) TableName() string {
	return "record_txt"
}

func init() {
	RegisterModel(&TXTRecord{})
}
