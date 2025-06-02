package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name  string `json:"name" gorm:"size:100;not null"`
	Email string `json:"email" gorm:"size:255;unique;not null"`
}

func (u *User) TableName() string {
	return "users"
}
