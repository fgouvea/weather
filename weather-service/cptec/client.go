package cptec

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/net/html/charset"
)

const (
	getCitiesURL = "/XML/listaCidades"
)

var ErrReadingResponse = errors.New("error reading api response")
var ErrFetchingCities = errors.New("error fetching cities")
var ErrCityNotFound = errors.New("city not found")
var ErrMultipleCities = errors.New("multiple cities found with name")

type Client struct {
	Client *http.Client

	getCitiesUrl string
}

func NewClient(httpClient *http.Client, basePath string) *Client {
	return &Client{
		Client: httpClient,

		getCitiesUrl: fmt.Sprintf("%s%s", basePath, getCitiesURL),
	}
}

func (c *Client) FindCity(name string) (City, error) {
	url := fmt.Sprintf("%s?city=%s", c.getCitiesUrl, url.QueryEscape(name))

	fmt.Println(url)
	response, err := c.Client.Get(url)

	if err != nil {
		return City{}, fmt.Errorf("%w: %w", ErrFetchingCities, err)
	}

	if response.StatusCode != http.StatusOK {
		return City{}, fmt.Errorf("%w: unexpected status code: %d", ErrFetchingCities, response.StatusCode)
	}

	var parsedResponse CitiesResponse

	decoder := xml.NewDecoder(response.Body)
	decoder.CharsetReader = charset.NewReaderLabel

	err = decoder.Decode(&parsedResponse)

	if err != nil {
		return City{}, fmt.Errorf("%w: %w", ErrReadingResponse, err)
	}

	if len(parsedResponse.Cities) == 0 {
		return City{}, ErrCityNotFound
	}

	if len(parsedResponse.Cities) > 1 {
		return City{}, ErrMultipleCities
	}

	return parsedResponse.Cities[0], nil
}
