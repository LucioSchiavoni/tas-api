package models

import "gorm.io/gorm"

type Comments struct {
	gorm.Model

	UserID  uint `gorm:"foreignKey:UserID"`
	PostID  uint `gorm:"foreignKey:PostID"`
	Content string
	User    User `gorm:"foreignKey:UserID"`
	Post    Post `gorm:"foreignKey:PostID"`
}
