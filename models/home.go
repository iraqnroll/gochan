package models

type HomeViewModel struct {
	Boards      []BoardDto
	RecentPosts []RecentPostsDto
	ErrMsg      string
}

func NewHomeViewModel(boards []BoardDto, posts []RecentPostsDto, err error) (m HomeViewModel) {
	m.Boards = boards
	m.RecentPosts = posts
	return m
}
