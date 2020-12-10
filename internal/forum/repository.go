package forum

import (
	"github.com/Kudesnjk/DB_TP/internal/models"
)

type ForumRepository interface {
	InsertForum(forum *models.Forum) error
	SelectByForumSlug(slug string) (*models.Forum, error)
	SelectForumUsers(slug string) ([]*models.User, error)
}
