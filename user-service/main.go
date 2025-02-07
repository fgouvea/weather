package main

import (
	"fmt"
	"net/http"

	"github.com/fgouvea/weather/user-service/api"
	"github.com/fgouvea/weather/user-service/temp"
	"github.com/fgouvea/weather/user-service/user"
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

	repository := temp.NewInMemoryRTepository()

	service := user.NewService(repository, repository)

	handler := api.UserHandler{
		Service: service,
		Logger:  *logger,
	}

	r := chi.NewRouter()

	r.Route("/user-service", func(r chi.Router) {
		r.Get("/health", api.Health)

		r.Route("/user", func(r chi.Router) {
			r.Get("/{userID}", handler.FindUser)
			r.Post("/", handler.CreateUser)
		})
	})

	logger.Info("application started", zap.String("port", port))
	defer logger.Info("application shutdown")

	http.ListenAndServe(port, r)
}
