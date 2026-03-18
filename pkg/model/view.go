package model

import (
	"time"
)

type View struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"not null;uniqueIndex" json:"name"`
	Category  string    `gorm:"type:varchar(20);not null;default:'acl';comment:类型: acl 或 geoip" json:"category"`
	Value     string    `gorm:"type:text;comment:匹配规则内容(CIDR列表或GeoIP标签)" json:"value"`
	Priority  int       `gorm:"type:int;not null;default:0;index;comment:匹配优先级(越小越优先)" json:"priority"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (View) TableName() string {
	return "view"
}

func init() {
	RegisterModel(&View{})
}
