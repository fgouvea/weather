package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fgouvea/weather/weather-service/api"
	"github.com/fgouvea/weather/weather-service/cptec"
	"github.com/fgouvea/weather/weather-service/temp"
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

	weatherService := weather.NewService(&temp.TempUserClient{}, cptecClient, cptecClient, cptecClient, &temp.TempNotifier{})

	r := chi.NewRouter()

	r.Route("/weather-service", func(r chi.Router) {
		r.Get("/health", api.Health)
	})

	logger.Info("application started", zap.String("port", port))
	defer logger.Info("application shutdown")

	// http.ListenAndServe(port, r)

	err := weatherService.NotifyUser("USER-123456", "rio de janeiro")

	if err != nil {
		fmt.Println(err.Error())
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
