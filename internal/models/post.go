package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Post struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key" example:"123e4567-e89b-12d3-a456-426614174000"`
	Title     string    `json:"title" gorm:"not null" example:"Sample Post Title"`
	Content   string    `json:"content" gorm:"not null" example:"This is the content of the post"`
	Author    string    `json:"author" gorm:"not null" example:"John Doe"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime" example:"2023-01-01T00:00:00Z"`
}

func (p *Post) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

type CreatePostRequest struct {
	Title   string `json:"title" binding:"required" example:"New Post Title"`
	Content string `json:"content" binding:"required" example:"Post content goes here"`
	Author  string `json:"author" binding:"required" example:"Jane Doe"`
}

type UpdatePostRequest struct {
	Title   string `json:"title,omitempty" example:"Updated Post Title"`
	Content string `json:"content,omitempty" example:"Updated post content"`
	Author  string `json:"author,omitempty" example:"Updated Author"`
}

func NewPost(req CreatePostRequest) *Post {
	now := time.Now()
	return &Post{
		ID:        uuid.New().String(),
		Title:     req.Title,
		Content:   req.Content,
		Author:    req.Author,
		CreatedAt: now,
		UpdatedAt: now,
	}
}