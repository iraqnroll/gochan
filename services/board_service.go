package services

import (
	"fmt"
	"strings"

	"github.com/iraqnroll/gochan/models"
)

type BoardRepository interface {
	GetAll() ([]models.BoardDto, error)
	GetById(id int) (models.BoardDto, error)
	GetByUri(uri string) (models.BoardDto, error)
	GetAllForAdmin() ([]models.Board, error)

	CreateNew(uri, name, description string, owner_id int) (models.Board, error)
	Delete(id int) error
}

type BoardService struct {
	boardRepo BoardRepository
	thService *ThreadService
}

// Fetches board and it's content (threads/posts) for a specified board uri
func (bs *BoardService) GetBoard(uri string) (*models.BoardDto, error) {
	result, err := bs.boardRepo.GetByUri(uri)
	if err != nil {
		return nil, fmt.Errorf("BoardService.GetBoard failed : %w", err)
	}

	result.Threads, err = bs.thService.GetBoardThreads(result.Id)
	if err != nil {
		return nil, fmt.Errorf("BoardService.GetBoard failed : %w", err)
	}

	return &result, nil
}

// Returns a list of all registered boards
func (bs *BoardService) GetBoardList() ([]models.BoardDto, error) {
	boards, err := bs.boardRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("BoardService.GetBoardList failed : %w", err)
	}

	return boards, nil
}

// Returns a list of all registered boards (with extra data for admins)
func (bs *BoardService) GetAdminBoardList() ([]models.Board, error) {
	boards, err := bs.boardRepo.GetAllForAdmin()
	if err != nil {
		return nil, fmt.Errorf("BoardService.GetAdminBoardList failed : %w", err)
	}

	return boards, nil
}

func (bs *BoardService) Create(uri, name, description string, ownerId int) (*models.Board, error) {
	uri = strings.ToLower(uri)
	board, err := bs.boardRepo.CreateNew(uri, name, description, ownerId)
	if err != nil {
		return nil, fmt.Errorf("BoardService.Create failed : %w", err)
	}

	return &board, nil
}

func (bs *BoardService) Delete(boardId int, boardUri string) error {
	err := bs.boardRepo.Delete(boardId)
	if err != nil {
		fmt.Println("BoardService.Delete failed : %w", err)
		return err
	}

	return nil
}
