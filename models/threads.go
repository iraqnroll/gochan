package models

import (
	"github.com/microcosm-cc/bluemonday"
)

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
	BoardView    bool
	ErrMsg       string
}

func NewThreadsViewModel(id, postsPerPage int, bannerUrl, board_uri, topic string, op_post PostDto, replies []PostDto, boardview bool, pPol *bluemonday.Policy) (t ThreadViewModel) {
	t.Id = id
	t.BannerUrl = bannerUrl
	t.OPPost = op_post
	t.Replies = replies
	t.PostsPerPage = postsPerPage
	t.BoardUri = board_uri
	t.Topic = topic
	t.BoardView = boardview

	//Sanitize post content
	//TODO: Implement error handling here
	t.OPPost.Content, _ = RenderSafeMarkdown(t.OPPost.Content, pPol)

	for i := range t.Replies {
		t.Replies[i].Content, _ = RenderSafeMarkdown(t.Replies[i].Content, pPol)
	}

	return t
}
