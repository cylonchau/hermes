package model

import "gorm.io/gorm"

type Record struct {
	gorm.Model
	ZoneID    uint   `gorm:"not null;index" json:"zone_id"`
	Name      string `gorm:"not null;index" json:"name"`    // record记录
	TTL       uint32 `gorm:"not null" json:"ttl"`           // 这条记录缓存时间为1小时（单位秒）
	Remark    string `gorm:"type:text" json:"remark"`       // 备注
	Tags      string `gorm:"size:500" json:"tags"`          // 标签，用逗号分隔
	Source    string `gorm:"size:100" json:"source"`        // 数据来源
	IsActive  bool   `gorm:"default:true" json:"is_active"` // 该记录是否活跃
	CreatedBy string `gorm:"size:100" json:"created_by"`    // 创建人
	UpdatedBy string `gorm:"size:100" json:"updated_by"`    // 更新人

	// 关联关系
	Zone        Zone         `gorm:"foreignKey:ZoneID" json:"zone,omitempty"`
	ARecord     *ARecord     `gorm:"foreignKey:RecordID" json:"a_record,omitempty"`
	AAAARecord  *AAAARecord  `gorm:"foreignKey:RecordID" json:"aaaa_record,omitempty"`
	CNAMERecord *CNAMERecord `gorm:"foreignKey:RecordID" json:"cname_record,omitempty"`
	MXRecord    *MXRecord    `gorm:"foreignKey:RecordID" json:"mx_record,omitempty"`
	TXTRecord   *TXTRecord   `gorm:"foreignKey:RecordID" json:"txt_record,omitempty"`
	SRVRecord   *SRVRecord   `gorm:"foreignKey:RecordID" json:"srv_record,omitempty"`
	SOARecord   *SOARecord   `gorm:"foreignKey:RecordID" json:"soa_record,omitempty"`
	NSRecord    *NSRecord    `gorm:"foreignKey:RecordID" json:"ns_record,omitempty"`
}
