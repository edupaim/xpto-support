package services

import (
	"context"
	"edupaim/xpto-support/app/domain"
	"errors"
	"github.com/arangodb/go-driver"
	"github.com/sirupsen/logrus"
	"go.elastic.co/apm"
)

var (
	SaveNegativeError  = errors.New("save negative")
	QueryNegativeError = errors.New("query negative from database")
)

type LocalRepository interface {
	SaveNegative(negative domain.Negative, ctx context.Context) error
	GetNegativeByCustomerDocument(document string, ctx context.Context) ([]domain.Negative, error)
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

func (localStorage *ArangoLocalRepository) SaveNegative(negative domain.Negative, ctx context.Context) error {
	span, ctx := apm.StartSpan(ctx, "SaveNegative()", "runtime.localRepository")
	defer span.End()
	logCtx := logrus.WithContext(ctx)
	meta, err := localStorage.negativesCollection.CreateDocument(nil, negative)
	if err != nil {
		logCtx.WithError(err).Errorln(SaveNegativeError.Error())
		return SaveNegativeError
	}
	logCtx.WithField("key", meta.Key).
		Debugln("create document on negatives collection")
	return nil
}

func (localStorage *ArangoLocalRepository) GetNegativeByCustomerDocument(document string, ctx context.Context) ([]domain.Negative, error) {
	span, ctx := apm.StartSpan(ctx, "GetNegativeByCustomerDocument()", "runtime.localRepository")
	defer span.End()
	logCtx := logrus.WithContext(ctx)
	cursor, err := localStorage.arangoDatabase.Query(nil, "FOR n IN @@coll FILTER n.customerDocument == @customerDocument RETURN n",
		map[string]interface{}{
			"@coll":            NegativesCollectionName,
			"customerDocument": document,
		})
	if err != nil {
		logCtx.WithError(err).Errorln(QueryNegativeError.Error())
		return nil, QueryNegativeError
	}
	var negatives []domain.Negative
	for cursor.HasMore() {
		negative := domain.Negative{}
		_, err = cursor.ReadDocument(nil, &negative)
		if err != nil {
			logCtx.WithError(err).Errorln(domain.DecodeNegativeJsonError.Error())
			return nil, domain.DecodeNegativeJsonError
		}
		negatives = append(negatives, negative)
	}
	return negatives, nil
}
