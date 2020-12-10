package app

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Api struct {
	httpServer *http.Server
}

func initializeApi() (*Api, error) {
	r := gin.Default()
	r.GET("/legacy/integrate", func(c *gin.Context) {
	})
	api := &Api{}
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	api.httpServer = srv
	return nil, nil
}

func (api *Api) run() <-chan error {
	apiErrorChan := make(chan error)
	go func() {
		err := api.httpServer.ListenAndServe()
		if err != nil {
			apiErrorChan <- err
		}
	}()
	return apiErrorChan
}

func (api *Api) shutdown() {
	timeout, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	err := api.httpServer.Shutdown(timeout)
	if err != nil {
		logrus.WithError(err).Errorln("shutdown http server")
	}
}
