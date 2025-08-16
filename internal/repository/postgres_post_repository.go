package repository

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"postService/internal/models"
)

type PostgresPostRepository struct {
	db *gorm.DB
}

func NewPostgresPostRepository(db *gorm.DB) PostRepository {
	return &PostgresPostRepository{db: db}
}

func (r *PostgresPostRepository) Create(post *models.Post) error {
	if err := r.db.Create(post).Error; err != nil {
		return err
	}
	return nil
}

func (r *PostgresPostRepository) GetByID(id string) (*models.Post, error) {
	var post models.Post
	if err := r.db.First(&post, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}
	return &post, nil
}

func (r *PostgresPostRepository) GetAll() ([]*models.Post, error) {
	var posts []*models.Post
	if err := r.db.Order("created_at DESC").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *PostgresPostRepository) Update(id string, req models.UpdatePostRequest) (*models.Post, error) {
	var post models.Post
	if err := r.db.First(&post, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("post not found")
		}
		return nil, err
	}

	// Update fields if provided
	updateData := make(map[string]interface{})
	if req.Title != "" {
		updateData["title"] = req.Title
	}
	if req.Content != "" {
		updateData["content"] = req.Content
	}
	if req.Author != "" {
		updateData["author"] = req.Author
	}
	updateData["updated_at"] = time.Now()

	if err := r.db.Model(&post).Updates(updateData).Error; err != nil {
		return nil, err
	}

	// Fetch updated record
	if err := r.db.First(&post, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostgresPostRepository) Delete(id string) error {
	result := r.db.Delete(&models.Post{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("post not found")
	}
	return nil
}