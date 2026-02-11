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
	Update(board_id int, uri, name, description string) (*models.BoardDto, error)
}

type IThreadService interface {
	CreateThread(board_id int, topic, identifier, content, fingerprint string) (*models.ThreadDto, error)
	GetThread(thread_id int) (*models.ThreadDto, error)
	GetBoardThreads(board_id int) ([]models.ThreadDto, error)
}

type IFileService interface {
	CreateBoardStatic(board_uri string) error
	RemoveBoardStatic(board_uri string) error
}

type BoardService struct {
	boardRepo BoardRepository
	thService IThreadService
	fService  IFileService
}

func NewBoardService(repo BoardRepository, thService IThreadService, fService IFileService) *BoardService {
	return &BoardService{boardRepo: repo, thService: thService, fService: fService}
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

// Fetches board by id (only board metadata, no threads/posts are attached.)
func (bs *BoardService) GetBoardById(id int) (*models.BoardDto, error) {
	result, err := bs.boardRepo.GetById(id)
	if err != nil {
		return nil, fmt.Errorf("BoardService.GetBoardById failed : %w", err)
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

	err = bs.fService.CreateBoardStatic(board.Uri)
	if err != nil {
		return nil, fmt.Errorf("BoardService.Create failed : %w", err)
	}

	return &board, nil
}

func (bs *BoardService) Delete(boardId int, boardUri string) error {
	err := bs.boardRepo.Delete(boardId)
	if err != nil {
		return fmt.Errorf("BoardService.Delete failed : %w", err)
	}

	err = bs.fService.RemoveBoardStatic(boardUri)
	if err != nil {
		return fmt.Errorf("BoardService.Delete failed : %w", err)
	}

	return nil
}

func (bs *BoardService) UpdateBoard(board_id int, uri, name, description string) (*models.BoardDto, error) {
	result, err := bs.boardRepo.Update(board_id, uri, name, description)
	if err != nil {
		return nil, fmt.Errorf("BoardService.UpdateBoard failed: %w", err)
	}

	return result, nil
}
