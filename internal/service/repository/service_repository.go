package repository

import (
	"database/sql"

	"github.com/Kudesnjk/DB_TP/internal/models"

	"github.com/Kudesnjk/DB_TP/internal/service"
)

type ServiceRepository struct {
	db *sql.DB
}

func NewServiceRepository(db *sql.DB) service.ServiceRepository {
	return &ServiceRepository{
		db: db,
	}
}

func (sr *ServiceRepository) SelectStatus() (*models.Service, error) {
	status := &models.Service{}

	err := sr.db.QueryRow(`
	select
	COALESCE(sum(posts_num), 0),
	COALESCE(sum(threads_num), 0),
	count(slug),
	count(distinct u.nickname) from forums
	full outer join users u on forums.user_nickname = u.nickname`).
		Scan(&status.PostsNum,
			&status.ThreadsNum,
			&status.ForumsNum,
			&status.UsersNum)

	if err != nil {
		return nil, err
	}

	return status, nil
}

func (sr *ServiceRepository) TruncateTables() error {
	tx, err := sr.db.Begin()

	if err != nil {
		return err
	}

	_, err = tx.Exec("truncate table users cascade")

	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("truncate table forums cascade")

	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("truncate table threads cascade")

	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("truncate table posts cascade")

	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("truncate table votes cascade")

	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
