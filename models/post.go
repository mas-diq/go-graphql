package models

import (
	"mas-diq/go-graphql/config"

	"gorm.io/gorm"
)

type PostStatus string

const (
	Draft     PostStatus = "draft"
	Published PostStatus = "published"
	Archived  PostStatus = "archived"
)

type Post struct {
	gorm.Model
	Title     string     `json:"title" gorm:"size:255;not null"`
	Subtitle  string     `json:"subtitle" gorm:"size:255"`
	Image     string     `json:"image" gorm:"size:255"` // Store image URL or path
	Content   string     `json:"content" gorm:"type:text;not null"`
	Status    PostStatus `json:"status" gorm:"type:enum('draft','published','archived');default:'draft'"`
	CreatedBy uint       `json:"createdBy" gorm:"not null"`        // Foreign key to User
	User      User       `json:"user" gorm:"foreignKey:CreatedBy"` // Relationship
}

const postTable = "posts"

func (p *Post) TableName() string {
	return postTable
}

func CreatePostData(m *Post) (err error) {
	if err = config.DB.Create(m).Error; err != nil {
		return err
	}
	return nil
}

func UpdatePostData(m *Post) (err error) {
	if err = config.DB.Save(m).Error; err != nil {
		return err
	}
	return nil
}

func GetListPost(m *[]Post, filter Post) (err error) {
	query := config.DB.Table(postTable)

	if filter.ID != 0 {
		query = query.Where("id = ?", filter.ID)
	}

	if filter.Title != "" {
		query = query.Where("title = ?", filter.Title)
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	query = query.Find(m)
	return query.Error
}

func GetOnePost(m *Post, id uint64) (err error) {
	query := config.DB.
		Table(postTable).
		Where("id = ?", id).
		First(m)
	return query.Error
}

func DeletePost(m *Post) (err error) {
	query := config.DB.Table(postTable).Delete(m)
	return query.Error
}
