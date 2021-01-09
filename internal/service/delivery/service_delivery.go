package delivery

import (
	"log"
	"net/http"

	"github.com/Kudesnjk/DB_TP/internal/tools"

	"github.com/Kudesnjk/DB_TP/internal/service"
	"github.com/labstack/echo/v4"
)

type ServiceDelivery struct {
	serviceUsecase service.ServiceUsecase
}

func NewServiceDelivery(serviceUsecase service.ServiceUsecase) *ServiceDelivery {
	return &ServiceDelivery{
		serviceUsecase: serviceUsecase,
	}
}

func (sd *ServiceDelivery) Configure(e *echo.Echo) {
	e.GET("api/service/status", sd.ServiceStatusHandler())
	e.POST("api/service/clear", sd.ClearServiceHandler())
}

func (sd *ServiceDelivery) ServiceStatusHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		service, err := sd.serviceUsecase.GetStatus()

		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		return ctx.JSON(http.StatusOK, service)
	}
}

func (sd *ServiceDelivery) ClearServiceHandler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		err := sd.serviceUsecase.ClearDatabase()

		if err != nil {
			log.Println(err)
			return ctx.JSON(http.StatusInternalServerError, tools.BadResponse{
				Message: tools.ConstInternalErrorMessage,
			})
		}

		return ctx.JSON(http.StatusOK, tools.BadResponse{
			Message: tools.ConstSomeMessage,
		})
	}
}
