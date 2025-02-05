package main

import (
	"fmt"

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

func main() {
	logger := buildLogger()
	defer logger.Sync()

	logger.Info("application started", zap.String("host", "localhost"))

}
