package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserId   int64  `gorm:"primaryKey;autoIncrement" json:"userId"`
	Name     string `json:"name"`
	Email    string `gorm:"unique;not null" json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
