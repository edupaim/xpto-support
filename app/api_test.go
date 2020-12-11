package app

import (
	"github.com/onsi/gomega"
	"net/http"
	"testing"
)

func TestApi_Run(t *testing.T) {
	g := gomega.NewWithT(t)
	t.Run("success integrate legacy database", func(t *testing.T) {
		api, err := InitializeApi(&Config{ServerConfig: ServerConfig{Port: 5051}})
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		defer api.Shutdown()
		errChan := api.Run()
		resp, err := http.Get("http://localhost:5051/legacy/integrate")
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		g.Expect(resp.StatusCode).Should(gomega.Equal(http.StatusOK))
		g.Expect(errChan).ShouldNot(gomega.Receive())
	})
}
