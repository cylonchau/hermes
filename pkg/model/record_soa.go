package model

// SOA记录表 (Start of Authority)
type SOARecord struct {
	ID        int64  `gorm:"type:bigint;primaryKey;autoIncrement;comment:主键id;" json:"id"`
	RecordID  int64  `gorm:"type:bigint;unique;not null;comment:关联record表的id;" json:"record_id"`
	PrimaryNS string `gorm:"type:varchar(255);not null;comment:主名称服务器;" json:"primary_ns"`  // 主名称服务器 (sns.dns.icann.org.)
	MBox      string `gorm:"type:varchar(255);not null;comment:管理员邮箱;" json:"mail_box"`     // 管理员邮箱 (noc.dns.icann.org.)
	Serial    uint32 `gorm:"type:int;not null;default:0;comment:序列号;" json:"serial"`        // 序列号 (2017042745)
	Refresh   uint32 `gorm:"type:int;not null;default:7200;comment:刷新间隔;" json:"refresh"`   // 刷新间隔 (7200)
	Retry     uint32 `gorm:"type:int;not null;default:3600;comment:重试间隔;" json:"retry"`     // 重试间隔 (3600)
	Expire    uint32 `gorm:"type:int;not null;default:1209600;comment:过期时间;" json:"expire"` // 过期时间 (1209600)
	MinTTL    uint32 `gorm:"type:int;not null;default:3600;comment:最小TTL;" json:"minttl"`   // 最小TTL (3600)
	Remark    string `gorm:"type:text;comment:备注;" json:"remark"`                           // SOA记录备注

	// 关联关系
	Record Record `gorm:"foreignKey:RecordID;constraint:OnDelete:CASCADE" json:"record,omitempty"`
}

func (SOARecord) TableName() string {
	return "record_soa"
}
