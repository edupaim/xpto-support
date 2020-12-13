package command

import (
	"context"
	"edupaim/xpto-support/app/services"
	"github.com/sirupsen/logrus"
)

type LegacyIntegrate interface {
	LegacyIntegrate(cmd *IntegrateCmd, ctx context.Context) error
}

type IntegrateCmd struct {
}

type LegacyIntegrateController struct {
	legacyRepository services.LegacyRepository
	localRepository  services.LocalRepository
}

func NewLegacyIntegrateController(
	legacyRepo services.LegacyRepository,
	localRepo services.LocalRepository,
) *LegacyIntegrateController {
	return &LegacyIntegrateController{
		legacyRepository: legacyRepo,
		localRepository:  localRepo,
	}
}

func (controller *LegacyIntegrateController) LegacyIntegrate(cmd *IntegrateCmd, ctx context.Context) error {
	negatives, err := controller.legacyRepository.GetAllNegatives(ctx)
	if err != nil {
		return err
	}
	logrus.WithField("amount", len(negatives)).Debugln("get all negatives from legacy repository")
	for _, negative := range negatives {
		negative.DatesToUTC()
		err = negative.EncryptDocuments()
		if err != nil {
			return err
		}
		err = controller.localRepository.SaveNegative(negative, ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
