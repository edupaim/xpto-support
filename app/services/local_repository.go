package services

import (
	"context"
	"edupaim/xpto-support/app/domain"
	"errors"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/sirupsen/logrus"
	"go.elastic.co/apm"
	"strings"
)

var (
	SaveNegativeError  = errors.New("save negative")
	QueryNegativeError = errors.New("query negative from database")
)

type LocalRepository interface {
	SaveNegative(negative domain.Negative, ctx context.Context) error
	GetNegativeByQuery(document map[string][]string, ctx context.Context) ([]domain.Negative, error)
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
	_, err := localStorage.arangoDatabase.Query(nil, `
	UPSERT {customerDocument: @customerDocument, companyDocument: @companyDocument}
	INSERT @negative
	UPDATE @negative IN @@coll`, map[string]interface{}{
		"customerDocument": negative.CustomerDocument,
		"companyDocument":  negative.CompanyDocument,
		"@coll":            NegativesCollectionName,
		"negative":         negative,
	})
	if err != nil {
		logCtx.WithError(err).Errorln(SaveNegativeError.Error())
		return SaveNegativeError
	}
	return nil
}

func (localStorage *ArangoLocalRepository) GetNegativeByQuery(queryNegative map[string][]string, ctx context.Context) ([]domain.Negative, error) {
	span, ctx := apm.StartSpan(ctx, "GetNegativeByQuery()", "runtime.localRepository")
	defer span.End()
	logCtx := logrus.WithContext(ctx)
	query := NewForQuery(NegativesCollectionName, "negative")
	for key, value := range queryNegative {
		query.Filter(key, "==", value[0])
	}
	query.Return()
	cursor, err := localStorage.arangoDatabase.Query(nil,
		query.String(),
		map[string]interface{}{})
	if err != nil {
		logCtx.WithError(err).Errorln(QueryNegativeError.Error())
		return nil, QueryNegativeError
	}
	return localStorage.iterateCursorAndReadNegatives(cursor, logCtx)
}

func (localStorage *ArangoLocalRepository) iterateCursorAndReadNegatives(cursor driver.Cursor, logCtx *logrus.Entry) ([]domain.Negative, error) {
	var negatives []domain.Negative
	for cursor.HasMore() {
		negative := domain.Negative{}
		_, err := cursor.ReadDocument(nil, &negative)
		if err != nil {
			logCtx.WithError(err).Errorln(domain.DecodeNegativeJsonError.Error())
			return nil, domain.DecodeNegativeJsonError
		}
		negatives = append(negatives, negative)
	}
	return negatives, nil
}

type ArangoQueryBuilder struct {
	query       strings.Builder
	counterName string
}

func NewForQuery(repositoryName string, counterName string) *ArangoQueryBuilder {
	aqb := &ArangoQueryBuilder{
		counterName: counterName,
	}
	aqb.query.WriteString(fmt.Sprintf("FOR %s IN %s", aqb.counterName, repositoryName))
	return aqb
}

func (aqb *ArangoQueryBuilder) Filter(fieldName string, condition string, value interface{}) *ArangoQueryBuilder {
	aqb.query.WriteString(" ")

	switch value.(type) {
	case string:
		aqb.query.WriteString(fmt.Sprintf("FILTER %s.%s %s %q", aqb.counterName, fieldName, condition, value))
	default:
		aqb.query.WriteString(fmt.Sprintf("FILTER %s.%s %s %v", aqb.counterName, fieldName, condition, value))
	}

	return aqb
}

func (aqb *ArangoQueryBuilder) Return() *ArangoQueryBuilder {
	aqb.query.WriteString(" ")
	aqb.query.WriteString(fmt.Sprintf("RETURN %s", aqb.counterName))
	return aqb
}

func (aqb *ArangoQueryBuilder) String() string {
	return aqb.query.String()
}
