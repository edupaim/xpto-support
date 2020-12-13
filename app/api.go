package app

import (
	"context"
	"edupaim/xpto-support/app/controllers/command"
	"edupaim/xpto-support/app/controllers/query"
	"edupaim/xpto-support/app/services"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Api struct {
	httpServer  *http.Server
	services    AppServices
	controllers ApiControllers
}

type AppServices struct {
	legacyRepository services.LegacyRepository
	localRepository  services.LocalRepository
}

type ApiControllers struct {
	legacyIntegrate command.LegacyIntegrate
	negativesQuery  query.NegativesQuery
}

func InitializeApi(c *Config) (*Api, error) {
	api := &Api{}
	err := api.initializeServices(c)
	if err != nil {
		return nil, err
	}
	err = api.initializeControllers()
	if err != nil {
		return nil, err
	}
	r := initializeApiRouter(api.controllers)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", c.ServerConfig.Port),
		Handler: r,
	}
	api.httpServer = srv
	return api, nil
}

func (api *Api) initializeControllers() error {
	legacyController := command.NewLegacyIntegrateController(api.services.legacyRepository, api.services.localRepository)
	api.controllers.legacyIntegrate = legacyController
	negativeController := query.NewNegativeQueryController(api.services.localRepository)
	api.controllers.negativesQuery = negativeController
	return nil
}

func (api *Api) initializeServices(c *Config) error {
	legacyRepository, err := services.InitializeApiLegacyRepository(c.LegacyXptoConfig.Address)
	if err != nil {
		return err
	}
	api.services.legacyRepository = legacyRepository
	localRepository, err := services.InitializeArangoLocalRepository(c.ArangoConfig)
	if err != nil {
		return err
	}
	api.services.localRepository = localRepository
	return nil
}

func logJsendWriteError(err error) {
	if err != nil {
		logrus.WithError(err).Errorln("write jsend on response writer")
	}
}

func (api *Api) Run() <-chan error {
	apiErrorChan := make(chan error)
	go func() {
		defer close(apiErrorChan)
		err := api.httpServer.ListenAndServe()
		if err != nil {
			apiErrorChan <- err
		}
	}()
	return apiErrorChan
}

func (api *Api) Shutdown() {
	timeout, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	err := api.httpServer.Shutdown(timeout)
	if err != nil {
		logrus.WithError(err).Errorln("shutdown http server")
	}
}
