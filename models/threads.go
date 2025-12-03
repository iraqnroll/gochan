package models

type ThreadDto struct {
	Id       int       `json:"id"`
	Posts    []PostDto `json:"posts"`
	Topic    string    `json:"topic"`
	Locked   bool      `json:"locked"`
	BoardId  int       `json:"board_id"`
	BoardUri string    `json:"board_uri"`
}

type ThreadViewModel struct {
	Id           int
	Topic        string
	BannerUrl    string
	OPPost       PostDto
	Replies      []PostDto
	PostsPerPage int
	BoardUri     string
	ErrMsg       string
}

func NewThreadsViewModel(id, postsPerPage int, bannerUrl, board_uri, topic string, op_post PostDto, replies []PostDto) (t ThreadViewModel) {
	t.Id = id
	t.BannerUrl = bannerUrl
	t.OPPost = op_post
	t.Replies = replies
	t.PostsPerPage = postsPerPage
	t.BoardUri = board_uri
	t.Topic = topic

	return t
}
