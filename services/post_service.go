package services

import (
	"fmt"

	"github.com/iraqnroll/gochan/models"
)

type PostRepository interface {
	GetAllByThread(thread_id int) ([]models.PostDto, error)
	CreateNew(thread_id int, identifier, content string, is_op bool) (models.PostDto, error)
}

type PostService struct {
	postRepo PostRepository
}

// Creates a post in a specific thread
func (ps *PostService) CreatePost(thread_id int, identifier, content string, is_op bool) (models.PostDto, error) {
	post, err := ps.postRepo.CreateNew(thread_id, identifier, content, false)
	if err != nil {
		return post, fmt.Errorf("PostService.CreateReply failed : %w", err)
	}

	return post, nil
}

// Fetches all posts of the specified thread
func (ps *PostService) GetThreadPosts(thread_id int) ([]models.PostDto, error) {
	posts, err := ps.postRepo.GetAllByThread(thread_id)
	if err != nil {
		return nil, fmt.Errorf("PostService.GetThreadPosts failed : %w", err)
	}

	return posts, nil
}
