package service

import (
	"postService/internal/models"
	"postService/internal/repository"
)

type PostService interface {
	CreatePost(req models.CreatePostRequest) (*models.Post, error)
	GetPost(id string) (*models.Post, error)
	GetAllPosts() ([]*models.Post, error)
	UpdatePost(id string, req models.UpdatePostRequest) (*models.Post, error)
	DeletePost(id string) error
}

type postService struct {
	repo repository.PostRepository
}

func NewPostService(repo repository.PostRepository) PostService {
	return &postService{repo: repo}
}

func (s *postService) CreatePost(req models.CreatePostRequest) (*models.Post, error) {
	post := models.NewPost(req)
	err := s.repo.Create(post)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *postService) GetPost(id string) (*models.Post, error) {
	return s.repo.GetByID(id)
}

func (s *postService) GetAllPosts() ([]*models.Post, error) {
	return s.repo.GetAll()
}

func (s *postService) UpdatePost(id string, req models.UpdatePostRequest) (*models.Post, error) {
	return s.repo.Update(id, req)
}

func (s *postService) DeletePost(id string) error {
	return s.repo.Delete(id)
}