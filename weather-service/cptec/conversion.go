package cptec

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fgouvea/weather/weather-service/weather"
)

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
		UpdatedAt: to.Date,
		Forecast:  forecast,
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

func buildCityWaveForecast(to CityWaveForecastTO) (weather.CityWaveForecast, error) {
	morning, err := buildWaveForecast(to.Morning)
	if err != nil {
		return weather.CityWaveForecast{}, fmt.Errorf("failed to read morning wave forecast: %w", err)
	}

	afternoon, err := buildWaveForecast(to.Afternoon)
	if err != nil {
		return weather.CityWaveForecast{}, fmt.Errorf("failed to read afternoon wave forecast: %w", err)
	}

	evening, err := buildWaveForecast(to.Evening)
	if err != nil {
		return weather.CityWaveForecast{}, fmt.Errorf("failed to read evening wave forecast: %w", err)
	}

	date, err := fixDate(strings.Split(to.Morning.Date, " ")[0])
	if err != nil {
		return weather.CityWaveForecast{}, fmt.Errorf("failed to parse date: %w", err)
	}

	updatedAt, err := fixDate(to.UpdatedAt)
	if err != nil {
		return weather.CityWaveForecast{}, fmt.Errorf("failed to parse update date: %w", err)
	}

	return weather.CityWaveForecast{
		UpdatedAt: updatedAt,
		Date:      date,
		Morning:   morning,
		Afternoon: afternoon,
		Evening:   evening,
	}, nil
}

func buildWaveForecast(to WaveForecastTO) (weather.WaveForecast, error) {
	height, err := strconv.ParseFloat(to.Height, 64)
	if err != nil {
		return weather.WaveForecast{}, fmt.Errorf("failed to read wave height: %s", to.Height)
	}

	wind, err := strconv.ParseFloat(to.Wind, 64)
	if err != nil {
		return weather.WaveForecast{}, fmt.Errorf("failed to read wind: %s", to.Wind)
	}

	return weather.WaveForecast{
		Swell:         to.Swell,
		Height:        height,
		Wind:          wind,
		WaveDirection: to.WaveDirection,
		WindDirection: to.WindDirection,
	}, nil
}

func fixDate(brDate string) (string, error) {
	date, err := time.Parse("02-01-2006", brDate)

	if err != nil {
		return "", err
	}

	return date.Format("2006-01-02"), nil
}
