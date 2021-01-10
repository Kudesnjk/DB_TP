package forum

import (
	"github.com/Kudesnjk/DB_TP/internal/models"
	"github.com/Kudesnjk/DB_TP/internal/tools"
)

type ForumRepository interface {
	InsertForum(forum *models.Forum) error
	SelectByForumSlug(slug string) (*models.Forum, error)
	SelectForumUsers(slug string, qpm *tools.QPM) ([]*models.User, error)
}
