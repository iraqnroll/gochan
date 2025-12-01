package models

type HomeViewModel struct {
	Boards []BoardDto
	ErrMsg string
}

func NewHomeViewModel(boards []BoardDto, err error) (m HomeViewModel) {
	if len(boards) <= 0 {
		m.ErrMsg = "Failed to retrieve any boards..."
	} else {
		m.Boards = boards
	}

	return m
}
