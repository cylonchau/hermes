package model

// MX记录表
type MXRecord struct {
	ID       int64  `gorm:"type:bigint;primaryKey;autoIncrement;comment:主键id;" json:"id"`
	RecordID int64  `gorm:"type:bigint;unique;not null;comment:关联record表的id;" json:"record_id"`
	Host     string `gorm:"type:varchar(255);not null;index;comment:邮件服务器域名;" json:"host"` // 邮件服务器域名
	Priority uint16 `gorm:"type:smallint;not null;index;comment:优先级;" json:"priority"`     // 优先级
	Remark   string `gorm:"type:text;comment:备注;" json:"remark"`                           // 备注
	Provider string `gorm:"type:varchar(100);comment:邮件服务提供商;" json:"provider"`            // 邮件服务提供商

	// 关联关系
	Record Record `gorm:"foreignKey:RecordID;constraint:OnDelete:CASCADE" json:"record,omitempty"`

	// 不映射为表字段，用于 Join 查询结果映射，避免嵌套对象解析
	TTL uint32 `gorm:"->" json:"ttl"`
}

func (MXRecord) TableName() string {
	return "record_mx"
}

func init() {
	RegisterModel(&MXRecord{})
}
