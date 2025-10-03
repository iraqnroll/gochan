package models

import (
	"database/sql"
	"fmt"
	"strings"
)

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
	Id          int
	Uri         string
	Name        string
	Description string
}

type BoardService struct {
	DB *sql.DB
}

func (bs *BoardService) Create(uri, name, description string, ownerId int) (*Board, error) {
	uri = strings.ToLower(uri)

	board := Board{
		Uri:         uri,
		Name:        name,
		Description: description,
		OwnerId:     ownerId,
	}

	row := bs.DB.QueryRow(`
		INSERT INTO boards (uri, name, description, ownerId)
		VALUES ($1, $2, $3, $4) RETURNING id, to_char(date_created, 'YYYY-MM-DD HH24:MI:SS')`, uri, name, description, ownerId)

	err := row.Scan(&board.Id, &board.Date_created)

	if err != nil {
		return nil, fmt.Errorf("BoardService.Create failed : %w", err)
	}

	return &board, nil
}

func (bs *BoardService) Delete(boardId int) error {
	_, err := bs.DB.Exec(`DELETE FROM boards WHERE id = $1`, boardId)
	if err != nil {
		fmt.Println("BoardService.Delete failed : %w", err)
		return err
	}

	return nil
}

func (bs *BoardService) GetAdminBoardList() ([]Board, error) {
	var result []Board

	rows, err := bs.DB.Query(`SELECT
		b.id,
		b.uri,
		b.name,
		b.description,
		COALESCE(to_char(b.date_created, 'YYYY-MM-DD HH24:MI:SS'), 'Never') AS date_created,
		COALESCE(to_char(b.date_updated, 'YYYY-MM-DD HH24:MI:SS'), 'Never') AS date_updated,
		usr.username AS ownerUsername
		FROM boards AS b
		INNER JOIN users AS usr ON usr.id = b.ownerId`)

	if err != nil {
		return nil, fmt.Errorf("BoardService.GetAdminBoardList failed : %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var board Board
		err := rows.Scan(&board.Id, &board.Uri, &board.Name, &board.Description, &board.Date_created, &board.Date_updated, &board.OwnerUsername)
		if err != nil {
			fmt.Println("BoardService.GetAdminBoardList loop failed : %w", err)
		}
		result = append(result, board)
	}

	return result, nil
}

func (bs *BoardService) GetBoardList() ([]BoardDto, error) {
	var result []BoardDto

	rows, err := bs.DB.Query(`SELECT id, uri, name, description FROM boards`)

	if err != nil {
		return nil, fmt.Errorf("BoardService.GetBoardList failed : %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var board BoardDto
		err := rows.Scan(&board.Id, &board.Uri, &board.Name, &board.Description)
		if err != nil {
			fmt.Println("BoardService.GetBoardList loop failed : %w", err)
		}
		result = append(result, board)
	}

	return result, nil
}

func (bs *BoardService) GetBoard(uri string) (*BoardDto, error) {
	var result BoardDto
	rows := bs.DB.QueryRow(`SELECT id, uri, name, description FROM boards WHERE uri = $1`, uri)
	err := rows.Scan(&result.Id, &result.Uri, &result.Name, &result.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}

		return nil, fmt.Errorf("BoardService.GetBoard failed : %w", err)
	}

	return &result, nil
}
