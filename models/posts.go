package models

type PostDto struct {
	Id            int    `json:"id"`
	ThreadId      int    `json:"thread_id"`
	Identifier    string `json:"identifier"`
	Content       string `json:"content"`
	PostTimestamp string `json:"post_timestamp"`
	IsOP          bool   `json:"is_op"`
}
