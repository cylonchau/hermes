package model

// CNAME记录表
type CNAMERecord struct {
	ID       int64  `gorm:"type:bigint;primaryKey;autoIncrement;comment:主键id;" json:"id"`
	RecordID int64  `gorm:"type:bigint;unique;not null;comment:关联record表的id;" json:"record_id"`
	Target   string `gorm:"type:varchar(255);not null;index;comment:目标域名;" json:"target"` // 目标域名
	Remark   string `gorm:"type:text;comment:CNAME记录备注;" json:"remark"`                   // CNAME记录备注

	// 关联关系
	Record Record `gorm:"foreignKey:RecordID;constraint:OnDelete:CASCADE" json:"record,omitempty"`
}

func (CNAMERecord) TableName() string {
	return "record_cname"
}
