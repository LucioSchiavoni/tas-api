package models

import "gorm.io/gorm"

type Comments struct {
	gorm.Model

	UserID    uint `gorm:"foreignKey:UserID"`
	PostID    uint `gorm:"foreignKey:PostID"`
	CreatorID uint
	Content   string
	User      User `gorm:"foreignKey:UserID"`
	Post      Post `gorm:"foreignKey:PostID"`
	Creator   User `gorm:"foreignKey:CreatorID"`
}
