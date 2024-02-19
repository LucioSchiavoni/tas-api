package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	UserID  uint
	User    User `gorm:"foreignKey:UserID"`
	Content string
}
