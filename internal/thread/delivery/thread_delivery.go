package delivery

import (
	"github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Kudesnjk/DB_TP/internal/models"

	"github.com/Kudesnjk/DB_TP/internal/forum"

	"github.com/Kudesnjk/DB_TP/internal/tools"

	"github.com/Kudesnjk/DB_TP/internal/user"

	"github.com/Kudesnjk/DB_TP/internal/thread"
	"github.com/labstack/echo/v4"
)

type ThreadDelivery struct {
	threadUsecase thread.ThreadUsecase
	userUsecase   user.UserUsecase
	forumUsecase  forum.ForumUsecase
}

func NewThreadDelivery(threadUsecase thread.ThreadUsecase, userUsecase user.UserUsecase, forumUsecase forum.ForumUsecase) *ThreadDelivery {
	return &ThreadDelivery{
		threadUsecase: threadUsecase,
		userUsecase:   userUsecase,
		forumUsecase:  forumUsecase,
	}
}

func (td *ThreadDelivery) Configure(e *echo.Echo) {
	e.POST("api/forum/:slug/create", td.CreateThreadHandler())
	e.GET("api/forum/:slug/threads", td.GetForumThreadsHandler())
	e.GET("api/thread/:slug_or_id/details", td.GetConcreteThreadHandler())
	e.POST("api/thread/:slug_or_id/details", td.UpdateThreadHandler())
	e.POST("api/thread/:slug_or_id/vote", td.VoteThreadHandler())
}

func (td *ThreadDelivery) CreateThreadHandler() echo.HandlerFunc {
	type Request struct {
		Author  string    `json:"author"`
		Slug    string    `json:"slug"`
		Created time.Time `json:"created"`
		Message string    `json:"message"`
		Title   string    `json:"title"`
		Forum   string    `json:"forum"`
	}

	return func(ctx echo.Context) error {
		slug := ctx.Param("slug")
		request := &Request{}
		err := ctx.Bind(request)

		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		forum, err := td.forumUsecase.GetForumInfo(slug)
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

		user, err := td.userUsecase.GetUserInfo(request.Author)
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

		if request.Slug != "" {
			thread, err := td.threadUsecase.GetThreadInfo(request.Slug)
			if err != nil {
				log.Println(err)
				return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
					Message: tools.ConstInternalErrorMessage,
				})
			}

			if thread != nil {
				return ctx.JSON(http.StatusConflict, thread)
			}
		}

		thread := &models.Thread{
			Author:    request.Author,
			Created:   request.Created,
			Message:   request.Message,
			Title:     request.Title,
			ForumSlug: forum.Slug,
		}

		if request.Slug != thread.ForumSlug {
			thread.Slug = request.Slug
		}

		err = td.threadUsecase.CreateThread(thread)
		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		return ctx.JSON(http.StatusCreated, thread)
	}
}

func (td *ThreadDelivery) GetForumThreadsHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		slug := ctx.Param("slug")
		qpm := tools.NewQPM(ctx)

		forum, err := td.forumUsecase.GetForumInfo(slug)

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

		threads, err := td.threadUsecase.GetThreadsByForumSlug(slug, qpm)

		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		return ctx.JSON(http.StatusOK, threads)
	}
}

func (td *ThreadDelivery) GetConcreteThreadHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		slugOrID := ctx.Param("slug_or_id")
		thread, err := td.threadUsecase.GetThreadInfo(slugOrID)

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

		return ctx.JSON(http.StatusOK, thread)
	}
}

func (td *ThreadDelivery) UpdateThreadHandler() echo.HandlerFunc {
	type Request struct {
		Message string `json:"message"`
		Title   string `json:"title"`
	}

	return func(ctx echo.Context) error {
		slugOrID := ctx.Param("slug_or_id")
		request := &Request{}
		err := ctx.Bind(request)

		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		thread, err := td.threadUsecase.GetThreadInfo(slugOrID)

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

		if request.Title == "" && request.Message == "" {
			return ctx.JSON(http.StatusOK, thread)
		}

		if request.Title != "" {
			thread.Title = request.Title
		}

		if request.Message != "" {
			thread.Message = request.Message
		}

		err = td.threadUsecase.UpdateThread(thread)

		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		return ctx.JSON(http.StatusOK, thread)
	}
}

func (td *ThreadDelivery) VoteThreadHandler() echo.HandlerFunc {
	type Request struct {
		Nickname string `json:"nickname"`
		Voice    int    `json:"voice"`
	}
	return func(ctx echo.Context) error {
		slugOrID := ctx.Param("slug_or_id")
		request := &Request{}
		err := ctx.Bind(request)

		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		threadID, isSlug := strconv.Atoi(slugOrID)
		var threadSluged *models.Thread

		if isSlug != nil {
			threadSluged, err = td.threadUsecase.GetThreadInfo(slugOrID)

			if err != nil {
				log.Println(err)
				return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
					Message: tools.ConstInternalErrorMessage,
				})
			}

			if threadSluged == nil {
				return ctx.JSON(http.StatusNotFound, tools.BadResponse{
					Message: tools.ConstNotFoundMessage,
				})
			}
		}

		if isSlug != nil {
			err = td.threadUsecase.VoteThread(int(threadSluged.ID), request.Nickname, request.Voice)
		} else {
			err = td.threadUsecase.VoteThread(threadID, request.Nickname, request.Voice)
		}

		if err != nil {
			if err.(*pq.Error).Code == tools.ConstUserNotFoundError {
				return ctx.JSON(http.StatusNotFound, tools.BadResponse{
					Message: tools.ConstSomeMessage,
				})
			}

			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		threadSluged, err = td.threadUsecase.GetThreadInfo(slugOrID)

		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		if threadSluged == nil {
			return ctx.JSON(http.StatusNotFound, tools.BadResponse{
				Message: tools.ConstNotFoundMessage,
			})
		}

		return ctx.JSON(http.StatusOK, threadSluged)
	}
}
