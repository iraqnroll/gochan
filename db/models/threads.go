package models

type ThreadDto struct {
	Id       int       `db:"id" json:"id" schema:"id"`
	Posts    []PostDto `json:"posts"`
	Topic    string    `db:"topic" json:"topic" schema:"topic"`
	Locked   bool      `db:"locked" json:"locked" schema:"locked"`
	BoardId  int       `db:"board_id" json:"board_id" schema:"boardId"`
	BoardUri string    `json:"board_uri" schema:"boardUri"`
	Pinned   bool      `db:"sticky"`
}

type ThreadViewModel struct {
	Id           int
	Topic        string
	BannerUrl    string
	OPPost       PostDto
	Replies      []PostDto
	PostsPerPage int
	BoardUri     string
	BoardView    bool
	ErrMsg       string
	Pinned       bool
}

func NewThreadsViewModel(id, postsPerPage int, bannerUrl, board_uri, topic string, op_post PostDto, replies []PostDto, boardview, pinned bool) (t ThreadViewModel) {
	t.Id = id
	t.BannerUrl = bannerUrl
	t.OPPost = op_post
	t.Replies = replies
	t.PostsPerPage = postsPerPage
	t.BoardUri = board_uri
	t.Topic = topic
	t.BoardView = boardview
	t.Pinned = pinned

	return t
}
