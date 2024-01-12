package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model

	ImagePost   string
	Description string
	UserID      uint
	User        User `gorm:"foreignKey:UserID"`
	Likes       []Likes
	Comments    []Comments
}
