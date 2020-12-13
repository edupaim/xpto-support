package services

import (
	"edupaim/xpto-support/app/domain"
	"errors"
	"github.com/arangodb/go-driver"
	"github.com/sirupsen/logrus"
)

var (
	SaveNegativeError  = errors.New("save negative")
	QueryNegativeError = errors.New("query negative from database")
)

type LocalRepository interface {
	SaveNegative(negative domain.Negative) error
	GetNegativeByCustomerDocument(document string) ([]domain.Negative, error)
}

type ArangoLocalRepository struct {
	arangoDatabase      driver.Database
	negativesCollection driver.Collection
}

func InitializeArangoLocalRepository(c ArangoConfig) (*ArangoLocalRepository, error) {
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
	return NewArangoLocalRepository(db, negativesColl), nil
}

func NewArangoLocalRepository(db driver.Database, negativesColl driver.Collection) *ArangoLocalRepository {
	return &ArangoLocalRepository{
		arangoDatabase:      db,
		negativesCollection: negativesColl,
	}
}

func (localStorage *ArangoLocalRepository) SaveNegative(negative domain.Negative) error {
	meta, err := localStorage.negativesCollection.CreateDocument(nil, negative)
	if err != nil {
		logrus.WithError(err).Errorln(SaveNegativeError.Error())
		return SaveNegativeError
	}
	logrus.WithField("key", meta.Key).
		Debugln("create document on negatives collection")
	return nil
}

func (localStorage *ArangoLocalRepository) GetNegativeByCustomerDocument(document string) ([]domain.Negative, error) {
	cursor, err := localStorage.arangoDatabase.Query(nil, "FOR n IN @@coll FILTER n.customerDocument == @customerDocument RETURN n",
		map[string]interface{}{
			"@coll":            negativesCollectionName,
			"customerDocument": document,
		})
	if err != nil {
		logrus.WithError(err).Errorln(QueryNegativeError.Error())
		return nil, QueryNegativeError
	}
	var negatives []domain.Negative
	for cursor.HasMore() {
		negative := domain.Negative{}
		_, err = cursor.ReadDocument(nil, &negative)
		if err != nil {
			logrus.WithError(err).Errorln(domain.DecodeNegativeJsonError.Error())
			return nil, domain.DecodeNegativeJsonError
		}
		negatives = append(negatives, negative)
	}
	return negatives, nil
}
