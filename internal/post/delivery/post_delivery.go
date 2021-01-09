package delivery

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Kudesnjk/DB_TP/internal/forum"

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
	forumUsecase  forum.ForumUsecase
}

func NewPostDelivery(postUsecase post.PostUsecase,
	userUsecase user.UserUsecase,
	threadUsecase thread.ThreadUsecase,
	forumUsecase forum.ForumUsecase) *PostDelivery {
	return &PostDelivery{
		postUsecase:   postUsecase,
		userUsecase:   userUsecase,
		threadUsecase: threadUsecase,
		forumUsecase:  forumUsecase,
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

		if request.Message == "" || request.Message == post.Message {
			return ctx.JSON(http.StatusOK, post)
		}

		post = &models.Post{
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

		if len(posts) == 0 {
			return ctx.JSON(http.StatusCreated, posts)
		}

		location, _ := time.LoadLocation("UTC")
		now := time.Now().In(location).Round(time.Microsecond)

		for _, post := range posts {
			user, err := pd.userUsecase.GetUserInfo(post.Author)
			if err != nil {
				log.Println(err)
				return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
					Message: tools.ConstInternalErrorMessage,
				})
			}

			if user == nil {
				return ctx.JSON(http.StatusNotFound, tools.BadResponse{
					Message: tools.ConstSomeMessage,
				})
			}

			post.Created = now
			post.ThreadSlug = thread.Slug
			post.ForumSlug = thread.ForumSlug
			post.ThreadID = thread.ID

			err = pd.postUsecase.CreatePost(post)

			if err == tools.ErrorParentPostNotFound {
				return ctx.JSON(http.StatusConflict, tools.BadResponse{
					Message: tools.ConstSomeMessage,
				})
			}

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

		relatedParams := strings.Split(ctx.QueryParam("related"), ",")
		var userParam, threadParam, forumParam bool

		for _, param := range relatedParams {
			switch param {
			case "user":
				userParam = true
			case "thread":
				threadParam = true
			case "forum":
				forumParam = true
			}
		}

		res := &Response{
			Post: post,
		}

		if userParam {
			user, err := pd.userUsecase.GetUserInfo(post.Author)

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

			res.User = user
		}

		if threadParam {
			thread, err := pd.threadUsecase.GetThreadInfo(strconv.Itoa(int(post.ThreadID)))

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

			res.Thread = thread
		}

		if forumParam {
			forum, err := pd.forumUsecase.GetForumInfo(post.ForumSlug)

			if err != nil {
				log.Println(err)
				return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
					Message: tools.ConstInternalErrorMessage,
				})
			}

			if forum == nil {
				return ctx.JSON(http.StatusNotFound, tools.BadResponse{
					Message: tools.ConstNotFoundMessage,
				})
			}

			res.Forum = forum
		}

		return ctx.JSON(http.StatusOK, res)
	}
}
