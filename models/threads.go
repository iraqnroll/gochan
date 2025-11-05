package models

type ThreadDto struct {
	Id       int       `json:"id"`
	Posts    []PostDto `json:"posts"`
	Topic    string    `json:"topic"`
	Locked   bool      `json:"locked"`
	BoardId  int       `json:"board_id"`
	BoardUri string    `json:"board_uri"`
}
