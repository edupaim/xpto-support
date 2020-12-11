package services

import (
	"edupaim/xpto-support/app/domain"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	RequestNegativesToLegacyApi = errors.New("request negatives to legacy api")
	ReadResponseFromLegacyApi   = errors.New("read response from legacy api")
	DecodeNegativesFromResponse = errors.New("decode negatives from response")
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
		return nil, RequestNegativesToLegacyApi
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, ReadResponseFromLegacyApi
	}
	var negatives []domain.Negative
	err = json.Unmarshal(body, &negatives)
	if err != nil {
		return nil, DecodeNegativesFromResponse
	}
	return negatives, nil
}
