package models

import "gorm.io/gorm"

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

func (p *Post) TableName() string {
	return "posts"
}
