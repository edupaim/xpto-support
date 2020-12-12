package app

import (
	"clevergo.tech/jsend"
	"edupaim/xpto-support/app/domain"
	"edupaim/xpto-support/app/services"
	"encoding/json"
	"github.com/arangodb/go-driver"
	arangohttp "github.com/arangodb/go-driver/http"
	"github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

var config = Config{
	ServerConfig: ServerConfig{
		Port: 5051,
	},
	ArangoConfig: services.ArangoConfig{
		Address:  "http://localhost:8529",
		Password: "root",
		User:     "root",
		Database: "xpto-support",
	},
	LegacyXptoConfig: LegacyXptoConfig{
		Address: "http://localhost:8080",
	},
}

func TestApi_Run(t *testing.T) {
	g := gomega.NewWithT(t)
	logrus.SetLevel(logrus.DebugLevel)
	c := getArangoClient(g)
	db := getDatabaseFromArango(g, c)
	_ = db.Remove(nil)

	api, err := InitializeApi(&config)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	defer api.Shutdown()
	errChan := api.Run()
	g.Consistently(errChan).ShouldNot(gomega.Receive())

	t.Run("success integrate legacy database", func(t *testing.T) {
		resp, err := http.Post("http://localhost:5051/legacy/integrate", "application/json", nil)
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		g.Expect(resp.StatusCode).Should(gomega.Equal(http.StatusOK))

		resp, err = http.Get("http://localhost:5051/negatives?customerDocument=51537476467")
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		g.Expect(resp.StatusCode).Should(gomega.Equal(http.StatusOK))
		response := readResponseBody(g, resp)
		expectedResponse := getJsendJsonFromNegative(g, []domain.Negative{
			{
				CompanyDocument:  "59291534000167",
				CompanyName:      "ABC S.A.",
				CustomerDocument: "51537476467",
				Value:            1235.23,
				Contract:         "bc063153-fb9e-4334-9a6c-0d069a42065b",
				DebtDate:         time.Date(2015, 11, 13, 23, 32, 51, 0, time.UTC),
				InclusionDate:    time.Date(2020, 11, 13, 23, 32, 51, 0, time.UTC),
			},
			{
				CompanyDocument:  "77723018000146",
				CompanyName:      "123 S.A.",
				CustomerDocument: "51537476467",
				Value:            400.0,
				Contract:         "5f206825-3cfe-412f-8302-cc1b24a179b0",
				DebtDate:         time.Date(2015, 10, 12, 23, 32, 51, 0, time.UTC),
				InclusionDate:    time.Date(2020, 10, 12, 23, 32, 51, 0, time.UTC),
			},
		})
		g.Expect(response).Should(gomega.MatchJSON(expectedResponse))
		g.Expect(errChan).ShouldNot(gomega.Receive())
	})
}

func readResponseBody(g *gomega.WithT, resp *http.Response) []byte {
	response, err := ioutil.ReadAll(resp.Body)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	return response
}

func getJsendJsonFromNegative(g *gomega.WithT, data interface{}) []byte {
	negativeExpected := jsend.New(data)
	negativeJsonExpected, err := json.Marshal(negativeExpected)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	return negativeJsonExpected
}

func getDatabaseFromArango(g *gomega.WithT, c driver.Client) driver.Database {
	db, err := c.Database(nil, config.ArangoConfig.Database)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	return db
}

func getArangoClient(g *gomega.WithT) driver.Client {
	conn, err := arangohttp.NewConnection(arangohttp.ConnectionConfig{
		Endpoints: []string{config.ArangoConfig.Address},
	})
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	c, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(config.ArangoConfig.User, config.ArangoConfig.Password),
	})
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	return c
}
