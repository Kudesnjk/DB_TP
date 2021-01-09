package service

import (
	"github.com/Kudesnjk/DB_TP/internal/models"
)

type ServiceUsecase interface {
	GetStatus() (*models.Service, error)
	ClearDatabase() error
}
