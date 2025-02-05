package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/fgouvea/weather/user-service/user"
	"github.com/go-chi/chi"
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
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseBody, err := json.Marshal(buildUserTO(result))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(responseBody))
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var body CreateUserRequestTO
	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := h.Service.Create(body.Name, body.WebNotificationID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseBody, err := json.Marshal(buildUserTO(result))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(responseBody))
}
