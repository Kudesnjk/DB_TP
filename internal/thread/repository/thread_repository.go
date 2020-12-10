package repository

import (
	"database/sql"

	"github.com/Kudesnjk/DB_TP/internal/tools"

	"github.com/Kudesnjk/DB_TP/internal/models"

	"github.com/Kudesnjk/DB_TP/internal/thread"
)

type ThreadRepository struct {
	db *sql.DB
}

func NewThreadRepository(db *sql.DB) thread.ThreadRepository {
	return &ThreadRepository{
		db: db,
	}
}

func (tr *ThreadRepository) InsertThread(thread *models.Thread) error {
	err := tr.db.QueryRow("insert into threads(id, title, message, created, thread_nickname, forum_slug, votes) values (default, $1, $2, $3, $4, $5, default) returning id",
		thread.Title,
		thread.Message,
		thread.Created,
		thread.Author,
		thread.ForumSlug,
	).Scan(&thread.ID)
	return err
}

func (tr *ThreadRepository) SelectBySlugOrID(slugOrID string) (*models.Thread, error) {
	thread := &models.Thread{}

	err := tr.db.QueryRow("select id, slug, title, message, created, thread_nickname, forum_slug, votes, slug from threads where id = $1 or slug = $1",
		slugOrID).Scan(
		&thread.ID,
		&thread.Slug,
		&thread.Title,
		&thread.Message,
		&thread.Created,
		&thread.Author,
		&thread.ForumSlug,
		&thread.Votes,
		&thread.Slug)

	if err != nil {
		return nil, err
	}

	return thread, nil
}

func (tr *ThreadRepository) SelectThreadsByForumSlug(slug string, qpm *tools.QPM) ([]*models.Thread, error) {
	query := "select id, slug, title, message, created, thread_nickname, forum_slug, votes, slug from threads where forum_slug = $1"
	query = qpm.CreateQuery(query)

	rows, err := tr.db.Query(query,
		slug)

	if err != nil {
		return nil, err
	}

	threads := make([]*models.Thread, 0)

	for rows.Next() {
		thread := &models.Thread{}
		if err := rows.Scan(
			&thread.ID,
			&thread.Slug,
			&thread.Title,
			&thread.Message,
			&thread.Created,
			&thread.Author,
			&thread.ForumSlug,
			&thread.Votes,
			&thread.Slug,
		); err != nil {
			return nil, err
		}

		threads = append(threads, thread)
	}

	return threads, nil
}
