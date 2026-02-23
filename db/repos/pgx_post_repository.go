package repos

import (
	"database/sql"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/iraqnroll/gochan/db/models"
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

func (r *PostgresPostRepository) dbInstance() *goqu.Database {
	return pgDialect.DB(r.db)
}

func (r *PostgresPostRepository) CreateNew(thread_id int, identifier, content, fingerprint string, is_op bool) (models.PostDto, error) {
	var result models.PostDto

	_, err := r.dbInstance().Insert("posts").
		Cols("thread_id", "identifier", "content", "is_op", "fingerprint").
		Vals([]interface{}{thread_id, identifier, content, is_op, fingerprint}).
		Returning("id").
		Executor().
		ScanStruct(&result)
	if err != nil {
		return result, fmt.Errorf("PostgresPostRepository.CreateNew error: %w", err)
	}

	result.ThreadId = thread_id
	result.Identifier = identifier
	result.Content = content
	result.IsOP = is_op
	result.Post_fprint = fingerprint
	return result, nil
}

func (r *PostgresPostRepository) GetAllByThread(thread_id int, for_mod bool) ([]models.PostDto, error) {
	var result []models.PostDto
	var whereExpr goqu.Ex

	cols := []interface{}{
		"id",
		"identifier",
		"content",
		goqu.L("COALESCE(to_char(post_timestamp, 'YYYY-MM-DD HH24:MI:SS'), 'Never')").As("post_timestamp"),
		"is_op",
		goqu.L("COALESCE(has_media, '')").As("has_media"),
	}

	if for_mod {
		cols = append(cols, goqu.L("COALESCE(fingerprint, '')").As("fingerprint"))
		cols = append(cols, "deleted")
		whereExpr = goqu.Ex{
			"thread_id": thread_id,
		}
	} else {
		whereExpr = goqu.Ex{
			"thread_id": thread_id,
			"deleted":   false,
		}
	}

	err := r.dbInstance().From("posts").
		Select(cols...).
		Where(whereExpr).
		Order(goqu.C("id").Asc()).
		ScanStructs(&result)
	if err != nil {
		return nil, fmt.Errorf("PostgresPostRepository.GetAllByThread error: %w", err)
	}

	for i := range result {
		result[i].ThreadId = thread_id
	}

	return result, nil
}

func (r *PostgresPostRepository) GetMostRecent(num_of_posts int) ([]models.RecentPostsDto, error) {
	var result []models.RecentPostsDto

	err := r.dbInstance().From("recent_posts").
		Select(
			"board_uri",
			"board_name",
			"thread_id",
			"thread_topic",
			"post_id",
			"post_ident",
			goqu.L("SUBSTRING(post_content for 100)").As("post_content"),
			goqu.L("COALESCE(to_char(post_timestamp, 'YYYY-MM-DD HH24:MI:SS'), 'Never')").As("post_timestamp"),
		).
		Order(goqu.C("post_timestamp").Desc()).
		Limit(uint(num_of_posts)).
		ScanStructs(&result)
	if err != nil {
		return nil, fmt.Errorf("PostgresPostRepository.GetMostRecent error: %w", err)
	}

	return result, nil
}

func (r *PostgresPostRepository) UpdateAttachedMedia(post_id int, attached_media, original_media string) error {
	_, err := r.dbInstance().Update("posts").
		Set(goqu.Record{
			"has_media": attached_media,
			"og_media":  original_media,
		}).
		Where(goqu.C("id").Eq(post_id)).
		Executor().Exec()
	if err != nil {
		return fmt.Errorf("PostgresPostRepository.UpdateAttachedMedia error: %w", err)
	}
	return nil
}

func (r *PostgresPostRepository) SoftDeletePost(post_id int) error {
	_, err := r.dbInstance().Update("posts").
		Set(goqu.Record{
			"deleted": true,
		}).
		Where(goqu.C("id").Eq(post_id)).
		Executor().Exec()

	if err != nil {
		return fmt.Errorf("PostgresPostRepository.SoftDeletePost error: %w", err)
	}
	return nil
}

func (r *PostgresPostRepository) RemoveSoftDeleteFromPost(post_id int) error {
	_, err := r.dbInstance().Update("posts").
		Set(goqu.Record{
			"deleted": false,
		}).
		Where(goqu.C("id").Eq(post_id)).
		Executor().Exec()

	if err != nil {
		return fmt.Errorf("PostgresPostRepository.RemoveSoftDeleteFromPost error: %w", err)
	}
	return nil
}
