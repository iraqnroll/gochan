package repos

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/iraqnroll/gochan/models"
)

const (
	b_get_all_query           = `SELECT id, uri, name, description FROM boards`
	b_get_by_id_query         = `SELECT uri, name, description FROM boards where id = $1`
	b_get_by_uri_query        = `SELECT id, name, description FROM boards where uri = $1`
	b_get_all_for_admin_query = `SELECT 
		b.id,
		b.uri,
		b.name,
		b.description,
		COALESCE(to_char(b.date_created, 'YYYY-MM-DD HH24:MI:SS'), 'Never') AS date_created,
		COALESCE(to_char(b.date_updated, 'YYYY-MM-DD HH24:MI:SS'), 'Never') AS date_updated,
		usr.username AS ownerUsername
		FROM boards AS b
		INNER JOIN users AS usr ON usr.id = b.ownerId`
	b_create_new_query = `INSERT INTO boards (uri, name, description, ownerId)
		VALUES ($1,$2,$3,$4) RETURNING id, to_char(date_created, 'YYYY-MM-DD HH24:MI:SS')`
	b_delete_query = `DELETE FROM boards WHERE id = $1`
)

// Implementation of the BoardRepository for pgx/Postgres
type PostgresBoardRepository struct {
	db *sql.DB
}

func NewPostgresBoardRepository(db *sql.DB) *PostgresBoardRepository {
	if db == nil {
		panic("Missing db")
	}

	return &PostgresBoardRepository{db: db}
}

func (r *PostgresBoardRepository) Delete(id int) error {
	_, err := r.db.Exec(b_delete_query, id)
	if err != nil {
		return fmt.Errorf("PostgresBoardRepository.Delete error : %w", err)
	}

	return nil
}

func (r *PostgresBoardRepository) CreateNew(uri, name, description string, owner_id int) (models.Board, error) {
	result := models.Board{
		Uri:         strings.ToLower(uri),
		Name:        name,
		Description: description,
		OwnerId:     owner_id,
	}

	row := r.db.QueryRow(b_create_new_query, result.Uri, result.Name, result.Description, result.OwnerId)
	err := row.Scan(&result.Id, &result.Date_created)
	if err != nil {
		return result, fmt.Errorf("PostgresBoardRepository.CreateNew error : %w", err)
	}

	return result, nil
}

func (r *PostgresBoardRepository) GetAllForAdmin() ([]models.Board, error) {
	var result []models.Board
	rows, err := r.db.Query(b_get_all_for_admin_query)
	if err != nil {
		return nil, fmt.Errorf("PostgresBoardRepository.GetAllForAdmin error : %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var board models.Board
		err := rows.Scan(
			&board.Id,
			&board.Uri,
			&board.Name,
			&board.Description,
			&board.Date_created,
			&board.Date_updated,
			&board.OwnerUsername)

		if err != nil {
			return result, fmt.Errorf("PostgresBoardRepository.GetAllForAdmin error : %w", err)
		}
		result = append(result, board)
	}

	return result, nil
}

func (r *PostgresBoardRepository) GetByUri(uri string) (models.BoardDto, error) {
	result := models.BoardDto{Uri: uri}

	rows := r.db.QueryRow(b_get_by_uri_query, result.Uri)
	err := rows.Scan(&result.Id, &result.Name, &result.Description)

	if err != nil {
		return result, fmt.Errorf("PostgresBoardRepository.GetByUri error : %w", err)
	}

	return result, nil
}

func (r *PostgresBoardRepository) GetById(id int) (models.BoardDto, error) {
	result := models.BoardDto{Id: id}

	rows := r.db.QueryRow(b_get_by_id_query, result.Id)
	err := rows.Scan(&result.Uri, &result.Name, &result.Description)

	if err != nil {
		return result, fmt.Errorf("PostgresBoardRepository.GetById error : %w", err)
	}

	return result, nil
}

func (r *PostgresBoardRepository) GetAll() ([]models.BoardDto, error) {
	var result []models.BoardDto

	rows, err := r.db.Query(b_get_all_query)
	if err != nil {
		return nil, fmt.Errorf("PostgresBoardRepository.GetAll error : %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var board models.BoardDto
		err := rows.Scan(&board.Id, &board.Uri, &board.Name, &board.Description)
		if err != nil {
			return nil, fmt.Errorf("PostgresBoardRepository.GetAll error : %w", err)
		}
		result = append(result, board)
	}

	return result, nil
}
