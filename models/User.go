package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Username      string
	Email         string
	Password      string
	Image         string
	ImageBg       string
	Description   string
	Post          []Post          `gorm:"foreignKey:UserID"`
	Notifications []Notifications `gorm:"foreignKey:NotificationID"`
}
