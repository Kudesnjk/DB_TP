package post

import (
	"github.com/Kudesnjk/DB_TP/internal/models"
	"github.com/Kudesnjk/DB_TP/internal/tools"
)

type PostRepository interface {
	InsertPost(post *models.Post) error
	SelectPosts(threadID uint64, qpm *tools.QPM) ([]*models.Post, error)
}
