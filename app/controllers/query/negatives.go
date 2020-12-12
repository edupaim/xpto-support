package query

import (
	"edupaim/xpto-support/app/domain"
	"edupaim/xpto-support/app/services"
)

type NegativesQuery interface {
	GetByCustomerDocument(customerDocument string) (*domain.Negative, error)
}

type NegativeQueryController struct {
	localRepository services.LocalRepository
}

func NewNegativeQueryController(
	localRepo services.LocalRepository) *NegativeQueryController {
	return &NegativeQueryController{
		localRepository: localRepo,
	}
}

func (controller *NegativeQueryController) GetByCustomerDocument(customerDocument string) (*domain.Negative, error) {
	return controller.localRepository.GetNegativeByCustomerDocument(customerDocument)
}
