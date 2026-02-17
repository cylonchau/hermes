package model

type SRVRecord struct {
	ID       int64  `gorm:"type:bigint;primaryKey;autoIncrement;comment:主键id;" json:"id"`
	RecordID int64  `gorm:"type:bigint;unique;not null;comment:关联record表的id;" json:"record_id"`
	Priority uint16 `gorm:"type:smallint;not null;index;comment:优先级;" json:"priority"`    // 优先级
	Weight   uint16 `gorm:"type:smallint;not null;index;comment:权重;" json:"weight"`       // 权重
	Port     uint16 `gorm:"type:smallint;not null;index;comment:端口;" json:"port"`         // 端口
	Target   string `gorm:"type:varchar(255);not null;index;comment:目标主机;" json:"target"` // 目标主机
	Remark   string `gorm:"type:text;comment:备注;" json:"remark"`                          // 备注
	Service  string `gorm:"size:100;comment:服务类型;" json:"service"`                        // 服务类型
	Protocol string `gorm:"size:50;comment:协议类型;" json:"protocol"`                        // 协议类型

	// 关联关系
	Record Record `gorm:"foreignKey:RecordID;constraint:OnDelete:CASCADE" json:"record,omitempty"`

	// 不映射为表字段，用于 Join 查询结果映射，避免嵌套对象解析
	TTL uint32 `gorm:"->" json:"ttl"`
}

func (SRVRecord) TableName() string {
	return "record_srv"
}

func init() {
	RegisterModel(&SRVRecord{})
}
