package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/fgouvea/weather/user-service/user"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type UserProcessor interface {
	Find(id string) (*user.User, error)
	Create(name, webNotificationId string) (*user.User, error)
	OptOutOfNotifications(id string) error
}

type UserHandler struct {
	Service UserProcessor
	Logger  zap.Logger
}

func (h *UserHandler) FindUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	result, err := h.Service.Find(userID)

	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			h.Logger.Info("user not found", zap.String("userID", userID))
			w.WriteHeader(http.StatusNotFound)
			return
		}

		h.Logger.Error("error finding user", zap.String("userID", userID), zap.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseBody, err := json.Marshal(buildUserTO(result))

	if err != nil {
		h.Logger.Error("error writing find response", zap.String("userID", userID), zap.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.Logger.Info("user found", zap.String("userID", userID), zap.String("userName", result.Name))
	w.Write([]byte(responseBody))
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var body CreateUserRequestTO
	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		h.Logger.Error("error reading create body", zap.String("error", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := h.Service.Create(body.Name, body.WebNotificationID)

	if err != nil {
		h.Logger.Error("error creating user", zap.String("userName", body.Name), zap.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseBody, err := json.Marshal(buildUserTO(result))

	if err != nil {
		h.Logger.Error("error writing create response", zap.String("userID", result.Id), zap.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.Logger.Info("user created", zap.String("userID", result.Id), zap.String("userName", result.Name))
	w.Write([]byte(responseBody))
}
