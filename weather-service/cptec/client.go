package cptec

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/fgouvea/weather/weather-service/weather"
	"golang.org/x/net/html/charset"
)

const (
	getCitiesURL  = "/XML/listaCidades?city=%s"
	getWeatherURL = "/XML/cidade/%s/previsao.xml"
)

var ErrReadingResponse = errors.New("error reading api response")
var ErrFetchingResponse = errors.New("error fetching cptec response")
var ErrCityNotFound = errors.New("city not found")
var ErrMultipleCities = errors.New("multiple cities found with name")

type Client struct {
	Client *http.Client

	getCitiesUrl  string
	getWeatherUrl string
}

func NewClient(httpClient *http.Client, basePath string) *Client {
	return &Client{
		Client: httpClient,

		getCitiesUrl:  fmt.Sprintf("%s%s", basePath, getCitiesURL),
		getWeatherUrl: fmt.Sprintf("%s%s", basePath, getWeatherURL),
	}
}

func (c *Client) FindCity(name string) (weather.City, error) {
	url := fmt.Sprintf(c.getCitiesUrl, url.QueryEscape(name))

	var parsedResponse CitiesResponseTO

	err := getFromAPI(c, url, &parsedResponse)

	if err != nil {
		return weather.City{}, err
	}

	if len(parsedResponse.Cities) == 0 {
		return weather.City{}, ErrCityNotFound
	}

	if len(parsedResponse.Cities) > 1 {
		return weather.City{}, ErrMultipleCities
	}

	return buildCity(parsedResponse.Cities[0]), nil
}

func (c *Client) GetForecast(id string) (weather.CityForecast, error) {
	url := fmt.Sprintf(c.getWeatherUrl, id)

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

func buildCity(to CityTO) weather.City {
	return weather.City{
		ID:    to.ID,
		Name:  to.Name,
		State: to.State,
	}
}

func buildCityForecast(to CityForecastTO) (weather.CityForecast, error) {
	forecast := make([]weather.Forecast, len(to.Forecast))

	for i, forecastTO := range to.Forecast {
		parsedForecast, err := buildForecast(forecastTO)

		if err != nil {
			return weather.CityForecast{}, nil
		}

		forecast[i] = parsedForecast
	}

	return weather.CityForecast{
		Name:     to.Name,
		State:    to.State,
		Date:     to.Date,
		Forecast: forecast,
	}, nil
}

func buildForecast(to ForecastTO) (weather.Forecast, error) {
	min, err := strconv.Atoi(to.MinTemperature)
	if err != nil {
		return weather.Forecast{}, fmt.Errorf("failed to read min temp: %s", to.MinTemperature)
	}

	max, err := strconv.Atoi(to.MaxTemperature)
	if err != nil {
		return weather.Forecast{}, fmt.Errorf("failed to read max temp: %s", to.MaxTemperature)
	}

	iuv, err := strconv.ParseFloat(to.IUV, 32)
	if err != nil {
		return weather.Forecast{}, fmt.Errorf("failed to read IUV: %s", to.IUV)
	}

	fullWeatherName, exists := WeatherName[to.Weather]
	if !exists {
		fullWeatherName = UnknownWeather
	}

	return weather.Forecast{
		Date:           to.Date,
		Weather:        fullWeatherName,
		MinTemperature: min,
		MaxTemperature: max,
		IUV:            iuv,
	}, nil
}
