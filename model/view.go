package model

import "gorm.io/gorm"

type View struct {
	gorm.Model
	Name string `json:"name"`
}
