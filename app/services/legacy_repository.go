package services

import (
	"edupaim/xpto-support/app/domain"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	RequestLegacyApiError = errors.New("request to legacy api")
	ReadResponseError     = errors.New("read legacy api response")
)

type LegacyRepository interface {
	GetAllNegatives() ([]domain.Negative, error)
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

func (repository *ApiLegacyRepository) GetAllNegatives() ([]domain.Negative, error) {
	resp, err := http.Get(repository.negativesUrl.String())
	if err != nil {
		logrus.WithError(err).Errorln(RequestLegacyApiError.Error())
		return nil, RequestLegacyApiError
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithError(err).Errorln(ReadResponseError.Error())
		return nil, ReadResponseError
	}
	var negatives []domain.Negative
	err = json.Unmarshal(body, &negatives)
	if err != nil {
		logrus.WithError(err).Errorln(domain.DecodeNegativeJsonError.Error())
		return nil, domain.DecodeNegativeJsonError
	}
	return negatives, nil
}
