package app

import (
	"edupaim/xpto-support/app/models"
	"github.com/arangodb/go-driver"
	arangohttp "github.com/arangodb/go-driver/http"
	"github.com/onsi/gomega"
	"net/http"
	"testing"
	"time"
)

var config = Config{
	ServerConfig: ServerConfig{
		Port: 5051,
	},
	ArangoConfig: ArangoConfig{
		Endpoint: "http://localhost:8529",
	},
}

func TestApi_Run(t *testing.T) {
	g := gomega.NewWithT(t)
	c := getArangoClient(g)
	t.Run("success integrate legacy database", func(t *testing.T) {
		api, err := InitializeApi(&config)
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		defer api.Shutdown()
		errChan := api.Run()
		resp, err := http.Get("http://localhost:5051/legacy/integrate")
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		g.Expect(resp.StatusCode).Should(gomega.Equal(http.StatusOK))
		g.Expect(errChan).ShouldNot(gomega.Receive())
		db := getDatabaseFromArango(g, c)
		negatives := queryAllNegativatesFromDatabase(g, db)
		g.Expect(negatives).Should(gomega.ContainElements(models.Negative{
			CompanyDocument:  "59291534000167",
			CompanyName:      "ABC S.A.",
			CustomerDocument: "51537476467",
			Value:            1235.23,
			Contract:         "bc063153-fb9e-4334-9a6c-0d069a42065b",
			DebtDate:         time.Date(2015, 11, 13, 23, 32, 51, 0, time.UTC),
			InclusionDate:    time.Date(2020, 11, 13, 23, 32, 51, 0, time.UTC),
		}))
	})
}

func queryAllNegativatesFromDatabase(g *gomega.WithT, db driver.Database) []models.Negative {
	cursor, err := db.Query(nil, "FOR n IN negative RETURN n", map[string]interface{}{})
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	var negatives []models.Negative
	for cursor.HasMore() {
		negative := models.Negative{}
		_, err := cursor.ReadDocument(nil, &negative)
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		negatives = append(negatives, negative)
	}
	return negatives
}

func getDatabaseFromArango(g *gomega.WithT, c driver.Client) driver.Database {
	db, err := c.Database(nil, config.ArangoConfig.Database)
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	return db
}

func getArangoClient(g *gomega.WithT) driver.Client {
	conn, err := arangohttp.NewConnection(arangohttp.ConnectionConfig{
		Endpoints: []string{config.ArangoConfig.Endpoint},
	})
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	c, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(config.ArangoConfig.User, config.ArangoConfig.Password),
	})
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	return c
}
