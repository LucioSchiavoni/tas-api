package models

import "gorm.io/gorm"

type Likes struct {
	gorm.Model

	UserID uint
	PostID uint
	User   User `gorm:"foreignKey:UserID"`
	Post   Post `gorm:"foreignKey:PostID"`
}
