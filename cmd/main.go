package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	forumRepository "github.com/Kudesnjk/DB_TP/internal/forum/repository"
	postRepository "github.com/Kudesnjk/DB_TP/internal/post/repository"
	serviceRepository "github.com/Kudesnjk/DB_TP/internal/service/repository"
	threadRepository "github.com/Kudesnjk/DB_TP/internal/thread/repository"
	userRepository "github.com/Kudesnjk/DB_TP/internal/user/repository"

	forumUsecase "github.com/Kudesnjk/DB_TP/internal/forum/usecase"
	postUsecase "github.com/Kudesnjk/DB_TP/internal/post/usecase"
	serviceUsecase "github.com/Kudesnjk/DB_TP/internal/service/usecase"
	threadUsecase "github.com/Kudesnjk/DB_TP/internal/thread/usecase"
	userUsecase "github.com/Kudesnjk/DB_TP/internal/user/usecase"

	forumDelivery "github.com/Kudesnjk/DB_TP/internal/forum/delivery"
	postDelivery "github.com/Kudesnjk/DB_TP/internal/post/delivery"
	serviceDelivery "github.com/Kudesnjk/DB_TP/internal/service/delivery"
	threadDelivery "github.com/Kudesnjk/DB_TP/internal/thread/delivery"
	userDelivery "github.com/Kudesnjk/DB_TP/internal/user/delivery"
)

func main() {
	dbConn, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "postgres", "postgres", "forum"))
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = dbConn.Close()
	}()

	if err := dbConn.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Printf("DB connected")

	userRep := userRepository.NewUserRepository(dbConn)
	forumRep := forumRepository.NewForumRepository(dbConn)
	threadRep := threadRepository.NewThreadRepository(dbConn)
	postRep := postRepository.NewPostRepository(dbConn)
	serviceRep := serviceRepository.NewServiceRepository(dbConn)

	userUsecase := userUsecase.NewUserUsecase(userRep)
	forumUsecase := forumUsecase.NewForumUsecase(forumRep)
	threadUsecase := threadUsecase.NewThreadUsecase(threadRep)
	postUsecase := postUsecase.NewPostUsecase(postRep)
	serviceUsecase := serviceUsecase.NewServiceUsecase(serviceRep)

	userDelivery := userDelivery.NewUserDelivery(userUsecase)
	forumDelivery := forumDelivery.NewForumDelivery(forumUsecase, userUsecase)
	threadDelivery := threadDelivery.NewThreadDelivery(threadUsecase, userUsecase, forumUsecase)
	postDelivery := postDelivery.NewPostDelivery(postUsecase, userUsecase, threadUsecase, forumUsecase)
	serviceDelivery := serviceDelivery.NewServiceDelivery(serviceUsecase)

	e := echo.New()
	//e.Use(tools.TimingLogMiddleware)

	userDelivery.Configure(e)
	forumDelivery.Configure(e)
	threadDelivery.Configure(e)
	postDelivery.Configure(e)
	serviceDelivery.Configure(e)
	e.Logger.Fatal(e.Start("0.0.0.0:5000"))
}
