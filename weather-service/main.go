package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fgouvea/weather/weather-service/api"
	"github.com/fgouvea/weather/weather-service/cptec"
	"github.com/fgouvea/weather/weather-service/temp"
	"github.com/fgouvea/weather/weather-service/user"
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

const port = ":8081"

func main() {
	logger := buildLogger()
	defer logger.Sync()

	cptecClient := cptec.NewClient(buildHttpClient(), "http://servicos.cptec.inpe.br")

	userClient := user.NewClient(buildHttpClient(), "http://localhost:8080")

	weatherService := weather.NewService(userClient, cptecClient, cptecClient, cptecClient, &temp.TempNotifier{})

	weatherHandler := &api.WeatherHandler{
		Notifier: weatherService,
		Logger:   logger,
	}

	r := chi.NewRouter()

	r.Route("/weather-service", func(r chi.Router) {
		r.Get("/health", api.Health)
		r.Post("/notify", weatherHandler.NotifyUser)
	})

	logger.Info("application started", zap.String("port", port))
	defer logger.Info("application shutdown")

	http.ListenAndServe(port, r)
}

func buildHttpClient() *http.Client {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	return &http.Client{Transport: tr}
}
