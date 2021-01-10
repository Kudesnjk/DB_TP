package repository

import (
	"database/sql"
	"fmt"
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

func (pr *PostRepository) InsertPost(posts []*models.Post, ad *models.AdditionalPostData) error {
	tx, err := pr.db.Begin()
	queryStr := `insert into posts(id, message, is_edited, created, parent_id, user_nickname, thread_id, forum_slug, path)
			values(default, $1, default, $2, $3, $4, $5, $6, $7) returning id`

	for _, post := range posts {
		err = tx.QueryRow(queryStr,
			post.Message,
			ad.Created,
			post.Parent,
			post.Author,
			ad.ThreadID,
			ad.ForumSlug,
			pq.Array([]int64{int64(post.Parent)})).
			Scan(
				&post.ID,
			)

		if err != nil {
			tx.Rollback()
			return err
		}

		post.Created = ad.Created
		post.ThreadID = ad.ThreadID
		post.ForumSlug = ad.ForumSlug
		post.ThreadSlug = ad.ThreadSlug
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
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

	if len(posts) > 0 && posts[0].ID == 813277 {
		fmt.Println()
		fmt.Println(qpm.Sort)
		fmt.Println(posts[0])
		fmt.Println(posts[0].ThreadID)
		fmt.Println(qpm)
		fmt.Println(query)
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
