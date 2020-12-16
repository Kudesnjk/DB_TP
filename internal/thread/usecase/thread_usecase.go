package usecase

import (
	"database/sql"

	"github.com/Kudesnjk/DB_TP/internal/tools"

	"github.com/Kudesnjk/DB_TP/internal/models"
	"github.com/Kudesnjk/DB_TP/internal/thread"
)

type ThreadUsecase struct {
	threadRep thread.ThreadRepository
}

func NewThreadUsecase(threadRep thread.ThreadRepository) thread.ThreadUsecase {
	return &ThreadUsecase{
		threadRep: threadRep,
	}
}

func (tu *ThreadUsecase) CreateThread(thread *models.Thread) error {
	err := tu.threadRep.InsertThread(thread)
	if err != nil {
		return err
	}
	return nil
}

func (tu *ThreadUsecase) GetThreadInfo(slugOrID string) (*models.Thread, error) {
	thread, err := tu.threadRep.SelectBySlugOrID(slugOrID)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return thread, nil
}

func (tu *ThreadUsecase) GetThreadsByForumSlug(slug string, qpm *tools.QPM) ([]*models.Thread, error) {
	threads, err := tu.threadRep.SelectThreadsByForumSlug(slug, qpm)
	if err != nil {
		return nil, err
	}
	return threads, nil
}

func (tu *ThreadUsecase) VoteThread(threadID int, nickname string, vote int) error {
	return tu.threadRep.InsertVote(nickname, threadID, vote)
}

func (tu *ThreadUsecase) UpdateThread(thread *models.Thread) error {
	return tu.threadRep.UpdateThread(thread)
}
