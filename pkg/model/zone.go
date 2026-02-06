package model

type Zone struct {
	ID          int64  `gorm:"type:bigint;primaryKey;autoIncrement;comment:主键id;" json:"id"`
	Name        string `gorm:"type:varchar(255);unique;not null;index;comment:zone名称;" json:"name"` // example.org.
	Serial      uint32 `gorm:"type:int;not null;default:0;comment:序列号;" json:"serial"`              // SOA记录的序列号
	Description string `gorm:"type:text;comment:描述信息;" json:"description"`                          // 描述信息
	Remark      string `gorm:"type:text;comment:备注信息;" json:"remark"`                               // 备注信息
	Contact     string `gorm:"type:varchar(255);comment:联系人;" json:"contact"`                       // 联系人
	Email       string `gorm:"type:varchar(255);comment:联系邮箱;" json:"email"`                        // 联系邮箱
	IsActive    bool   `gorm:"default:true;comment:该zone是否活跃;" json:"is_active"`                    // 该zone是否活跃

	// 关联关系
	Records []Record `gorm:"foreignKey:ZoneID;constraint:OnDelete:CASCADE" json:"records,omitempty"`
}

func (Zone) TableName() string {
	return "zone"
}

func init() {
	RegisterModel(&Zone{})
}
