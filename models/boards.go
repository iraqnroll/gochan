package models

type Board struct {
	Id            int
	Uri           string
	Name          string
	Description   string
	Date_created  string
	Date_updated  string
	OwnerId       int
	OwnerUsername string
}

type BoardDto struct {
	Id          int         `json:"id"`
	Uri         string      `json:"uri"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Threads     []ThreadDto `json:"threads"`
}
