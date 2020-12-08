package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	userRepository "github.com/Kudesnjk/DB_TP/internal/user/repository"
	"github.com/labstack/echo/v4"

	userUsecase "github.com/Kudesnjk/DB_TP/internal/user/usecase"

	userDelivery "github.com/Kudesnjk/DB_TP/internal/user/delivery"
)

func main() {
	dbConn, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "db_forum_user", "9680nM", "db_forum"))
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

	userUsecase := userUsecase.NewUserUsecase(userRep)

	userDelivery := userDelivery.NewUserDelivery(userUsecase)

	e := echo.New()

	userDelivery.Configure(e)
	e.Logger.Fatal(e.Start("127.0.0.1:8080"))
}
