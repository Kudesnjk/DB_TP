package usecase

import (
	"github.com/Kudesnjk/DB_TP/internal/models"
	"github.com/Kudesnjk/DB_TP/internal/post"
	"github.com/Kudesnjk/DB_TP/internal/tools"
)

type PostUsecase struct {
	postRep post.PostRepository
}

func NewPostUsecase(postRep post.PostRepository) post.PostUsecase {
	return &PostUsecase{
		postRep: postRep,
	}
}

func (pu *PostUsecase) CreatePost(post *models.Post) error {
	err := pu.postRep.InsertPost(post)
	if err != nil {
		return err
	}
	return nil
}

func (pu *PostUsecase) GetPosts(threadID uint64, qpm *tools.QPM) ([]*models.Post, error) {
	posts, err := pu.postRep.SelectPosts(threadID, qpm)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
