package models

type PostDto struct {
	Id            int    `json:"id"`
	ThreadId      int    `json:"thread_id"`
	Identifier    string `json:"identifier"`
	Content       string `json:"content"`
	PostTimestamp string `json:"post_timestamp"`
	IsOP          bool   `json:"is_op"`
	HasMedia      string
}

type RecentPostsDto struct {
	Board_uri      string `json:"board_uri"`
	Board_name     string `json:"board_name"`
	Thread_id      int    `json:"thread_id"`
	Thread_topic   string `json:"thread_topic"`
	Post_id        int    `json:"post_id"`
	Post_ident     string `json:"post_ident"`
	Post_content   string `json:"post_content"`
	Post_timestamp string `json:"post_timestamp"`
	HasMedia       string
}
