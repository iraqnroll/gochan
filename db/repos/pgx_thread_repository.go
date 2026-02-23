package repos

import (
	"database/sql"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/iraqnroll/gochan/db/models"
)

var pgDialect = goqu.Dialect("postgres")

type PostgresThreadRepository struct {
	db *sql.DB
}

func NewPostgresThreadRepository(db *sql.DB) *PostgresThreadRepository {
	if db == nil {
		panic("Missing db")
	}

	return &PostgresThreadRepository{db: db}
}

func (r *PostgresThreadRepository) dbInstance() *goqu.Database {
	return pgDialect.DB(r.db)
}

func (r *PostgresThreadRepository) GetById(thread_id int) (models.ThreadDto, error) {
	var result models.ThreadDto

	found, err := r.dbInstance().From("threads").
		Select("locked", "board_id", "topic").
		Where(goqu.C("id").Eq(thread_id)).
		ScanStruct(&result)
	if err != nil {
		return result, fmt.Errorf("PostgresThreadRepository.GetById error: %w", err)
	}
	if !found {
		return result, fmt.Errorf("PostgresThreadRepository.GetById: thread not found")
	}

	result.Id = thread_id
	return result, nil
}

func (r *PostgresThreadRepository) CreateNew(board_id int, topic string) (models.ThreadDto, error) {
	var result models.ThreadDto

	_, err := r.dbInstance().Insert("threads").
		Cols("board_id", "topic").
		Vals([]interface{}{board_id, topic}).
		Returning("id").
		Executor().
		ScanStruct(&result)
	if err != nil {
		return result, fmt.Errorf("PostgresThreadRepository.CreateNew error: %w", err)
	}

	result.BoardId = board_id
	result.Topic = topic
	return result, nil
}

func (r *PostgresThreadRepository) GetAllByBoard(board_id int) ([]models.ThreadDto, error) {
	var result []models.ThreadDto

	err := r.dbInstance().From("threads").
		Select("id", "locked", "topic").
		Where(goqu.C("board_id").Eq(board_id)).
		Order(goqu.C("id").Desc()).
		ScanStructs(&result)
	if err != nil {
		return nil, fmt.Errorf("PostgresThreadRepository.GetAllByBoard error: %w", err)
	}

	for i := range result {
		result[i].BoardId = board_id
	}

	return result, nil
}
