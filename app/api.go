package app

import (
	"context"
	"edupaim/xpto-support/app/controllers/command"
	"edupaim/xpto-support/app/services"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Api struct {
	httpServer *http.Server
}

func InitializeApi(c *Config) (*Api, error) {
	r := gin.Default()
	legacyRepository, err := services.InitializeApiLegacyRepository(c.LegacyXptoConfig.Address)
	if err != nil {
		return nil, err
	}
	localRepository, err := services.InitializeArangoLocalStorage(c.ArangoConfig)
	if err != nil {
		return nil, err
	}
	legacyController := command.NewLegacyIntegrateController(legacyRepository, localRepository)
	r.GET("/legacy/integrate", func(c *gin.Context) {
		err := legacyController.LegacyIntegrate(nil)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusOK)
	})
	api := &Api{}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", c.ServerConfig.Port),
		Handler: r,
	}
	api.httpServer = srv
	return api, nil
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
