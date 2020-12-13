package services

import (
	"context"
	"edupaim/xpto-support/app/domain"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"go.elastic.co/apm"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	RequestLegacyApiError = errors.New("request to legacy api")
	ReadResponseError     = errors.New("read legacy api response")
)

type LegacyRepository interface {
	GetAllNegatives(context.Context) ([]domain.Negative, error)
}

type ApiLegacyRepository struct {
	negativesUrl *url.URL
}

func InitializeApiLegacyRepository(rawApiUrl string) (*ApiLegacyRepository, error) {
	apiUrl, err := url.Parse(rawApiUrl)
	if err != nil {
		return nil, err
	}
	return NewApiLegacyRepository(apiUrl), nil
}

func NewApiLegacyRepository(apiUrl *url.URL) *ApiLegacyRepository {
	negativesUrl := &url.URL{}
	*negativesUrl = *apiUrl
	negativesUrl.Path = "/negatives"
	return &ApiLegacyRepository{negativesUrl: negativesUrl}
}

func (repository *ApiLegacyRepository) GetAllNegatives(ctx context.Context) ([]domain.Negative, error) {
	span, ctx := apm.StartSpan(ctx, "GetAllNegatives()", "runtime.legacyRepository")
	defer span.End()
	logCtx := logrus.WithContext(ctx)
	resp, err := http.Get(repository.negativesUrl.String())
	if err != nil {
		logCtx.WithError(err).Errorln(RequestLegacyApiError.Error())
		return nil, RequestLegacyApiError
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logCtx.WithError(err).Errorln(ReadResponseError.Error())
		return nil, ReadResponseError
	}
	var negatives []domain.Negative
	err = json.Unmarshal(body, &negatives)
	if err != nil {
		logCtx.WithError(err).Errorln(domain.DecodeNegativeJsonError.Error())
		return nil, domain.DecodeNegativeJsonError
	}
	return negatives, nil
}
