package query

import (
	"edupaim/xpto-support/app/domain"
	"edupaim/xpto-support/app/services"
	"github.com/sirupsen/logrus"
)

type NegativesQuery interface {
	GetByCustomerDocument(customerDocument string) (*domain.Negative, error)
}

type NegativeQueryController struct {
	localRepository services.LocalRepository
}

func NewNegativeQueryController(
	localRepo services.LocalRepository,
) *NegativeQueryController {
	return &NegativeQueryController{
		localRepository: localRepo,
	}
}

func (controller *NegativeQueryController) GetByCustomerDocument(customerDocument string) ([]domain.Negative, error) {
	cd := domain.CustomerDocument(customerDocument)
	cd, err := cd.Encrypt()
	if err != nil {
		logrus.WithError(err).Errorln("crypt customer document")
		return nil, err
	}
	negatives, err := controller.localRepository.GetNegativeByCustomerDocument(string(cd))
	if err != nil {
		return nil, err
	}
	for i, _ := range negatives {
		err = negatives[i].DecryptCustomerDocument()
		if err != nil {
			logrus.WithError(err).Errorln("decrypt customer document from negative")
			return nil, err
		}
	}
	return negatives, nil
}
