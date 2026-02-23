package models

type PostDto struct {
	Id            int    `db:"id" json:"id" schema:"id"`
	ThreadId      int    `db:"thread_id" json:"thread_id" schema:"threadId"`
	Identifier    string `db:"identifier" json:"identifier" schema:"identifier"`
	Content       string `db:"content" json:"content" schema:"content"`
	PostTimestamp string `db:"post_timestamp" json:"post_timestamp" schema:"post_timestamp"`
	IsOP          bool   `db:"is_op" json:"is_op" schema:"is_op"`
	Post_fprint   string `db:"fingerprint" json:"post_fprint" schema:"post_fprint"`
	Deleted       bool   `db:"deleted" json:"deleted" schema:"deleted"`
	OgMedia       string `db:"og_media" schema:"og_media"`
	HasMedia      string `db:"has_media" schema:"has_media"`
	Tripcode      string `db:"tripcode"`
}

type RecentPostsDto struct {
	Board_uri      string `db:"board_uri" json:"board_uri" schema:"board_uri"`
	Board_name     string `db:"board_name" json:"board_name" schema:"board_name"`
	Thread_id      int    `db:"thread_id" json:"thread_id" schema:"thread_id"`
	Thread_topic   string `db:"thread_topic" json:"thread_topic" schema:"thread_topic"`
	Post_id        int    `db:"post_id" json:"post_id" schema:"post_id"`
	Post_ident     string `db:"post_ident" json:"post_ident" schema:"post_ident"`
	Post_content   string `db:"post_content" json:"post_content" schema:"post_content"`
	Post_timestamp string `db:"post_timestamp" json:"post_timestamp" schema:"post_timestamp"`
	HasMedia       string `db:"has_media" schema:"has_media"`
}
