package repos

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/doug-martin/goqu/v9"
	"github.com/iraqnroll/gochan/models"
)

type PostgresBoardRepository struct {
	db *sql.DB
}

func NewPostgresBoardRepository(db *sql.DB) *PostgresBoardRepository {
	if db == nil {
		panic("Missing db")
	}

	return &PostgresBoardRepository{db: db}
}

func (r *PostgresBoardRepository) dbInstance() *goqu.Database {
	return pgDialect.DB(r.db)
}

func (r *PostgresBoardRepository) Delete(id int) error {
	_, err := r.dbInstance().Delete("boards").
		Where(goqu.C("id").Eq(id)).
		Executor().Exec()
	if err != nil {
		return fmt.Errorf("PostgresBoardRepository.Delete error: %w", err)
	}

	return nil
}

func (r *PostgresBoardRepository) CreateNew(uri, name, description string, owner_id int) (models.Board, error) {
	var result models.Board

	uri = strings.ToLower(uri)

	_, err := r.dbInstance().Insert("boards").
		Cols("uri", "name", "description", "ownerId").
		Vals([]interface{}{uri, name, description, owner_id}).
		Returning("id", goqu.L("to_char(date_created, 'YYYY-MM-DD HH24:MI:SS')").As("date_created")).
		Executor().
		ScanStruct(&result)
	if err != nil {
		return result, fmt.Errorf("PostgresBoardRepository.CreateNew error: %w", err)
	}

	result.Uri = uri
	result.Name = name
	result.Description = description
	result.OwnerId = owner_id
	return result, nil
}

func (r *PostgresBoardRepository) GetAllForAdmin() ([]models.Board, error) {
	var result []models.Board

	err := r.dbInstance().From("boards").
		Select(
			"boards.id",
			"boards.uri",
			"boards.name",
			"boards.description",
			goqu.L("COALESCE(to_char(boards.date_created, 'YYYY-MM-DD HH24:MI:SS'), 'Never')").As("date_created"),
			goqu.L("COALESCE(to_char(boards.date_updated, 'YYYY-MM-DD HH24:MI:SS'), 'Never')").As("date_updated"),
			goqu.I("users.username").As("owner_username"),
		).
		Join(goqu.T("users"), goqu.On(goqu.I("users.id").Eq(goqu.I("boards.ownerid")))).
		ScanStructs(&result)
	if err != nil {
		return nil, fmt.Errorf("PostgresBoardRepository.GetAllForAdmin error: %w", err)
	}

	return result, nil
}

func (r *PostgresBoardRepository) GetByUri(uri string) (models.BoardDto, error) {
	var result models.BoardDto

	found, err := r.dbInstance().From("boards").
		Select("id", "name", "description").
		Where(goqu.C("uri").Eq(uri)).
		ScanStruct(&result)
	if err != nil {
		return result, fmt.Errorf("PostgresBoardRepository.GetByUri error: %w", err)
	}
	if !found {
		return result, fmt.Errorf("PostgresBoardRepository.GetByUri: board not found")
	}

	result.Uri = uri
	return result, nil
}

func (r *PostgresBoardRepository) GetById(id int) (models.BoardDto, error) {
	var result models.BoardDto

	found, err := r.dbInstance().From("boards").
		Select("uri", "name", "description").
		Where(goqu.C("id").Eq(id)).
		ScanStruct(&result)
	if err != nil {
		return result, fmt.Errorf("PostgresBoardRepository.GetById error: %w", err)
	}
	if !found {
		return result, fmt.Errorf("PostgresBoardRepository.GetById: board not found")
	}

	result.Id = id
	return result, nil
}

func (r *PostgresBoardRepository) GetAll() ([]models.BoardDto, error) {
	var result []models.BoardDto

	err := r.dbInstance().From("boards").
		Select("id", "uri", "name", "description").
		ScanStructs(&result)
	if err != nil {
		return nil, fmt.Errorf("PostgresBoardRepository.GetAll error: %w", err)
	}

	return result, nil
}

func (r *PostgresBoardRepository) Update(board_id int, uri, name, description string) (*models.BoardDto, error) {
	_, err := r.dbInstance().Update("boards").
		Set(goqu.Record{
			"uri":          uri,
			"name":         name,
			"description":  description,
			"date_updated": goqu.L("NOW()"),
		}).
		Where(goqu.C("id").Eq(board_id)).
		Executor().Exec()
	if err != nil {
		return nil, fmt.Errorf("PostgresBoardRepository.Update error: %w", err)
	}

	result := &models.BoardDto{
		Id:          board_id,
		Uri:         uri,
		Name:        name,
		Description: description,
	}
	return result, nil
}
