package command

import (
	"edupaim/xpto-support/app/services"
	"github.com/sirupsen/logrus"
)

type LegacyIntegrate interface {
	LegacyIntegrate(cmd IntegrateCmd) error
}

type IntegrateCmd struct {
}

type LegacyIntegrateController struct {
	legacyRepository services.LegacyRepository
	localRepository  services.LocalRepository
}

func NewLegacyIntegrateController(
	legacyRepo services.LegacyRepository,
	localRepo services.LocalRepository) *LegacyIntegrateController {
	return &LegacyIntegrateController{
		legacyRepository: legacyRepo,
		localRepository:  localRepo,
	}
}

func (controller *LegacyIntegrateController) LegacyIntegrate(cmd *IntegrateCmd) error {
	negatives, err := controller.legacyRepository.GetAllNegatives()
	if err != nil {
		logrus.WithError(err).Errorln(err.Error())
	}
	logrus.WithField("amount", len(negatives)).Debugln("get all negatives from legacy repository")
	for _, negative := range negatives {
		negative.DatesToUTC()
		err = controller.localRepository.SaveNegative(negative)
		if err != nil {
			return err
		}
	}
	return nil
}
