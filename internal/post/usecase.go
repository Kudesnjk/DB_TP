package post

import (
	"github.com/Kudesnjk/DB_TP/internal/models"
	"github.com/Kudesnjk/DB_TP/internal/tools"
)

type PostUsecase interface {
	CreatePost(posts []*models.Post, ad *models.AdditionalPostData) error
	GetPosts(threadID uint64, qpm *tools.QPM) ([]*models.Post, error)
	GetPost(postID uint64) (*models.Post, error)
	UpdatePost(post *models.Post) error
}
