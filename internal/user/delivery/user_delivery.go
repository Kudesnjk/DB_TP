package delivery

import (
	"log"
	"net/http"

	"github.com/Kudesnjk/DB_TP/internal/models"

	"github.com/Kudesnjk/DB_TP/internal/tools"

	"github.com/Kudesnjk/DB_TP/internal/user"
	"github.com/labstack/echo/v4"
)

type UserDelivery struct {
	userUsecase user.UserUsecase
}

func NewUserDelivery(userUsecase user.UserUsecase) *UserDelivery {
	return &UserDelivery{
		userUsecase: userUsecase,
	}
}

func (ud *UserDelivery) Configure(e *echo.Echo) {
	e.POST("api/user/:nickname/create", ud.CreateUserHandler())
	e.GET("api/user/:nickname/profile", ud.GetUserHandler())
	e.POST("api/user/:nickname/profile", ud.UpdateUserHandler())
}

func (ud *UserDelivery) CreateUserHandler() echo.HandlerFunc {
	type Request struct {
		About    string `json:"about"`
		Email    string `json:"email"`
		FullName string `json:"fullname"`
	}

	return func(ctx echo.Context) error {
		nickname := ctx.Param("nickname")

		request := &Request{}

		err := ctx.Bind(request)

		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		users, err := ud.userUsecase.GetByNicknameOrEmail(nickname, request.Email)
		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		if len(users) > 0 {
			return ctx.JSON(http.StatusConflict, users)
		}

		user := &models.User{
			Nickname: nickname,
			Email:    request.Email,
			Fullname: request.FullName,
			About:    request.About,
		}

		err = ud.userUsecase.CreateUser(user)
		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}
		return ctx.JSON(http.StatusCreated, user)
	}
}

func (ud *UserDelivery) GetUserHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		nickname := ctx.Param("nickname")

		user, err := ud.userUsecase.GetUserInfo(nickname)

		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		if user == nil {
			return ctx.JSON(http.StatusNotFound, tools.BadResponse{
				Message: tools.ConstNotFoundMessage,
			})
		}

		return ctx.JSON(http.StatusOK, user)
	}
}

func (ud *UserDelivery) UpdateUserHandler() echo.HandlerFunc {
	type Request struct {
		About    string `json:"about"`
		Email    string `json:"email"`
		FullName string `json:"fullname"`
	}

	return func(ctx echo.Context) error {
		nickname := ctx.Param("nickname")

		request := &Request{}

		err := ctx.Bind(request)

		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		user, err := ud.userUsecase.GetUserInfo(nickname)

		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		if user == nil {
			log.Println(err)
			return ctx.JSON(http.StatusNotFound, tools.BadResponse{
				Message: tools.ConstNotFoundMessage,
			})
		}

		if request.About == "" && request.FullName == "" && request.Email == "" {
			return ctx.JSON(http.StatusOK, user)
		}

		exists, err := ud.userUsecase.CheckEmailExists(request.Email)
		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		if exists {
			return ctx.JSON(http.StatusConflict, tools.BadResponse{
				Message: tools.ConstSomeMessage,
			})
		}

		if request.About != "" {
			user.About = request.About
		}
		if request.FullName != "" {
			user.Fullname = request.FullName
		}
		if request.Email != "" {
			user.Email = request.Email
		}

		err = ud.userUsecase.UpdateUserInfo(user)
		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		return ctx.JSON(http.StatusOK, user)
	}
}
