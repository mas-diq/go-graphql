package models

import (
	"mas-diq/go-graphql/config"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name  string `json:"name" gorm:"size:100;not null"`
	Email string `json:"email" gorm:"size:255;unique;not null"`
}

func (u *User) TableName() string {
	return "users"
}

func CreateUserData(m *User) (err error) {
	if err = config.DB.Create(m).Error; err != nil {
		return err
	}
	return nil
}

func UpdateUserData(m *User) (err error) {
	if err = config.DB.Save(m).Error; err != nil {
		return err
	}
	return nil
}

func GetListUser(m *[]User, filter User) (err error) {
	query := config.DB.Table("users")

	if filter.ID != 0 {
		query = query.Where("id = ?", filter.ID)
	}

	if filter.Name != "" {
		query = query.Where("name = ?", filter.Name)
	}

	if filter.Email != "" {
		query = query.Where("email = ?", filter.Email)
	}

	query = query.Find(m)
	return query.Error
}

func GetOneUser(m *User, id uint64) (err error) {
	query := config.DB.
		Table("users").
		Where("id = ?", id).
		First(m)
	return query.Error
}

func DeleteUser(m *User) (err error) {
	query := config.DB.Table("users").Delete(m)
	return query.Error
}
