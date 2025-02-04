package main

import (
	"fmt"

	"go.uber.org/zap"
)

func buildLogger() *zap.SugaredLogger {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	logger, err := config.Build()

	if err != nil {
		panic(fmt.Errorf("error creating logger: %w", err))
	}

	return logger.Sugar()
}

func main() {
	logger := buildLogger()
	defer logger.Sync()

	logger.Infow("application started", "host", "localhost")
}
