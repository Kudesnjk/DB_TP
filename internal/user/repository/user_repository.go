package repository

import (
	"database/sql"

	"github.com/Kudesnjk/DB_TP/internal/models"

	"github.com/Kudesnjk/DB_TP/internal/user"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) user.UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) SelectByNicknameOrEmail(nickname string, email string) ([]*models.User, error) {
	rows, err := ur.db.Query("select fullname, nickname, email, about from users where lower(nickname) = lower($1) or lower(email) = lower($2)",
		nickname, email)

	if err != nil {
		return nil, err
	}

	users := make([]*models.User, 0)

	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(
			&user.Fullname,
			&user.Nickname,
			&user.Email,
			&user.About,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (ur *UserRepository) InsertUser(user *models.User) error {
	_, err := ur.db.Exec("insert into users (nickname, email, fullname, about) values($1, $2, $3, $4)",
		user.Nickname,
		user.Email,
		user.Fullname,
		user.About)

	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) SelectUserByNickname(nickname string) (*models.User, error) {
	user := &models.User{}

	err := ur.db.QueryRow("select fullname, nickname, email, about from users where lower(nickname) = lower($1)", nickname).Scan(
		&user.Fullname,
		&user.Nickname,
		&user.Email,
		&user.About,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) SelectUserByEmail(email string) (*models.User, error) {
	user := &models.User{}

	err := ur.db.QueryRow("select fullname, nickname, email, about from users where lower(email) = lower($1)", email).Scan(
		&user.Fullname,
		&user.Nickname,
		&user.Email,
		&user.About,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *UserRepository) UpdateUser(user *models.User) error {
	_, err := ur.db.Exec("update users set fullname=$1, about=$2, email=$3 where nickname=$4",
		user.Fullname,
		user.About,
		user.Email,
		user.Nickname,
	)

	if err != nil {
		return err
	}

	return nil
}
