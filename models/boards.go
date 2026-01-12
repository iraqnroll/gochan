package models

import "github.com/microcosm-cc/bluemonday"

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

type BoardViewModel struct {
	Id             int
	Uri            string
	Name           string
	Description    string
	BannerUrl      string
	Threads        []ThreadViewModel
	ThreadsPerPage int
	ErrMsg         string
}

func NewBoardViewModel(id, threads_per_page int, uri, name, desc, banner_url string, threads []ThreadDto, pPol *bluemonday.Policy) (m BoardViewModel) {
	m.Id = id
	m.ThreadsPerPage = threads_per_page
	m.Uri = uri
	m.Name = name
	m.Description = desc
	m.BannerUrl = banner_url

	for _, thread := range threads {
		result := NewThreadsViewModel(thread.Id, 10, banner_url, uri, thread.Topic, thread.Posts[0], thread.Posts[1:], true, pPol)
		m.Threads = append(m.Threads, result)
	}

	return m
}
