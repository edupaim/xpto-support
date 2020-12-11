package services

import (
	"edupaim/xpto-support/app/domain"
	"errors"
	"github.com/arangodb/go-driver"
	"github.com/sirupsen/logrus"
)

var (
	SaveNegativeError = errors.New("save negative")
)

type LocalRepository interface {
	SaveNegative(negative domain.Negative) error
}

type ArangoLocalStorage struct {
	arangoDatabase      driver.Database
	negativesCollection driver.Collection
}

func InitializeArangoLocalStorage(c ArangoConfig) (*ArangoLocalStorage, error) {
	arangoClient, err := GetArangoClient(c)
	if err != nil {
		return nil, err
	}
	db, err := GetDatabase(c.Database, arangoClient)
	if err != nil {
		return nil, err
	}
	negativesColl, err := GetNegativesCollection(db)
	if err != nil {
		return nil, err
	}
	return NewArangoLocalStorage(db, negativesColl), nil
}

func NewArangoLocalStorage(db driver.Database, negativesColl driver.Collection) *ArangoLocalStorage {
	return &ArangoLocalStorage{
		arangoDatabase:      db,
		negativesCollection: negativesColl,
	}
}

func (localStorage *ArangoLocalStorage) SaveNegative(negative domain.Negative) error {
	meta, err := localStorage.negativesCollection.CreateDocument(nil, negative)
	if err != nil {
		return SaveNegativeError
	}
	logrus.WithField("key", meta.Key).
		Debugln("create document on negatives collection")
	return nil
}
