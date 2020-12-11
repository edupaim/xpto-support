package app

import (
	"edupaim/xpto-support/app/domain"
	"edupaim/xpto-support/app/services"
	"github.com/arangodb/go-driver"
	arangohttp "github.com/arangodb/go-driver/http"
	"github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
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
	t.Run("success integrate legacy database", func(t *testing.T) {
		api, err := InitializeApi(&config)
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		defer api.Shutdown()
		errChan := api.Run()
		resp, err := http.Get("http://localhost:5051/legacy/integrate")
		g.Expect(err).ShouldNot(gomega.HaveOccurred())
		g.Expect(resp.StatusCode).Should(gomega.Equal(http.StatusOK))
		g.Expect(errChan).ShouldNot(gomega.Receive())
		negatives := queryAllNegativesFromDatabase(g, db)
		g.Expect(negatives).Should(gomega.ContainElements(domain.Negative{
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

func queryAllNegativesFromDatabase(g *gomega.WithT, db driver.Database) []domain.Negative {
	cursor, err := db.Query(nil, "FOR n IN negatives RETURN n", map[string]interface{}{})
	g.Expect(err).ShouldNot(gomega.HaveOccurred())
	var negatives []domain.Negative
	for cursor.HasMore() {
		negative := domain.Negative{}
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
