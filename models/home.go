package models

type HomeViewModel struct {
	Boards      []BoardDto
	RecentPosts []RecentPostsDto
	ErrMsg      string
}

func NewHomeViewModel(boards []BoardDto, posts []RecentPostsDto, err error) (m HomeViewModel) {
	if len(boards) <= 0 || len(posts) <= 0 {
		m.ErrMsg = "Failed to retrieve homepage data..."
	} else {
		m.Boards = boards
		m.RecentPosts = posts
	}

	return m
}
