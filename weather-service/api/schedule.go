package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/fgouvea/weather/weather-service/schedule"
	"go.uber.org/zap"
)

type WeatherScheduler interface {
	Schedule(userID, cityName string, scheduleTime time.Time) error
}

type ScheduleHandler struct {
	Scheduler WeatherScheduler
	Logger    *zap.Logger
}

func (h *ScheduleHandler) Schedule(w http.ResponseWriter, r *http.Request) {
	var body ScheduleRequest
	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		h.Logger.Error("error reading request body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	scheduleTime, err := time.Parse(time.RFC3339, body.Time)

	if err != nil {
		h.Logger.Error("error reading schedule time", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.Scheduler.Schedule(body.UserID, body.City, scheduleTime)

	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(schedule.ErrScheduleInThePast, err) {
			status = http.StatusBadRequest
		}

		h.Logger.Error("error scheduling weather info", zap.String("userID", body.UserID), zap.String("city", body.City), zap.Error(err))
		w.WriteHeader(status)
		return
	}

	h.Logger.Info("weather info scheduled", zap.String("userID", body.UserID), zap.String("city", body.City))
	w.WriteHeader(http.StatusOK)
}
