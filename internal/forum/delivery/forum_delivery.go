package delivery

import (
	"log"
	"net/http"

	"github.com/Kudesnjk/DB_TP/internal/models"
	"github.com/Kudesnjk/DB_TP/internal/user"

	"github.com/Kudesnjk/DB_TP/internal/tools"

	"github.com/Kudesnjk/DB_TP/internal/forum"
	"github.com/labstack/echo/v4"
)

type ForumDelivery struct {
	forumUsecase forum.ForumUsecase
	userUsecase  user.UserUsecase
}

func NewForumDelivery(forumUsecase forum.ForumUsecase, userUsecase user.UserUsecase) *ForumDelivery {
	return &ForumDelivery{
		forumUsecase: forumUsecase,
		userUsecase:  userUsecase,
	}
}

func (fd *ForumDelivery) Configure(e *echo.Echo) {
	e.POST("api/forum/create", fd.CreateForumHandler())
	e.GET("api/forum/:slug/details", fd.GetForumHandler())
	e.GET("api/forum/:slug/users", fd.GetForumUsersHandler())
}

func (fd *ForumDelivery) GetForumUsersHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		slug := ctx.Param("slug")
		forum, err := fd.forumUsecase.GetForumInfo(slug)

		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		if forum == nil {
			return ctx.JSON(http.StatusNotFound, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		qpm := tools.NewQPM(ctx)
		users, err := fd.forumUsecase.GetForumUsers(slug, qpm)
		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		return ctx.JSON(http.StatusOK, users)
	}
}

func (fd *ForumDelivery) GetForumHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		slug := ctx.Param("slug")
		forum, err := fd.forumUsecase.GetForumInfo(slug)

		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		if forum == nil {
			return ctx.JSON(http.StatusNotFound, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		return ctx.JSON(http.StatusOK, forum)
	}
}

func (fd *ForumDelivery) CreateForumHandler() echo.HandlerFunc {
	type Request struct {
		Slug  string `json:"slug"`
		User  string `json:"user"`
		Title string `json:"title"`
	}

	type Response struct {
		Slug  string `json:"slug"`
		User  string `json:"user"`
		Title string `json:"title"`
	}

	return func(ctx echo.Context) error {
		request := &Request{}
		err := ctx.Bind(request)

		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		user, err := fd.userUsecase.GetUserInfo(request.User)

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

		forum, err := fd.forumUsecase.GetForumInfo(request.Slug)

		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		if forum != nil {
			return ctx.JSON(http.StatusConflict, forum)
		}

		newForum := &models.Forum{
			User:  user.Nickname,
			Title: request.Title,
			Slug:  request.Slug,
		}

		err = fd.forumUsecase.CreateForum(newForum)
		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		response := &Response{
			User:  newForum.User,
			Title: newForum.Title,
			Slug:  newForum.Slug,
		}
		return ctx.JSON(http.StatusCreated, response)
	}
}
