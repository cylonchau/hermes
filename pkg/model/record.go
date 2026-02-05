package model

type Record struct {
	ID       int64  `gorm:"type:bigint;primaryKey;autoIncrement;comment:主键id;" json:"id"`
	ZoneID   int64  `gorm:"type:bigint;not null;index;comment:关联zone表的id;" json:"zone_id"`
	Name     string `gorm:"type:varchar(255);not null;index;comment:record记录;" json:"name"` // record记录
	Type     string `gorm:"type:varchar(50);not null;index;comment:记录类型;" json:"type"`      // 记录类型
	TTL      uint32 `gorm:"type:int;not null;comment:这st条记录缓存时间为1小时（单位秒）;" json:"ttl"`      // 这条记录缓存时间为1小时（单位秒）
	Remark   string `gorm:"type:text;comment:备注;" json:"remark"`                            // 备注
	Tags     string `gorm:"size:500;comment:标签，用逗号分隔;" json:"tags"`                         // 标签，用逗号分隔
	Source   string `gorm:"size:100;comment:数据来源;" json:"source"`                           // 数据来源
	IsActive bool   `gorm:"default:true;comment:该记录是否活跃;" json:"is_active"`                 // 该记录是否活跃

	// 关联关系
	// 这里全部使用了指针类型，对序列化更友好
	Zone        *Zone        `gorm:"foreignKey:ZoneID" json:"zone,omitempty"`
	ARecord     *ARecord     `gorm:"foreignKey:RecordID" json:"a_record,omitempty"`
	AAAARecord  *AAAARecord  `gorm:"foreignKey:RecordID" json:"aaaa_record,omitempty"`
	CNAMERecord *CNAMERecord `gorm:"foreignKey:RecordID" json:"cname_record,omitempty"`
	MXRecord    *MXRecord    `gorm:"foreignKey:RecordID" json:"mx_record,omitempty"`
	TXTRecord   *TXTRecord   `gorm:"foreignKey:RecordID" json:"txt_record,omitempty"`
	SRVRecord   *SRVRecord   `gorm:"foreignKey:RecordID" json:"srv_record,omitempty"`
	SOARecord   *SOARecord   `gorm:"foreignKey:RecordID" json:"soa_record,omitempty"`
	NSRecord    *NSRecord    `gorm:"foreignKey:RecordID" json:"ns_record,omitempty"`
}

func (Record) TableName() string {
	return "record"
}
