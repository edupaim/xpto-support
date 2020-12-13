package app

import (
	"clevergo.tech/jsend"
	"edupaim/xpto-support/app/domain"
	"edupaim/xpto-support/app/services"
	"encoding/json"
	"fmt"
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
	coll, err := db.Collection(nil, services.NegativesCollectionName)
	if err == nil {
		_ = coll.Truncate(nil)
	}

	api, err := InitializeApi(&config)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	defer api.Shutdown()
	errChan := api.Run()
	g.Consistently(errChan).ShouldNot(gomega.Receive())
	applicationEndpoint := fmt.Sprintf("http://localhost:%d", config.ServerConfig.Port)

	t.Run("integrate and require negatives by customer document", func(t *testing.T) {
		resp, err := http.Post(applicationEndpoint+"/legacy/integrate", "application/json", nil)
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		g.Expect(resp.StatusCode).Should(gomega.Equal(http.StatusOK))

		customerDocument1 := "51537476467"
		negativesRoute := fmt.Sprintf("/negatives?customerDocument=%s", customerDocument1)
		resp, err = http.Get(applicationEndpoint + negativesRoute)
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		g.Expect(resp.StatusCode).Should(gomega.Equal(http.StatusOK))
		response := readResponseBody(g, resp)
		expectedResponse := getJsendJson(g, []domain.Negative{
			{
				CompanyDocument:  "59291534000167",
				CompanyName:      "ABC S.A.",
				CustomerDocument: domain.CryptDocument(customerDocument1),
				Value:            1235.23,
				Contract:         "bc063153-fb9e-4334-9a6c-0d069a42065b",
				DebtDate:         time.Date(2015, 11, 13, 23, 32, 51, 0, time.UTC),
				InclusionDate:    time.Date(2020, 11, 13, 23, 32, 51, 0, time.UTC),
			},
			{
				CompanyDocument:  "77723018000146",
				CompanyName:      "123 S.A.",
				CustomerDocument: domain.CryptDocument(customerDocument1),
				Value:            400.0,
				Contract:         "5f206825-3cfe-412f-8302-cc1b24a179b0",
				DebtDate:         time.Date(2015, 10, 12, 23, 32, 51, 0, time.UTC),
				InclusionDate:    time.Date(2020, 10, 12, 23, 32, 51, 0, time.UTC),
			},
		})
		g.Expect(response).Should(gomega.MatchJSON(expectedResponse))

		customerDocument2 := "26658236674"
		negativesRoute = fmt.Sprintf("/negatives?customerDocument=%s", customerDocument2)
		resp, err = http.Get(applicationEndpoint + negativesRoute)
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		g.Expect(resp.StatusCode).Should(gomega.Equal(http.StatusOK))
		response = readResponseBody(g, resp)
		expectedResponse = getJsendJson(g, []domain.Negative{
			{
				CompanyDocument:  "04843574000182",
				CompanyName:      "DBZ S.A.",
				CustomerDocument: domain.CryptDocument(customerDocument2),
				Value:            59.99,
				Contract:         "3132f136-3889-4efb-bf92-e1efbb3fe15e",
				DebtDate:         time.Date(2015, 9, 11, 23, 32, 51, 0, time.UTC),
				InclusionDate:    time.Date(2020, 9, 11, 23, 32, 51, 0, time.UTC),
			},
		})
		g.Expect(response).Should(gomega.MatchJSON(expectedResponse))

		customerDocument3 := "62824334010"
		negativesRoute = fmt.Sprintf("/negatives?customerDocument=%s", customerDocument3)
		resp, err = http.Get(applicationEndpoint + negativesRoute)
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		g.Expect(resp.StatusCode).Should(gomega.Equal(http.StatusOK))
		response = readResponseBody(g, resp)
		expectedResponse = getJsendJson(g, []domain.Negative{
			{
				CompanyDocument:  "23993551000107",
				CompanyName:      "XPTO S.A.",
				CustomerDocument: domain.CryptDocument(customerDocument3),
				Value:            230.50,
				Contract:         "8b441dbb-3bb4-4fc9-9b46-bdaad00a7a98",
				DebtDate:         time.Date(2015, 8, 10, 23, 32, 51, 0, time.UTC),
				InclusionDate:    time.Date(2020, 8, 10, 23, 32, 51, 0, time.UTC),
			},
		})
		g.Expect(response).Should(gomega.MatchJSON(expectedResponse))

		customerDocument4 := "25124543043"
		negativesRoute = fmt.Sprintf("/negatives?customerDocument=%s", customerDocument4)
		resp, err = http.Get(applicationEndpoint + negativesRoute)
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		g.Expect(resp.StatusCode).Should(gomega.Equal(http.StatusOK))
		response = readResponseBody(g, resp)
		expectedResponse = getJsendJson(g, []domain.Negative{
			{
				CompanyDocument:  "70170935000100",
				CompanyName:      "ASD S.A.",
				CustomerDocument: domain.CryptDocument(customerDocument4),
				Value:            10340.67,
				Contract:         "d6628a0e-d4dd-4f14-8591-2ddc7f1bbeff",
				DebtDate:         time.Date(2015, 7, 9, 23, 32, 51, 0, time.UTC),
				InclusionDate:    time.Date(2020, 7, 9, 23, 32, 51, 0, time.UTC),
			},
		})
		g.Expect(response).Should(gomega.MatchJSON(expectedResponse))
	})

	t.Run("2 times integrate, should not duplicate persisted negatives", func(t *testing.T) {
		resp, err := http.Post(applicationEndpoint+"/legacy/integrate", "application/json", nil)
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		g.Expect(resp.StatusCode).Should(gomega.Equal(http.StatusOK))

		resp, err = http.Post(applicationEndpoint+"/legacy/integrate", "application/json", nil)
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		g.Expect(resp.StatusCode).Should(gomega.Equal(http.StatusOK))

		customerDocument5 := "25124543043"
		negativesRoute := fmt.Sprintf("/negatives?customerDocument=%s", customerDocument5)
		resp, err = http.Get(applicationEndpoint + negativesRoute)
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		g.Expect(resp.StatusCode).Should(gomega.Equal(http.StatusOK))
		response := readResponseBody(g, resp)
		expectedResponse := getJsendJson(g, []domain.Negative{
			{
				CompanyDocument:  "70170935000100",
				CompanyName:      "ASD S.A.",
				CustomerDocument: domain.CryptDocument(customerDocument5),
				Value:            10340.67,
				Contract:         "d6628a0e-d4dd-4f14-8591-2ddc7f1bbeff",
				DebtDate:         time.Date(2015, 7, 9, 23, 32, 51, 0, time.UTC),
				InclusionDate:    time.Date(2020, 7, 9, 23, 32, 51, 0, time.UTC),
			},
		})
		g.Expect(response).Should(gomega.MatchJSON(expectedResponse))
	})

	t.Run("integrate and requite negative by company document", func(t *testing.T) {
		resp, err := http.Post(applicationEndpoint+"/legacy/integrate", "application/json", nil)
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		g.Expect(resp.StatusCode).Should(gomega.Equal(http.StatusOK))

		companyDocument := "59291534000167"
		negativesRoute := fmt.Sprintf("/negatives?companyDocument=%s", companyDocument)
		resp, err = http.Get(applicationEndpoint + negativesRoute)
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		g.Expect(resp.StatusCode).Should(gomega.Equal(http.StatusOK))
		response := readResponseBody(g, resp)
		expectedResponse := getJsendJson(g, []domain.Negative{
			{
				CompanyDocument:  domain.CryptDocument(companyDocument),
				CompanyName:      "ABC S.A.",
				CustomerDocument: "51537476467",
				Value:            1235.23,
				Contract:         "bc063153-fb9e-4334-9a6c-0d069a42065b",
				DebtDate:         time.Date(2015, 11, 13, 23, 32, 51, 0, time.UTC),
				InclusionDate:    time.Date(2020, 11, 13, 23, 32, 51, 0, time.UTC),
			},
		})
		g.Expect(response).Should(gomega.MatchJSON(expectedResponse))
	})

	t.Run("integrate and requite all negatives", func(t *testing.T) {
		resp, err := http.Post(applicationEndpoint+"/legacy/integrate", "application/json", nil)
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		g.Expect(resp.StatusCode).Should(gomega.Equal(http.StatusOK))

		negativesRoute := fmt.Sprintf("/negatives")
		resp, err = http.Get(applicationEndpoint + negativesRoute)
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		g.Expect(resp.StatusCode).Should(gomega.Equal(http.StatusOK))
		response := readResponseBody(g, resp)
		expectedResponse := getJsendJson(g, []domain.Negative{
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
			{
				CompanyDocument:  "04843574000182",
				CompanyName:      "DBZ S.A.",
				CustomerDocument: "26658236674",
				Value:            59.99,
				Contract:         "3132f136-3889-4efb-bf92-e1efbb3fe15e",
				DebtDate:         time.Date(2015, 9, 11, 23, 32, 51, 0, time.UTC),
				InclusionDate:    time.Date(2020, 9, 11, 23, 32, 51, 0, time.UTC),
			},
			{
				CompanyDocument:  "23993551000107",
				CompanyName:      "XPTO S.A.",
				CustomerDocument: "62824334010",
				Value:            230.50,
				Contract:         "8b441dbb-3bb4-4fc9-9b46-bdaad00a7a98",
				DebtDate:         time.Date(2015, 8, 10, 23, 32, 51, 0, time.UTC),
				InclusionDate:    time.Date(2020, 8, 10, 23, 32, 51, 0, time.UTC),
			},
			{
				CompanyDocument:  "70170935000100",
				CompanyName:      "ASD S.A.",
				CustomerDocument: "25124543043",
				Value:            10340.67,
				Contract:         "d6628a0e-d4dd-4f14-8591-2ddc7f1bbeff",
				DebtDate:         time.Date(2015, 7, 9, 23, 32, 51, 0, time.UTC),
				InclusionDate:    time.Date(2020, 7, 9, 23, 32, 51, 0, time.UTC),
			},
		})
		g.Expect(response).Should(gomega.MatchJSON(expectedResponse))
	})

	g.Expect(errChan).ShouldNot(gomega.Receive())
}

func readResponseBody(g *gomega.WithT, resp *http.Response) []byte {
	response, err := ioutil.ReadAll(resp.Body)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	return response
}

func getJsendJson(g *gomega.WithT, data interface{}) []byte {
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
