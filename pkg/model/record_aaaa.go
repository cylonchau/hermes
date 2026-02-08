package model

type AAAARecord struct {
	ID       int64  `gorm:"type:bigint;primaryKey;autoIncrement;comment:主键id;" json:"id"`
	RecordID int64  `gorm:"type:bigint;unique;not null;comment:关联record表的id;" json:"record_id"` // 关联的record_id
	IP       []byte `gorm:"type:BINARY(16);not null;index;comment:IPv6地址;" json:"ip"`           // IPv6地址
	Remark   string `gorm:"type:varchar(256);comment:备注;" json:"remark"`                        // 备注

	// 关联关系
	Record Record `gorm:"foreignKey:RecordID;constraint:OnDelete:CASCADE" json:"record,omitempty"`

	// 不映射为表字段，用于 Join 查询结果映射，避免嵌套对象解析
	TTL uint32 `gorm:"->" json:"ttl"`
}

func (AAAARecord) TableName() string {
	return "record_aaaa"
}

func init() {
	RegisterModel(&AAAARecord{})
}
