package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/fgouvea/weather/weather-service/api"
	"github.com/fgouvea/weather/weather-service/cptec"
	"github.com/fgouvea/weather/weather-service/temp"
	"github.com/fgouvea/weather/weather-service/user"
	"github.com/fgouvea/weather/weather-service/weather"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type AppConfig struct {
	Port             string
	UserServiceHost  string
	CPTECServiceHost string
}

func readConfigFromEnv() AppConfig {
	return AppConfig{
		Port:             fmt.Sprintf(":%s", readFromEnv("PORT", "8080")),
		UserServiceHost:  readFromEnv("USER_SERVICE_HOST", "http://localhost:8080"),
		CPTECServiceHost: readFromEnv("CPTEC_HOST", "http://servicos.cptec.inpe.br"),
	}
}

func main() {
	logger := buildLogger()
	defer logger.Sync()

	config := readConfigFromEnv()

	cptecClient := cptec.NewClient(buildHttpClient(), config.CPTECServiceHost)
	userClient := user.NewClient(buildHttpClient(), config.UserServiceHost)

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

	logger.Info("application started", zap.Any("config", config))
	defer logger.Info("application shutdown")

	http.ListenAndServe(config.Port, r)
}

func readFromEnv(env, def string) string {
	if value := os.Getenv(env); value != "" {
		return value
	}

	return def
}

func buildLogger() *zap.Logger {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	logger, err := config.Build()

	if err != nil {
		panic(fmt.Errorf("error creating logger: %w", err))
	}

	return logger
}

func buildHttpClient() *http.Client {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	return &http.Client{Transport: tr}
}
