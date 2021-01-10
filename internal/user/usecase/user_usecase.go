package usecase

import (
	"database/sql"

	"github.com/Kudesnjk/DB_TP/internal/models"
	"github.com/Kudesnjk/DB_TP/internal/user"
)

type UserUsecase struct {
	userRep user.UserRepository
}

func NewUserUsecase(userRep user.UserRepository) user.UserUsecase {
	return &UserUsecase{
		userRep: userRep,
	}
}

func (uu *UserUsecase) GetByNicknameOrEmail(nickname string, email string) ([]*models.User, error) {
	users, err := uu.userRep.SelectByNicknameOrEmail(nickname, email)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (uu *UserUsecase) CreateUser(user *models.User) error {
	err := uu.userRep.InsertUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (uu *UserUsecase) GetUserInfo(nickname string) (*models.User, error) {
	user, err := uu.userRep.SelectUserByNickname(nickname)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (uu *UserUsecase) UpdateUserInfo(user *models.User) error {
	err := uu.userRep.UpdateUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (uu *UserUsecase) CheckEmailExists(email string) (bool, error) {
	_, err := uu.userRep.SelectUserByEmail(email)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
