package services

import (
	"fmt"

	"github.com/iraqnroll/gochan/models"
)

type ThreadRepository interface {
	CreateNew(board_id int, topic string) (models.ThreadDto, error)
	GetById(thread_id int) (models.ThreadDto, error)
	GetAllByBoard(board_id int) ([]models.ThreadDto, error)
}

type IPostService interface {
	CreatePost(thread_id int, identifier, content, fingerprint string, is_op bool) (models.PostDto, error)
	GetThreadPosts(thread_id int) ([]models.PostDto, error)
	UpdateAttachedMedia(post_id int, attached_media, original_media string) error
}

type ThreadService struct {
	threadRepo  ThreadRepository
	postService IPostService
}

func NewThreadService(repo ThreadRepository, postService IPostService) *ThreadService {
	return &ThreadService{threadRepo: repo, postService: postService}
}

// Creates a new thread in the specified board with a OP post.
func (ts *ThreadService) CreateThread(board_id int, topic, identifier, content, fingerprint string) (*models.ThreadDto, error) {
	result, err := ts.threadRepo.CreateNew(board_id, topic)
	if err != nil {
		return nil, fmt.Errorf("ThreadService.CreateThread failed while creating a new thread : %w", err)
	}

	//New thread created, create and attach the OP post to it.
	post, err := ts.postService.CreatePost(result.Id, identifier, content, fingerprint, true)
	if err != nil {
		return nil, fmt.Errorf("ThreadService.CreateThread failed while creating a new thread : %w", err)
	}
	result.Posts = append(result.Posts, post)

	return &result, nil
}

// Gets a specified thread and it's content (posts)
func (ts *ThreadService) GetThread(thread_id int) (*models.ThreadDto, error) {
	result, err := ts.threadRepo.GetById(thread_id)
	if err != nil {
		return nil, fmt.Errorf("ThreadService.GetThread failed : %w", err)
	}

	//We found a valid active thread, fetch posts
	posts, err := ts.postService.GetThreadPosts(thread_id)
	if err != nil {
		return nil, fmt.Errorf("ThreadService.GetThread failed : %w", err)
	}

	result.Posts = posts

	return &result, nil
}

// Fetches threads and their respective content (posts) for a specified board.
// TODO: Refactor post fetching, right now we query the db for each thread when we could fetch all of it in one query (if its even needed).
func (ts *ThreadService) GetBoardThreads(board_id int) ([]models.ThreadDto, error) {
	result, err := ts.threadRepo.GetAllByBoard(board_id)
	if err != nil {
		return nil, fmt.Errorf("ThreadService.GetBoardThreads failed : %w", err)
	}

	for i, thread := range result {
		posts, err := ts.postService.GetThreadPosts(thread.Id)
		if err != nil {
			return nil, fmt.Errorf("ThreadService.GetBoardThreads failed : %w", err)
		}
		result[i].Posts = posts
	}

	return result, nil
}

// Sorts posts into their respective threads (Probably wont need this ever...)
func (ts *ThreadService) SortPostsIntoThreads(threads []models.ThreadDto, posts []models.PostDto) {
	postHashMap := make(map[int][]models.PostDto)

	for _, post := range posts {
		postHashMap[post.ThreadId] = append(postHashMap[post.ThreadId], post)
	}

	for i := range threads {
		threadId := threads[i].Id
		threads[i].Posts = postHashMap[threadId]
	}
}

func (ts *ThreadService) UpdateAttachedMedia(post_id int, attached_media, original_media string) error {
	return ts.postService.UpdateAttachedMedia(post_id, attached_media, original_media)
}
