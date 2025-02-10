package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/fgouvea/weather/user-service/api"
	"github.com/fgouvea/weather/user-service/db"
	"github.com/fgouvea/weather/user-service/user"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type AppConfig struct {
	Port       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBDatabase string
}

func readConfigFromEnv() AppConfig {
	return AppConfig{
		Port:       fmt.Sprintf(":%s", readFromEnv("PORT", "8080")),
		DBHost:     readFromEnv("DB_HOST", "localhost"),
		DBPort:     readFromEnv("DB_PORT", "5432"),
		DBUser:     readFromEnv("DB_USER", "admin"),
		DBPassword: readFromEnv("DB_PASSWORD", "admin"),
		DBDatabase: readFromEnv("DB_DATABASE", "weather"),
	}
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

func main() {
	logger := buildLogger()
	defer logger.Sync()

	config := readConfigFromEnv()

	repository, err := db.NewUserRepository(config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBDatabase)

	if err != nil {
		panic(fmt.Sprintf("failed to initialize user repository: %s", err.Error()))
	}

	service := user.NewService(repository, repository)

	handler := api.UserHandler{
		Service: service,
		Logger:  *logger,
	}

	r := chi.NewRouter()

	r.Route("/user-service", func(r chi.Router) {
		r.Get("/health", api.Health)

		r.Route("/user", func(r chi.Router) {
			r.Post("/", handler.CreateUser)
			r.Get("/{userID}", handler.FindUser)
			r.Post("/{userID}/optout", handler.OutOutOfNotifications)
		})
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
