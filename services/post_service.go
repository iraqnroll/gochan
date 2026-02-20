package services

import (
	"fmt"

	"github.com/iraqnroll/gochan/models"
	"github.com/iraqnroll/gochan/rand"
)

type PostRepository interface {
	GetAllByThread(thread_id int, for_mod bool) ([]models.PostDto, error)
	CreateNew(thread_id int, identifier, content, fingerprint string, is_op bool) (models.PostDto, error)
	GetMostRecent(num_of_posts int) ([]models.RecentPostsDto, error)
	UpdateAttachedMedia(post_id int, attached_media, original_media string) error
	SoftDeletePost(post_id int) error
}

type PostService struct {
	PostRepo        PostRepository
	FingerprintSalt string
}

func NewPostService(repo PostRepository, fprintSalt string) *PostService {
	return &PostService{PostRepo: repo, FingerprintSalt: fprintSalt}
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
func (ps *PostService) GetThreadPosts(thread_id int, for_mod bool) ([]models.PostDto, error) {
	posts, err := ps.PostRepo.GetAllByThread(thread_id, for_mod)
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

// TODO: Deal with attached post media, if we hide the post, we want to hide the attached content from /static as well
func (ps *PostService) SoftDeletePost(post_id int) error {
	return ps.PostRepo.SoftDeletePost(post_id)
}

// TODO: Maybe split this into a separate service ? right now thread service wraps this
func (ps *PostService) UpdateAttachedMedia(post_id int, attached_media, original_media string) error {
	return ps.PostRepo.UpdateAttachedMedia(post_id, attached_media, original_media)
}

// TODO: Thread service wraps this too.... either my board handler is structured wrong or i need a separate service...
func (ps *PostService) GenerateFingerprint(ip string) string {
	return rand.GenerateFingerprint(ip, ps.FingerprintSalt)
}
