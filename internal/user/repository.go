package user

import (
	"github.com/Kudesnjk/DB_TP/internal/models"
)

type UserRepository interface {
	SelectByNicknameOrEmail(nickname string, email string) ([]*models.User, error)
	InsertUser(user *models.User) error
	SelectUserByNickname(nickname string) (*models.User, error)
	SelectUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error
}
