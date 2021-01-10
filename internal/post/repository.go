package post

import (
	"github.com/Kudesnjk/DB_TP/internal/models"
	"github.com/Kudesnjk/DB_TP/internal/tools"
)

type PostRepository interface {
	InsertPost(posts []*models.Post, ad *models.AdditionalPostData) error
	SelectPosts(threadID uint64, qpm *tools.QPM) ([]*models.Post, error)
	SelectPost(postID uint64) (*models.Post, error)
	UpdatePost(post *models.Post) error
}
