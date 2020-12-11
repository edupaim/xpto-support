package services

import (
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

const negativesCollectionName = "negatives"

type ArangoConfig struct {
	Address  string
	Password string
	User     string
	Database string
}

func GetDatabase(dbName string, arangoClient driver.Client) (driver.Database, error) {
	exist, err := arangoClient.DatabaseExists(nil, dbName)
	if err != nil {
		return nil, err
	}
	if !exist {
		db, err := arangoClient.CreateDatabase(nil, dbName, nil)
		if err != nil {
			return nil, err
		}
		return db, nil
	}
	db, err := arangoClient.Database(nil, dbName)
	return db, err
}

func GetNegativesCollection(db driver.Database) (driver.Collection, error) {
	exist, err := db.CollectionExists(nil, negativesCollectionName)
	if err != nil {
		return nil, err
	}
	if !exist {
		coll, err := db.CreateCollection(nil, negativesCollectionName, nil)
		if err != nil {
			return nil, err
		}
		return coll, nil
	}
	return db.Collection(nil, negativesCollectionName)
}

func GetArangoClient(c ArangoConfig) (driver.Client, error) {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{c.Address},
	})
	if err != nil {
		return nil, err
	}
	arangoClient, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(c.User, c.Password),
	})
	if err != nil {
		return nil, err
	}
	return arangoClient, nil
}
