package repository

import (
	"errors"
	"sync"
	"time"

	"postService/internal/models"
)

type PostRepository interface {
	Create(post *models.Post) error
	GetByID(id string) (*models.Post, error)
	GetAll() ([]*models.Post, error)
	Update(id string, req models.UpdatePostRequest) (*models.Post, error)
	Delete(id string) error
}

type InMemoryPostRepository struct {
	posts map[string]*models.Post
	mutex sync.RWMutex
}

func NewInMemoryPostRepository() *InMemoryPostRepository {
	return &InMemoryPostRepository{
		posts: make(map[string]*models.Post),
	}
}

func (r *InMemoryPostRepository) Create(post *models.Post) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	r.posts[post.ID] = post
	return nil
}

func (r *InMemoryPostRepository) GetByID(id string) (*models.Post, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	post, exists := r.posts[id]
	if !exists {
		return nil, errors.New("post not found")
	}
	return post, nil
}

func (r *InMemoryPostRepository) GetAll() ([]*models.Post, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	posts := make([]*models.Post, 0, len(r.posts))
	for _, post := range r.posts {
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *InMemoryPostRepository) Update(id string, req models.UpdatePostRequest) (*models.Post, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	post, exists := r.posts[id]
	if !exists {
		return nil, errors.New("post not found")
	}
	
	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}
	if req.Author != "" {
		post.Author = req.Author
	}
	post.UpdatedAt = time.Now()
	
	r.posts[id] = post
	return post, nil
}

func (r *InMemoryPostRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	_, exists := r.posts[id]
	if !exists {
		return errors.New("post not found")
	}
	
	delete(r.posts, id)
	return nil
}