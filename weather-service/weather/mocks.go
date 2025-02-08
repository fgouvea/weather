package weather

import (
	"errors"

	"github.com/fgouvea/weather/weather-service/user"
)

var runtimeError = errors.New("runtime error")

var testUser = user.User{
	ID:   "user-id",
	Name: "Fulano",
}

var testCity = City{
	ID:   "city-id",
	Name: "Test City",
}

var testWeatherForecast = CityForecast{
	Forecast: []Forecast{
		Forecast{Date: "2025-02-07", MinTemperature: 1, MaxTemperature: 31},
		Forecast{Date: "2025-02-08", MinTemperature: 2, MaxTemperature: 32},
		Forecast{Date: "2025-02-09", MinTemperature: 3, MaxTemperature: 33},
		Forecast{Date: "2025-02-10", MinTemperature: 4, MaxTemperature: 34},
	},
}

var testWavesForecast = CityWaveForecast{
	Date:      "2025-02-07",
	Morning:   WaveForecast{Swell: "Fraca", Height: 0.1},
	Afternoon: WaveForecast{Swell: "Moderada", Height: 0.23},
	Evening:   WaveForecast{Swell: "Forte", Height: 0.456},
}

type mockClient struct {
	findUserCalls  []string
	findUserResult user.User
	findUserError  error

	findCityCalls  []string
	findCityResult City
	findCityError  error

	getForecastCalls  []string
	getForecastResult CityForecast
	getForecastError  error

	getWaveForecastCalls  []string
	getWaveForecastResult CityWaveForecast
	getWaveForecastError  error

	notifyCallsUserID  []string
	notifyCallsContent []string
	notifyError        error
}

var _ UserFinder = (*mockClient)(nil)
var _ CityFinder = (*mockClient)(nil)
var _ WeatherForecaster = (*mockClient)(nil)
var _ WaveForecaster = (*mockClient)(nil)
var _ Notifier = (*mockClient)(nil)

func (m *mockClient) FindUser(id string) (user.User, error) {
	m.findUserCalls = append(m.findUserCalls, id)
	return m.findUserResult, m.findUserError
}

func (m *mockClient) FindCity(name string) (City, error) {
	m.findCityCalls = append(m.findCityCalls, name)
	return m.findCityResult, m.findCityError
}

func (m *mockClient) GetForecast(id string) (CityForecast, error) {
	m.getForecastCalls = append(m.getForecastCalls, id)
	return m.getForecastResult, m.getForecastError
}

func (m *mockClient) GetWaveForecast(id string) (CityWaveForecast, error) {
	m.getWaveForecastCalls = append(m.getWaveForecastCalls, id)
	return m.getWaveForecastResult, m.getWaveForecastError
}

func (m *mockClient) Notify(userID string, content string) error {
	m.notifyCallsUserID = append(m.notifyCallsUserID, userID)
	m.notifyCallsContent = append(m.notifyCallsContent, content)
	return m.notifyError
}
