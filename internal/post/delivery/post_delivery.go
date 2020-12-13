package delivery

import (
	"log"
	"net/http"

	"github.com/Kudesnjk/DB_TP/internal/models"
	"github.com/Kudesnjk/DB_TP/internal/user"

	"github.com/Kudesnjk/DB_TP/internal/tools"

	"github.com/Kudesnjk/DB_TP/internal/post"
	"github.com/labstack/echo/v4"
)

type PostDelivery struct {
	postUsecase post.PostUsecase
	userUsecase user.UserUsecase
}

func NewPostDelivery(postUsecase post.PostUsecase, userUsecase user.UserUsecase) *PostDelivery {
	return &PostDelivery{
		postUsecase: postUsecase,
		userUsecase: userUsecase,
	}
}

func (fd *PostDelivery) Configure(e *echo.Echo) {
	e.POST("api/thread/:slug_or_id/create", fd.CreatePostHandler())
	e.GET("api/post/:slug/details", fd.GetPostHandler())
	e.GET("api/post/:slug/users", fd.GetPostUsersHandler())
}

func (fd *PostDelivery) GetPostUsersHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		slug := ctx.Param("slug")
		post, err := fd.postUsecase.GetPostInfo(slug)

		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		if post == nil {
			log.Println(err)
			return ctx.JSON(http.StatusNotFound, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		users, err := fd.postUsecase.GetPostUsers(slug)
		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		return ctx.JSON(http.StatusOK, users)
	}
}

func (fd *PostDelivery) GetPostHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		slug := ctx.Param("slug")
		post, err := fd.postUsecase.GetPostInfo(slug)

		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		if post == nil {
			log.Println(err)
			return ctx.JSON(http.StatusNotFound, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		return ctx.JSON(http.StatusOK, post)
	}
}

func (fd *PostDelivery) CreatePostHandler() echo.HandlerFunc {
	type Request struct {
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
			log.Println(err)
			return ctx.JSON(http.StatusNotFound, tools.BadResponse{
				Message: tools.ConstNotFoundMessage,
			})
		}

		post, err := fd.postUsecase.GetPostInfo(request.Slug)

		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		if post != nil {
			log.Println(err)
			return ctx.JSON(http.StatusConflict, tools.BadResponse{
				Message: tools.ConstNotFoundMessage,
			})
		}

		newPost := &models.Post{
			User:  request.User,
			Title: request.Title,
			Slug:  request.Slug,
		}

		err = fd.postUsecase.CreatePost(newPost)
		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}
		return ctx.JSON(http.StatusCreated, newPost)
	}
}
