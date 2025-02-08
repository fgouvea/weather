package notification

import (
	"errors"
	"fmt"

	"github.com/fgouvea/weather/notification-service/user"
	"go.uber.org/zap"
)

var (
	ErrFailedToProcess = errors.New("failed to process notification")
	ErrUserOptOut      = errors.New("user opted out of receiving notifications")
)

type Notification struct {
	UserID  string
	Content string
}

type UserFinder interface {
	FindUser(id string) (user.User, error)
}

type Sender interface {
	Send(recipient user.User, content string) error
}

type Service struct {
	UserFinder UserFinder
	Senders    map[string]Sender
	Logger     *zap.Logger
}

func NewService(userFinder UserFinder, senders map[string]Sender, logger *zap.Logger) *Service {
	return &Service{
		UserFinder: userFinder,
		Senders:    senders,
		Logger:     logger,
	}
}

func (s *Service) Process(notification Notification) error {

	recipient, err := s.UserFinder.FindUser(notification.UserID)

	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToProcess, err)
	}

	if !recipient.NotificationConfig.Enabled {
		return nil
	}

	for senderName, sender := range s.Senders {
		err = sender.Send(recipient, notification.Content)

		if !errors.Is(ErrUserOptOut, err) {
			continue
		}

		if err != nil {
			s.Logger.Error("failed to send", zap.String("sender", senderName), zap.String("userID", recipient.ID), zap.Error(err))
			continue
		}

		s.Logger.Info("notification sent", zap.String("sender", senderName), zap.String("userID", recipient.ID))
	}

	return nil
}
