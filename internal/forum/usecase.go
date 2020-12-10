package forum

import (
	"github.com/Kudesnjk/DB_TP/internal/models"
)

type ForumUsecase interface {
	CreateForum(forum *models.Forum) error
	GetForumInfo(slug string) (*models.Forum, error)
	GetForumUsers(slug string) ([]*models.User, error)
}
