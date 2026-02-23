package models

type Board struct {
	Id            int    `db:"id" schema:"id"`
	Uri           string `db:"uri" schema:"uri"`
	Name          string `db:"name" schema:"name"`
	Description   string `db:"description" schema:"description"`
	Date_created  string `db:"date_created" schema:"date_created"`
	Date_updated  string `db:"date_updated" schema:"date_updated"`
	OwnerId       int    `db:"ownerId" schema:"ownerId"`
	OwnerUsername string `db:"owner_username" schema:"owner_username"`
}

type BoardDto struct {
	Id          int         `db:"id" json:"id" schema:"id"`
	Uri         string      `db:"uri" json:"uri" schema:"uri"`
	Name        string      `db:"name" json:"name" schema:"name"`
	Description string      `db:"description" json:"description" schema:"description"`
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

func NewBoardViewModel(id, threads_per_page int, uri, name, desc, banner_url string, threads []ThreadDto) (m BoardViewModel) {
	m.Id = id
	m.ThreadsPerPage = threads_per_page
	m.Uri = uri
	m.Name = name
	m.Description = desc
	m.BannerUrl = banner_url

	for _, thread := range threads {
		result := NewThreadsViewModel(thread.Id, 10, banner_url, uri, thread.Topic, thread.Posts[0], thread.Posts[1:], true)
		m.Threads = append(m.Threads, result)
	}

	return m
}
