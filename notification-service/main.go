package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/fgouvea/weather/notification-service/api"
	"github.com/fgouvea/weather/notification-service/user"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type AppConfig struct {
	Port                   string
	UserServiceHost        string
	WebNotificationAPIHost string
	RabbitHost             string
	NotificationQueue      string
}

func readConfigFromEnv() AppConfig {
	return AppConfig{
		Port:                   fmt.Sprintf(":%s", readFromEnv("PORT", "8082")),
		UserServiceHost:        readFromEnv("USER_SERVICE_HOST", "http://localhost:8080"),
		WebNotificationAPIHost: readFromEnv("WEB_NOTIFICATION_API_HOST", "http://localhost:8083"),
		RabbitHost:             readFromEnv("RABBIT_HOST", "amqp://guest:guest@localhost:5672/"),
		NotificationQueue:      readFromEnv("NOTIFICATION_QUEUE", "notifications"),
	}
}

func main() {
	logger := buildLogger()
	defer logger.Sync()

	config := readConfigFromEnv()

	// Message broker

	// Clients

	_ = user.NewClient(buildHttpClient(), config.UserServiceHost)

	// Services

	// API

	r := chi.NewRouter()

	r.Route("/notification-service", func(r chi.Router) {
		r.Get("/health", api.Health)
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
