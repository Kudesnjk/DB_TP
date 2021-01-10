package thread

import (
	"github.com/Kudesnjk/DB_TP/internal/models"
	"github.com/Kudesnjk/DB_TP/internal/tools"
)

type ThreadUsecase interface {
	GetThreadInfo(slugOrID string) (*models.Thread, error)
	CreateThread(thread *models.Thread) error
	GetThreadsByForumSlug(slug string, qpm *tools.QPM) ([]*models.Thread, error)
	VoteThread(threadID int, nickname string, vote int) error
	UpdateThread(thread *models.Thread) error
}
