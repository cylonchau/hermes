package model

type NSRecord struct {
	ID         int64  `gorm:"type:bigint;primaryKey;autoIncrement;comment:主键id;" json:"id"`
	RecordID   int64  `gorm:"type:bigint;unique;not null;comment:关联record表的id;" json:"record_id"`
	NameServer string `gorm:"type:varchar(255);not null;index;comment:名称服务器;" json:"name_server"` // 名称服务器 (a.iana-servers.net.)
	Remark     string `gorm:"type:text;comment:备注;" json:"remark"`                                // 备注
	IsGlue     bool   `gorm:"type:tinyint;default:false;comment:是否为胶水记录;" json:"is_glue"`         // 是否为胶水记录

	// 关联关系
	Record Record `gorm:"foreignKey:RecordID;constraint:OnDelete:CASCADE" json:"record,omitempty"`
}

func (NSRecord) TableName() string {
	return "record_ns"
}
