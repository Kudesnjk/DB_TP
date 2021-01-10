package usecase

import (
	"github.com/Kudesnjk/DB_TP/internal/models"
	"github.com/Kudesnjk/DB_TP/internal/service"
)

type ServiceUsecase struct {
	serviceRep service.ServiceRepository
}

func NewServiceUsecase(serviceRep service.ServiceRepository) service.ServiceUsecase {
	return &ServiceUsecase{
		serviceRep: serviceRep,
	}
}

func (su *ServiceUsecase) GetStatus() (*models.Service, error) {
	return su.serviceRep.SelectStatus()
}

func (su *ServiceUsecase) ClearDatabase() error {
	return su.serviceRep.TruncateTables()
}
