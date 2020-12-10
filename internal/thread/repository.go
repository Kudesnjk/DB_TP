package thread

import (
	"github.com/Kudesnjk/DB_TP/internal/models"
	"github.com/Kudesnjk/DB_TP/internal/tools"
)

type ThreadRepository interface {
	SelectBySlugOrID(slugOrID string) (*models.Thread, error)
	InsertThread(thread *models.Thread) error
	SelectThreadsByForumSlug(slug string, qpm *tools.QPM) ([]*models.Thread, error)
}
