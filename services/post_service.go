package services

import (
	"fmt"

	"github.com/iraqnroll/gochan/models"
)

type PostRepository interface {
	GetAllByThread(thread_id int) ([]models.PostDto, error)
	CreateNew(thread_id int, identifier, content, fingerprint string, is_op bool) (models.PostDto, error)
	GetMostRecent(num_of_posts int) ([]models.RecentPostsDto, error)
	UpdateAttachedMedia(post_id int, attached_media, original_media string) error
}

type PostService struct {
	PostRepo PostRepository
}

func NewPostService(repo PostRepository) *PostService {
	return &PostService{PostRepo: repo}
}

// Creates a post in a specific thread
func (ps *PostService) CreatePost(thread_id int, identifier, content, fingerprint string, is_op bool) (models.PostDto, error) {
	post, err := ps.PostRepo.CreateNew(thread_id, identifier, content, fingerprint, is_op)
	if err != nil {
		return post, fmt.Errorf("PostService.CreateReply failed : %w", err)
	}

	return post, nil
}

// Fetches all posts of the specified thread
func (ps *PostService) GetThreadPosts(thread_id int) ([]models.PostDto, error) {
	posts, err := ps.PostRepo.GetAllByThread(thread_id)
	if err != nil {
		return nil, fmt.Errorf("PostService.GetThreadPosts failed : %w", err)
	}

	return posts, nil
}

func (ps *PostService) GetMostRecent(num_of_posts int) ([]models.RecentPostsDto, error) {
	posts, err := ps.PostRepo.GetMostRecent(num_of_posts)
	if err != nil {
		return nil, fmt.Errorf("PostService.GetMostRecent failed : %s", err)
	}

	return posts, nil
}

// TODO: Maybe split this into a separate service ? right now thread service wraps this
func (ps *PostService) UpdateAttachedMedia(post_id int, attached_media, original_media string) error {
	return ps.PostRepo.UpdateAttachedMedia(post_id, attached_media, original_media)
}
