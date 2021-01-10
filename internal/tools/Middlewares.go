package tools

import (
	"github.com/labstack/echo/v4"
	"log"
	"time"
)

func TimingLogMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		start := time.Now()

		err := next(ctx)

		end := time.Now()

		log.Println()
		log.Println("Worktime: ", end.Sub(start))
		log.Println("Path: ", ctx.Request())

		return err
	}
}
