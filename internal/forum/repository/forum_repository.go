package repository

import (
	"database/sql"

	"github.com/Kudesnjk/DB_TP/internal/models"

	"github.com/Kudesnjk/DB_TP/internal/forum"
)

type ForumRepository struct {
	db *sql.DB
}

func NewForumRepository(db *sql.DB) forum.ForumRepository {
	return &ForumRepository{
		db: db,
	}
}

func (fr *ForumRepository) InsertForum(forum *models.Forum) error {
	_, err := fr.db.Exec("insert into forums (slug, title, user_nickname) values($1, $2, $3)",
		forum.Slug,
		forum.Title,
		forum.User)
	return err
}

func (fr *ForumRepository) SelectByForumSlug(slug string) (*models.Forum, error) {
	forum := &models.Forum{}

	err := fr.db.QueryRow("select slug, title, user_nickname, threads_num, posts_num from forums where lower(slug) = lower($1)", slug).Scan(
		&forum.Slug,
		&forum.Title,
		&forum.User,
		&forum.Threads,
		&forum.Posts,
	)

	if err != nil {
		return nil, err
	}

	return forum, nil
}

func (fr *ForumRepository) SelectForumUsers(slug string) ([]*models.User, error) {
	rows, err := fr.db.Query(`select u.nickname from users u 
	join threads t on u.nickname = t.user_nickname where forum_slug = $1
	union all select u.nickname from users u 
	join posts p on u.nickname = p.user_nickname 
	join threads t on p.thread_id = t.id where forum_slug = $1`,
		slug)

	if err != nil {
		return nil, err
	}

	users := make([]*models.User, 0)

	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(
			&user.Fullname,
			&user.Nickname,
			&user.Email,
			&user.About,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
