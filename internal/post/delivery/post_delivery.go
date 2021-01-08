package delivery

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Kudesnjk/DB_TP/internal/models"

	"github.com/Kudesnjk/DB_TP/internal/thread"

	"github.com/Kudesnjk/DB_TP/internal/user"

	"github.com/Kudesnjk/DB_TP/internal/tools"

	"github.com/Kudesnjk/DB_TP/internal/post"
	"github.com/labstack/echo/v4"
)

type PostDelivery struct {
	postUsecase   post.PostUsecase
	userUsecase   user.UserUsecase
	threadUsecase thread.ThreadUsecase
}

func NewPostDelivery(postUsecase post.PostUsecase, userUsecase user.UserUsecase, threadUsecase thread.ThreadUsecase) *PostDelivery {
	return &PostDelivery{
		postUsecase:   postUsecase,
		userUsecase:   userUsecase,
		threadUsecase: threadUsecase,
	}
}

func (pd *PostDelivery) Configure(e *echo.Echo) {
	e.POST("api/thread/:slug_or_id/create", pd.CreatePostHandler())
	e.GET("api/thread/:slug_or_id/posts", pd.GetPostsHandler())
	e.GET("api/post/:id/details", pd.GetPostDetailsHandler())
	e.POST("api/post/:id/details", pd.UpdatePostHandler())
}

func (pd *PostDelivery) UpdatePostHandler() echo.HandlerFunc {
	type Request struct {
		Message string `json:"message"`
	}

	return func(ctx echo.Context) error {
		postID, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		request := &Request{}
		err = ctx.Bind(request)
		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		post := &models.Post{
			ID:      uint64(postID),
			Message: request.Message,
		}

		err = pd.postUsecase.UpdatePost(post)

		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, tools.BadResponse{
				Message: tools.ConstNotFoundMessage,
			})
		}

		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		return ctx.JSON(http.StatusOK, post)
	}
}

func (pd *PostDelivery) GetPostsHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		slugOrID := ctx.Param("slug_or_id")

		thread, err := pd.threadUsecase.GetThreadInfo(slugOrID)

		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		if thread == nil {
			return ctx.JSON(http.StatusNotFound, tools.BadResponse{
				Message: tools.ConstNotFoundMessage,
			})
		}

		qpm := tools.NewQPM(ctx)
		posts, err := pd.postUsecase.GetPosts(thread.ID, qpm)

		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		return ctx.JSON(http.StatusOK, posts)
	}
}

func (pd *PostDelivery) CreatePostHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		slugOrID := ctx.Param("slug_or_id")
		posts := make([]*models.Post, 0)

		result, err := ioutil.ReadAll(ctx.Request().Body)
		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		err = json.Unmarshal(result, &posts)
		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		if len(posts) == 0 {
			return ctx.JSON(http.StatusCreated, posts)
		}

		thread, err := pd.threadUsecase.GetThreadInfo(slugOrID)

		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		if thread == nil {
			return ctx.JSON(http.StatusNotFound, tools.BadResponse{
				Message: tools.ConstNotFoundMessage,
			})
		}

		now := time.Now()

		for _, post := range posts {
			post.Created = now
			post.ThreadSlug = thread.Slug
			post.ForumSlug = thread.ForumSlug
			post.ThreadID = thread.ID
			err := pd.postUsecase.CreatePost(post)
			if err != nil {
				log.Println(err)
				return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
					Message: tools.ConstInternalErrorMessage,
				})
			}
		}

		return ctx.JSON(http.StatusCreated, posts)
	}
}

func (pd *PostDelivery) GetPostDetailsHandler() echo.HandlerFunc {
	type Response struct {
		Post   *models.Post   `json:"post"`
		User   *models.User   `json:"author,omitempty"`
		Thread *models.Thread `json:"thread,omitempty"`
		Forum  *models.Forum  `json:"forum,omitempty"`
	}

	return func(ctx echo.Context) error {
		postID, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		post, err := pd.postUsecase.GetPost(uint64(postID))

		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		if post == nil {
			return ctx.JSON(http.StatusNotFound, tools.BadResponse{
				Message: tools.ConstNotFoundMessage,
			})
		}

		res := &Response{
			Post: post,
		}

		return ctx.JSON(http.StatusOK, res)
	}
}
