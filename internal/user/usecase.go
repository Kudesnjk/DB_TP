package user

import (
	"github.com/Kudesnjk/DB_TP/internal/models"
)

type UserUsecase interface {
	GetByNicknameOrEmail(nickname string, email string) ([]*models.User, error)
	CheckEmailExists(email string) (bool, error)
	CreateUser(user *models.User) error
	GetUserInfo(nickname string) (*models.User, error)
	UpdateUserInfo(user *models.User) error
}
