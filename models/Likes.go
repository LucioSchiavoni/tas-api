package models

import "gorm.io/gorm"

type Likes struct {
	gorm.Model

	UserID    uint
	PostID    uint
	CreatorID uint
	User      User `gorm:"foreignKey:UserID"`
	Post      Post `gorm:"foreignKey:PostID"`
	Creator   User `gorm:"foreignKey:CreatorID"`
}
