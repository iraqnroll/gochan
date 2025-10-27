package repos

import (
	"database/sql"
	"fmt"

	"github.com/iraqnroll/gochan/models"
)

const (
	p_get_all_by_thread_query = `SELECT
		p.id,
		p.identifier,
		p.content,
		COALESCE(to_char(p.post_timestamp, 'YYYY-MM-DD HH24:MI:SS'), 'Never') AS post_timestamp,
		p.is_op
		FROM posts AS p
		INNER JOIN threads AS th ON th.id = p.thread_id
		WHERE th.id = $1`
	p_create_new_query = `INSERT INTO posts(thread_id, identifier, content, is_op) VALUES ($1, $2, $3, $4) RETURNING id`
)

type PostgresPostRepository struct {
	db *sql.DB
}

func NewPostgresPostRepository(db *sql.DB) *PostgresPostRepository {
	if db == nil {
		panic("Missing db")
	}

	return &PostgresPostRepository{db: db}
}
func (r *PostgresPostRepository) CreateNew(thread_id int, identifier, content string, is_op bool) (models.PostDto, error) {
	result := models.PostDto{ThreadId: thread_id, Identifier: identifier, Content: content, IsOP: is_op}
	row := r.db.QueryRow(p_create_new_query, result.ThreadId, result.Identifier, result.Content, result.IsOP)
	err := row.Scan(&result.Id)
	if err != nil {
		return result, fmt.Errorf("PostgresPostRepository.CreateNew error : %w", err)
	}

	return result, nil
}

func (r *PostgresPostRepository) GetAllByThread(thread_id int) ([]models.PostDto, error) {
	var result []models.PostDto
	rows, err := r.db.Query(p_get_all_by_thread_query, thread_id)
	if err != nil {
		return nil, fmt.Errorf("PostgresPostRepository.GetAllByThread error : %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		post := models.PostDto{ThreadId: thread_id}
		err := rows.Scan(&post.Id, &post.Identifier, &post.Content, &post.PostTimestamp, &post.IsOP)
		if err != nil {
			return nil, fmt.Errorf("PostgresPostRepository.GetAllByThread error : %w", err)
		}

		result = append(result, post)
	}

	return result, nil
}
