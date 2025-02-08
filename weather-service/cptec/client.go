package cptec

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/fgouvea/weather/weather-service/weather"
	"golang.org/x/net/html/charset"
)

const (
	getCitiesURL  = "/XML/listaCidades?city=%s"
	getWeatherURL = "/XML/cidade/%s/previsao.xml"
	getWaveURL    = "/XML/cidade/%s/dia/%d/ondas.xml"
)

var ErrReadingResponse = errors.New("error reading api response")
var ErrFetchingResponse = errors.New("error fetching cptec response")
var ErrCityNotFound = errors.New("city not found")
var ErrMultipleCities = errors.New("multiple cities found with name")

type Client struct {
	Client *http.Client

	getCitiesURL  string
	getWeatherURL string
	getWaveURL    string
}

func NewClient(httpClient *http.Client, basePath string) *Client {
	return &Client{
		Client: httpClient,

		getCitiesURL:  fmt.Sprintf("%s%s", basePath, getCitiesURL),
		getWeatherURL: fmt.Sprintf("%s%s", basePath, getWeatherURL),
		getWaveURL:    fmt.Sprintf("%s%s", basePath, getWaveURL),
	}
}

func (c *Client) FindCity(name string) (weather.City, error) {
	url := fmt.Sprintf(c.getCitiesURL, url.QueryEscape(name))

	var parsedResponse CitiesResponseTO

	err := getFromAPI(c, url, &parsedResponse)

	if err != nil {
		return weather.City{}, err
	}

	parsedCity, err := chooseCity(parsedResponse, name)

	if err != nil {
		return weather.City{}, err
	}

	return buildCity(parsedCity), nil
}

func (c *Client) GetForecast(id string) (weather.CityForecast, error) {
	url := fmt.Sprintf(c.getWeatherURL, id)

	var parsedResponse CityForecastTO

	err := getFromAPI(c, url, &parsedResponse)

	if err != nil {
		return weather.CityForecast{}, err
	}

	if parsedResponse.Name == "null" {
		return weather.CityForecast{}, ErrCityNotFound
	}

	forecast, err := buildCityForecast(parsedResponse)

	if err != nil {
		return weather.CityForecast{}, err
	}

	return forecast, nil
}

func (c *Client) GetWaveForecast(id string) (weather.CityWaveForecast, error) {
	url := fmt.Sprintf(c.getWaveURL, id, 0)

	var parsedResponse CityWaveForecastTO

	err := getFromAPI(c, url, &parsedResponse)

	if err != nil {
		return weather.CityWaveForecast{}, err
	}

	if parsedResponse.Name == "undefined" {
		return weather.CityWaveForecast{}, ErrCityNotFound
	}

	forecast, err := buildCityWaveForecast(parsedResponse)

	if err != nil {
		return weather.CityWaveForecast{}, err
	}

	return forecast, nil
}

func getFromAPI[T any](c *Client, url string, parsedResponse *T) error {
	response, err := c.Client.Get(url)

	if err != nil {
		return fmt.Errorf("%w: %w", ErrFetchingResponse, err)
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: unexpected status code: %d", ErrFetchingResponse, response.StatusCode)
	}

	decoder := xml.NewDecoder(response.Body)
	decoder.CharsetReader = charset.NewReaderLabel

	err = decoder.Decode(&parsedResponse)

	if err != nil {
		return fmt.Errorf("%w: %w", ErrReadingResponse, err)
	}

	return nil
}

func chooseCity(response CitiesResponseTO, name string) (CityTO, error) {
	if len(response.Cities) == 0 {
		return CityTO{}, ErrCityNotFound
	}

	if len(response.Cities) == 1 {
		return response.Cities[0], nil
	}

	for _, city := range response.Cities {
		if strings.ToLower(city.Name) == strings.ToLower(name) {
			return city, nil
		}
	}

	return CityTO{}, ErrMultipleCities
}
