package models

import "gorm.io/gorm"

type Notifications struct {
	gorm.Model

	UserID   uint `gorm:"foreignKey:UserID"`
	User     User `gorm:"foreignKey:UserID"`
	Type     string
	Post     Post `gorm:"foreignKey:PostID"`
	PostID   uint `gorm:"foreignKey:PostID"`
	Check    bool
	UserFrom uint
}
