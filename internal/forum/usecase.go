package forum

import (
	"github.com/Kudesnjk/DB_TP/internal/models"
	"github.com/Kudesnjk/DB_TP/internal/tools"
)

type ForumUsecase interface {
	CreateForum(forum *models.Forum) error
	GetForumInfo(slug string) (*models.Forum, error)
	GetForumUsers(slug string, qpm *tools.QPM) ([]*models.User, error)
}
