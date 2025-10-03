package controllers

import (
	"net/http"

	"github.com/iraqnroll/gochan/models"
)

type Home struct {
	Home         Template
	BoardService *models.BoardService
}

func (h Home) HomePage(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Boards []models.BoardDto
		Test   string
	}

	boards, err := h.BoardService.GetBoardList()
	if err != nil {
		http.Error(w, "Unable to fetch board list"+err.Error(), http.StatusInternalServerError)
	}

	data.Boards = boards
	data.Test = "String"

	h.Home.Execute(w, r, data)
}
