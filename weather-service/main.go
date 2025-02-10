package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/fgouvea/weather/weather-service/api"
	"github.com/fgouvea/weather/weather-service/cptec"
	"github.com/fgouvea/weather/weather-service/db"
	"github.com/fgouvea/weather/weather-service/notification"
	"github.com/fgouvea/weather/weather-service/queue"
	"github.com/fgouvea/weather/weather-service/schedule"
	"github.com/fgouvea/weather/weather-service/user"
	"github.com/fgouvea/weather/weather-service/weather"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type AppConfig struct {
	Port              string
	UserServiceHost   string
	CPTECServiceHost  string
	RabbitHost        string
	NotificationQueue string
	ScheduleQueue     string
	ScheduleConsumers int
	JobInterval       time.Duration
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBDatabase        string
}

func readConfigFromEnv() AppConfig {
	scheduleConsumers, err := strconv.Atoi(readFromEnv("SCHEDULE_CONSUMERS", "1"))

	if err != nil {
		panic("number of consumers must be integer")
	}

	jobInterval, err := time.ParseDuration(readFromEnv("JOB_INTERVAL", "5m"))

	if err != nil {
		panic("job interval must be duration")
	}

	return AppConfig{
		Port:              fmt.Sprintf(":%s", readFromEnv("PORT", "8081")),
		UserServiceHost:   readFromEnv("USER_SERVICE_HOST", "http://localhost:8080"),
		CPTECServiceHost:  readFromEnv("CPTEC_HOST", "http://servicos.cptec.inpe.br"),
		RabbitHost:        readFromEnv("RABBIT_HOST", "amqp://guest:guest@localhost:5672/"),
		NotificationQueue: readFromEnv("NOTIFICATION_QUEUE", "notifications"),
		ScheduleQueue:     readFromEnv("SCHEDULE_QUEUE", "schedules"),
		ScheduleConsumers: scheduleConsumers,
		JobInterval:       jobInterval,
		DBHost:            readFromEnv("DB_HOST", "localhost"),
		DBPort:            readFromEnv("DB_PORT", "5432"),
		DBUser:            readFromEnv("DB_USER", "admin"),
		DBPassword:        readFromEnv("DB_PASSWORD", "admin"),
		DBDatabase:        readFromEnv("DB_DATABASE", "weather"),
	}
}

func main() {
	logger := buildLogger()
	defer logger.Sync()

	config := readConfigFromEnv()

	// Message broker

	notificationPublisher, err := notification.NewPublisher(config.RabbitHost, config.NotificationQueue)

	if err != nil {
		panic(fmt.Sprintf("failed do initialize notification publisher: %s", err.Error()))
	}

	defer notificationPublisher.Close()

	schedulePublisher, err := queue.NewPublisher(config.RabbitHost, config.ScheduleQueue)

	if err != nil {
		panic(fmt.Sprintf("failed do initialize schedule publisher: %s", err.Error()))
	}

	defer schedulePublisher.Close()

	// Clients

	cptecClient := cptec.NewClient(buildHttpClient(), config.CPTECServiceHost)
	userClient := user.NewClient(buildHttpClient(), config.UserServiceHost)

	// Repositories

	scheduleRepository, err := db.NewScheduleRepository(config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBDatabase)

	// Services

	weatherService := weather.NewService(userClient, cptecClient, cptecClient, cptecClient, notificationPublisher)

	scheduleService := schedule.NewService(scheduleRepository, weatherService, weatherService)

	// Consumers

	scheduleConsumer, err := schedule.NewConsumer(config.RabbitHost, config.ScheduleQueue, config.ScheduleConsumers, scheduleService, logger)

	if err != nil {
		panic(fmt.Sprintf("failed do initialize schedule consumer: %s", err.Error()))
	}

	// Jobs

	scheduleJob := schedule.NewJob(config.JobInterval, schedule.NewPublisher(scheduleRepository, schedulePublisher), scheduleRepository, logger)

	// Handlers

	weatherHandler := &api.WeatherHandler{
		Notifier: weatherService,
		Logger:   logger,
	}

	scheduleHandler := &api.ScheduleHandler{
		Scheduler: scheduleService,
		Logger:    logger,
	}

	r := chi.NewRouter()

	r.Route("/weather-service", func(r chi.Router) {
		r.Get("/health", api.Health)
		r.Post("/notify", weatherHandler.NotifyUser)
		r.Post("/schedule", scheduleHandler.Schedule)
	})

	logger.Info("application started", zap.Any("config", config))
	defer logger.Info("application shutdown")

	scheduleConsumer.Start()
	scheduleJob.Start()

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
