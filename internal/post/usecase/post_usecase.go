package usecase

import (
	"database/sql"
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

func (pu *PostUsecase) CreatePost(posts []*models.Post, ad *models.AdditionalPostData) error {
	err := pu.postRep.InsertPost(posts, ad)

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

func (pu *PostUsecase) UpdatePost(post *models.Post) error {
	return pu.postRep.UpdatePost(post)
}

func (pu *PostUsecase) GetPost(postID uint64) (*models.Post, error) {
	post, err := pu.postRep.SelectPost(postID)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return post, nil
}
