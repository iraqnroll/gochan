package repos

import (
	"database/sql"
	"fmt"

	"github.com/iraqnroll/gochan/models"
)

const (
	t_get_by_id_query    = `SELECT locked, board_id, topic FROM threads WHERE id = $1`
	t_get_by_board_query = `SELECT id, locked, topic FROM threads WHERE board_id = $1`
	t_create_new_query   = `INSERT INTO threads(board_id, topic) VALUES ($1, $2) RETURNING id`
)

type PostgresThreadRepository struct {
	db *sql.DB
}

func NewPostgresThreadRepository(db *sql.DB) *PostgresThreadRepository {
	if db == nil {
		panic("Missing db")
	}

	return &PostgresThreadRepository{db: db}
}

func (r *PostgresThreadRepository) GetById(thread_id int) (models.ThreadDto, error) {
	result := models.ThreadDto{Id: thread_id}
	row := r.db.QueryRow(t_get_by_id_query, result.Id)
	err := row.Scan(&result.Locked, &result.BoardId, &result.Topic)
	if err != nil {
		return result, fmt.Errorf("PostgresThreadRepository.GetById error : %w", err)
	}

	return result, nil
}

func (r *PostgresThreadRepository) CreateNew(board_id int, topic string) (models.ThreadDto, error) {
	result := models.ThreadDto{BoardId: board_id, Topic: topic}
	row := r.db.QueryRow(t_create_new_query, result.BoardId, result.Topic)
	err := row.Scan(&result.Id)
	if err != nil {
		return result, fmt.Errorf("PostgresThreadRepository.CreateNew error : %w", err)
	}

	return result, nil
}

func (r *PostgresThreadRepository) GetAllByBoard(board_id int) ([]models.ThreadDto, error) {
	var result []models.ThreadDto
	rows, err := r.db.Query(t_get_by_board_query, board_id)
	if err != nil {
		return nil, fmt.Errorf("PostgresThreadRepository.GetAllByBoard error : %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		thread := models.ThreadDto{BoardId: board_id}
		err := rows.Scan(&thread.Id, &thread.Locked, &thread.Topic)
		if err != nil {
			return result, fmt.Errorf("PostgresThreadRepository.GetAllByBoard error : %w", err)
		}
		result = append(result, thread)
	}

	return result, nil
}
