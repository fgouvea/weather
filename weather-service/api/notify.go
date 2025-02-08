package api

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type WeatherNotifier interface {
	NotifyUser(userID string, city string) error
}

type WeatherHandler struct {
	Notifier WeatherNotifier
	Logger   *zap.Logger
}

func (h *WeatherHandler) NotifyUser(w http.ResponseWriter, r *http.Request) {
	var body NotifyUserRequest
	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		h.Logger.Error("error reading request body", zap.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.Notifier.NotifyUser(body.UserID, body.City)

	if err != nil {
		h.Logger.Error("error notifying user", zap.String("userID", body.UserID), zap.String("city", body.City), zap.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.Logger.Info("weather info sent to user", zap.String("userID", body.UserID), zap.String("city", body.City))
	w.WriteHeader(http.StatusOK)
}
