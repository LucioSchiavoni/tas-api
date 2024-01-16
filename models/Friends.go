package models

import "gorm.io/gorm"

type Friends struct {
	gorm.Model

	UserID    uint
	User      User `gorm:"foreignKey:UserID"`
	CreatorID uint
	Creator   User `gorm:"foreignKey:CreatorID"`
}
