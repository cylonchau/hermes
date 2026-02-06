package model

// CAA (Certificate Authority Authorization)
type CAARecord struct {
	ID       int64  `gorm:"type:bigint;primaryKey;autoIncrement;comment:主键id;" json:"id"`
	RecordID int64  `gorm:"type:bigint;unique;not null;comment:关联record表的id;" json:"record_id"`
	Flag     uint8  `gorm:"type:tinyint;not null;comment:CAA标志位(0-255);" json:"flag"`
	Tag      string `gorm:"type:varchar(64);not null;comment:标签;" json:"tag"`
	Value    string `gorm:"type:varchar(256);not null;comment:值;" json:"value"`

	// 关联关系
	Record Record `gorm:"foreignKey:RecordID;constraint:OnDelete:CASCADE" json:"record,omitempty"`
}

func (CAARecord) TableName() string {
	return "record_caa"
}

func init() {
	RegisterModel(&CAARecord{})
}
