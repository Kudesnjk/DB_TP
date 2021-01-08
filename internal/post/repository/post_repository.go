package repository

import (
	"database/sql"

	"github.com/lib/pq"

	"github.com/Kudesnjk/DB_TP/internal/tools"

	"github.com/Kudesnjk/DB_TP/internal/models"

	"github.com/Kudesnjk/DB_TP/internal/post"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) post.PostRepository {
	return &PostRepository{
		db: db,
	}
}

func (pr *PostRepository) InsertPost(post *models.Post) error {
	if post.Parent != 0 {
		err := pr.db.QueryRow(`with parent_path as (select path from posts where posts.id = $3 and thread_id = $5)
			insert into posts(id, message, is_edited, created, parent_id, user_nickname, thread_id, forum_slug, path)
			values(default, $1, default, $2, $3, $4, $5, $6, (select path from parent_path)) returning id`,
			post.Message,
			post.Created,
			post.Parent,
			post.Author,
			post.ThreadID,
			post.ForumSlug).
			Scan(
				&post.ID,
			)

		return err
	}

	err := pr.db.QueryRow(`insert into posts 
	(id, message, is_edited, created, parent_id, user_nickname, thread_id, forum_slug, path) 
	values(default, $1, default, $2, $3, $4, $5, $6, default) returning id`,
		post.Message,
		post.Created,
		post.Parent,
		post.Author,
		post.ThreadID,
		post.ForumSlug).
		Scan(
			&post.ID,
		)
	return err
}

func (pr *PostRepository) SelectPosts(threadID uint64, qpm *tools.QPM) ([]*models.Post, error) {
	query := `select id, message, is_edited, created, parent_id, user_nickname, forum_slug, thread_id, path from posts where thread_id = $1`
	query = qpm.UpdatePostQuery(query)

	rows, err := pr.db.Query(query, threadID)

	if err != nil {
		return nil, err
	}

	posts := make([]*models.Post, 0)

	for rows.Next() {
		post := &models.Post{}
		if err := rows.Scan(
			&post.ID,
			&post.Message,
			&post.IsEdited,
			&post.Created,
			&post.Parent,
			&post.Author,
			&post.ForumSlug,
			&post.ThreadID,
			pq.Array(&post.Path),
		); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (pr *PostRepository) UpdatePost(post *models.Post) error {
	query := `update posts set message = $1, is_edited = true where id = $2 returning id, message, is_edited, created, parent_id, user_nickname, forum_slug, thread_id, path`

	err := pr.db.QueryRow(query, post.Message, post.ID).
		Scan(&post.ID,
			&post.Message,
			&post.IsEdited,
			&post.Created,
			&post.Parent,
			&post.Author,
			&post.ForumSlug,
			&post.ThreadID,
			pq.Array(&post.Path))

	if err != nil {
		return err
	}
	return nil
}

func (pr *PostRepository) SelectPost(postID uint64) (*models.Post, error) {
	query := `select id, message, is_edited, created, parent_id, user_nickname, forum_slug, thread_id, path from posts where id = $1`
	post := &models.Post{}

	err := pr.db.QueryRow(query, postID).
		Scan(&post.ID,
			&post.Message,
			&post.IsEdited,
			&post.Created,
			&post.Parent,
			&post.Author,
			&post.ForumSlug,
			&post.ThreadID,
			pq.Array(&post.Path))

	if err != nil {
		return nil, err
	}

	return post, nil
}
