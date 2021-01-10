package service

import (
	"github.com/Kudesnjk/DB_TP/internal/models"
)

type ServiceRepository interface {
	SelectStatus() (*models.Service, error)
	TruncateTables() error
}
