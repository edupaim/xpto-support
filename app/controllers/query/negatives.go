package query

import (
	"context"
	"edupaim/xpto-support/app/domain"
	"edupaim/xpto-support/app/services"
	"github.com/sirupsen/logrus"
)

type NegativesQuery interface {
	GetByQuery(customerDocument map[string][]string, ctx context.Context) ([]domain.Negative, error)
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

func (controller *NegativeQueryController) GetByQuery(customerDocument map[string][]string, ctx context.Context) ([]domain.Negative, error) {
	for key, value := range customerDocument {
		if key == "customerDocument" {
			cd := domain.CryptDocument(value[0])
			err := cd.Encrypt()
			if err != nil {
				logrus.WithError(err).Errorln("crypt customer document")
				return nil, err
			}
			customerDocument[key][0] = string(cd)
		}
		if key == "companyDocument" {
			cd := domain.CryptDocument(value[0])
			err := cd.Encrypt()
			if err != nil {
				logrus.WithError(err).Errorln("crypt company document")
				return nil, err
			}
			customerDocument[key][0] = string(cd)
		}
	}
	negatives, err := controller.localRepository.GetNegativeByQuery(customerDocument, ctx)
	if err != nil {
		return nil, err
	}
	for i, _ := range negatives {
		err = negatives[i].DecryptDocuments()
		if err != nil {
			logrus.WithError(err).Errorln("decrypt customer document from negative")
			return nil, err
		}
	}
	return negatives, nil
}
