package post

import (
	"github.com/Kudesnjk/DB_TP/internal/models"
	"github.com/Kudesnjk/DB_TP/internal/tools"
)

type PostUsecase interface {
	CreatePost(post *models.Post) error
	GetPosts(threadID uint64, qpm *tools.QPM) ([]*models.Post, error)
}
