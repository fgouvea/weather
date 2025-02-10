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
	ErrUnknownChannel  = errors.New("unknown channel")
)

type Notification struct {
	UserID  string
	Content string
	Channel string
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

	sender, exists := s.Senders[notification.Channel]

	if !exists {
		return fmt.Errorf("%w: %s", ErrUnknownChannel, notification.Channel)
	}

	err = sender.Send(recipient, notification.Content)

	if errors.Is(ErrUserOptOut, err) {
		s.Logger.Info("notification skipped", zap.String("sender", notification.Channel), zap.String("userID", recipient.ID))
		return nil
	}

	if err != nil {
		s.Logger.Error("failed to send", zap.String("sender", notification.Channel), zap.String("userID", recipient.ID), zap.Error(err))
		return fmt.Errorf("%w: %w", ErrFailedToProcess, err)
	}

	s.Logger.Info("notification sent", zap.String("sender", notification.Channel), zap.String("userID", recipient.ID))

	return nil
}
