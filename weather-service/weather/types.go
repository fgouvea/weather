package weather

import "github.com/fgouvea/weather/weather-service/user"

type UserFinder interface {
	FindUser(id string) (user.User, error)
}

type CityFinder interface {
	FindCity(name string) (City, error)
}

type WeatherForecaster interface {
	GetForecast(id string) (CityForecast, error)
}

type WaveForecaster interface {
	GetWaveForecast(id string) (CityWaveForecast, error)
}

type Notifier interface {
	Notify(userID, content string) error
}
