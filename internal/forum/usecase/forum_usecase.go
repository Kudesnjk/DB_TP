package usecase

import (
	"database/sql"

	"github.com/Kudesnjk/DB_TP/internal/tools"

	"github.com/Kudesnjk/DB_TP/internal/forum"
	"github.com/Kudesnjk/DB_TP/internal/models"
)

type ForumUsecase struct {
	forumRep forum.ForumRepository
}

func NewForumUsecase(forumRep forum.ForumRepository) forum.ForumUsecase {
	return &ForumUsecase{
		forumRep: forumRep,
	}
}

func (fu *ForumUsecase) CreateForum(forum *models.Forum) error {
	err := fu.forumRep.InsertForum(forum)
	if err != nil {
		return err
	}
	return nil
}

func (fu *ForumUsecase) GetForumInfo(slug string) (*models.Forum, error) {
	forum, err := fu.forumRep.SelectByForumSlug(slug)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return forum, nil
}

func (fu *ForumUsecase) GetForumUsers(slug string, qpm *tools.QPM) ([]*models.User, error) {
	users, err := fu.forumRep.SelectForumUsers(slug, qpm)
	if err != nil {
		return nil, err
	}
	return users, nil
}
