package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name string `gorm:"unique"`
	Pass string `gorm:"not null"`
}

func (User) TableName() string {
	return "backend.user"
}
