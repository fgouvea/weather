package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fgouvea/weather/weather-service/api"
	"github.com/fgouvea/weather/weather-service/cptec"
	"github.com/fgouvea/weather/weather-service/weather"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func buildLogger() *zap.Logger {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	logger, err := config.Build()

	if err != nil {
		panic(fmt.Errorf("error creating logger: %w", err))
	}

	return logger
}

const port = ":8080"

func main() {
	logger := buildLogger()
	defer logger.Sync()

	cptecClient := cptec.NewClient(buildHttpClient(), "http://servicos.cptec.inpe.br")

	r := chi.NewRouter()

	r.Route("/weather-service", func(r chi.Router) {
		r.Get("/health", api.Health)
	})

	logger.Info("application started", zap.String("port", port))
	defer logger.Info("application shutdown")

	// http.ListenAndServe(port, r)

	city, err := cptecClient.FindCity("santos")

	if err != nil {
		fmt.Printf("%s\n", err.Error())
	} else {
		fmt.Println(city.ID)
		fmt.Println(city.Name)
		fmt.Println(city.State)
	}

	fmt.Println()

	forecast, err := cptecClient.GetForecast(city.ID)

	if err != nil {
		fmt.Printf("%s\n", err.Error())
	} else {
		for _, f := range forecast.Forecast {
			fmt.Println(f.Date)
			fmt.Println(f.Weather)
			fmt.Printf("%d - %d\n", f.MinTemperature, f.MaxTemperature)
			fmt.Println()
		}
	}

	waveForecast, err := cptecClient.GetWaveForecast(city.ID)

	if err != nil {
		if err == cptec.ErrCityNotFound {
			fmt.Println("Cidade não litorânea")
		} else {
			fmt.Printf("%s\n", err.Error())
		}
	} else {
		printWaveForecast(waveForecast.Morning, waveForecast.Date, "Manhã")
		printWaveForecast(waveForecast.Afternoon, waveForecast.Date, "Tarde")
		printWaveForecast(waveForecast.Evening, waveForecast.Date, "Noite")
	}
}

func printWaveForecast(f weather.WaveForecast, date, period string) {
	fmt.Printf("%s: %s\n", period, date)
	fmt.Printf("Agitação: %s\n", f.Swell)
	fmt.Printf("Ondas: %fm %s\n", f.Height, f.WaveDirection)
	fmt.Printf("Vento: %f %s\n", f.Wind, f.WindDirection)
	fmt.Println()
}

func buildHttpClient() *http.Client {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	return &http.Client{Transport: tr}
}
