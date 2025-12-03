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

type BoardViewModel struct {
	Id             int
	Uri            string
	Name           string
	Description    string
	BannerUrl      string
	Threads        []ThreadDto
	ThreadsPerPage int
	ErrMsg         string
}

func NewBoardViewModel(id, threads_per_page int, uri, name, desc, banner_url string, threads []ThreadDto) (m BoardViewModel) {
	m.Id = id
	m.ThreadsPerPage = threads_per_page
	m.Uri = uri
	m.Name = name
	m.Description = desc
	m.BannerUrl = banner_url
	m.Threads = threads

	return m
}
