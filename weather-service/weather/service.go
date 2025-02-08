package weather

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/fgouvea/weather/weather-service/user"
)

type Service struct {
	UserFinder        UserFinder
	CityFinder        CityFinder
	WeatherForecaster WeatherForecaster
	WaveForecaster    WaveForecaster
	Notifier          Notifier
}

func NewService(
	userFinder UserFinder,
	cityFinder CityFinder,
	weatherForecaster WeatherForecaster,
	waveForecaster WaveForecaster,
	notifier Notifier,
) *Service {
	return &Service{
		UserFinder:        userFinder,
		CityFinder:        cityFinder,
		WeatherForecaster: weatherForecaster,
		WaveForecaster:    waveForecaster,
		Notifier:          notifier,
	}
}

func (s *Service) NotifyUser(userID, cityName string) error {
	userEntry, err := s.UserFinder.FindUser(userID)

	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return err
		}

		return fmt.Errorf("unexpected error fetching user: %w", err)
	}

	city, err := s.CityFinder.FindCity(cityName)

	if err != nil {
		if errors.Is(err, ErrCityNotFound) || errors.Is(err, ErrMultipleCities) {
			return err
		}

		return fmt.Errorf("unexpected error fetching city: %w", err)
	}

	weatherForecast, err := s.WeatherForecaster.GetForecast(city.ID)

	if err != nil {
		return fmt.Errorf("unexpected error fetching weather forecast: %w", err)
	}

	waveForecast, err := s.WaveForecaster.GetWaveForecast(city.ID)

	if err != nil && !errors.Is(err, ErrCityNotFound) {
		return fmt.Errorf("unexpected error fetching wave forecast: %w", err)
	}

	return s.sendNotification(userEntry, city, weatherForecast, waveForecast)
}

func (s *Service) sendNotification(
	userEntry user.User,
	city City,
	weatherForecast CityForecast,
	waveForecast CityWaveForecast,
) error {
	var buffer strings.Builder

	buffer.WriteString(fmt.Sprintf("%s, aqui está a previsão do tempo para %s\n\n", userEntry.Name, city.Name))

	for i, forecast := range weatherForecast.Forecast {
		date := formatDate(forecast.Date)

		buffer.WriteString(fmt.Sprintf("%s: %d - %d", date, forecast.MinTemperature, forecast.MaxTemperature))
		if i < (len(weatherForecast.Forecast) - 1) {
			buffer.WriteString("\n")
		}
	}

	if waveForecast != (CityWaveForecast{}) {
		date := formatDate(waveForecast.Date)
		buffer.WriteString("\n\n")
		buffer.WriteString(fmt.Sprintf("Ondas para o dia %s:\n", date))
		buffer.WriteString(fmt.Sprintf("Manhã: %s %.2fm\n", waveForecast.Morning.Swell, waveForecast.Morning.Height))
		buffer.WriteString(fmt.Sprintf("Tarde: %s %.2fm\n", waveForecast.Afternoon.Swell, waveForecast.Afternoon.Height))
		buffer.WriteString(fmt.Sprintf("Noite: %s %.2fm", waveForecast.Evening.Swell, waveForecast.Evening.Height))
	}

	content := buffer.String()

	s.Notifier.Notify(userEntry.ID, content)

	return nil
}

func formatDate(date string) string {
	d, _ := time.Parse("2006-01-02", date)
	return d.Format("02/01/2006")
}
