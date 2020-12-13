package app

import (
	"clevergo.tech/jsend"
	"github.com/gin-gonic/gin"
	"go.elastic.co/apm/module/apmgin"
	"net/http"
)

func initializeApiRouter(controllers ApiControllers) http.Handler {
	r := gin.Default()
	r.Use(apmgin.Middleware(r))
	r.POST("/legacy/integrate", func(c *gin.Context) {
		err := controllers.legacyIntegrate.LegacyIntegrate(nil, c.Request.Context())
		if err != nil {
			err = jsend.Error(c.Writer, "integrate application with legacy api", http.StatusInternalServerError)
			logJsendWriteError(err)
			return
		}
		c.Status(http.StatusOK)
	})
	r.GET("/negatives", func(c *gin.Context) {
		customerQuery := c.Request.URL.Query()
		negative, err := controllers.negativesQuery.GetByQuery(customerQuery, c.Request.Context())
		if err != nil {
			err = jsend.Error(c.Writer, "get negative by customer document", http.StatusInternalServerError)
			logJsendWriteError(err)
			return
		}
		err = jsend.Success(c.Writer, negative, http.StatusOK)
		logJsendWriteError(err)
	})
	return r
}
