package user

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Saver interface {
	Save(user *User) error
}

type Finder interface {
	Find(id string) (*User, error)
}

type Service struct {
	Saver  Saver
	Finder Finder
	Logger zap.Logger
}

func NewService(saver Saver, finder Finder) *Service {
	return &Service{
		Saver:  saver,
		Finder: finder,
	}
}

func (s *Service) Create(name, webNotificationId string) (*User, error) {
	user := &User{
		ID:   "USER-" + uuid.New().String(),
		Name: name,
		NotificationConfig: NotificationConfig{
			Enabled: true,
			Web: WebNotificationConfig{
				Enabled: len(webNotificationId) > 0,
				Id:      webNotificationId,
			},
		},
	}

	err := s.save(user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Find(id string) (*User, error) {
	user, err := s.Finder.Find(id)

	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, err
		}

		return nil, fmt.Errorf("unexpected error fetching user: %w", err)
	}

	return user, nil
}

func (s *Service) OptOutOfNotifications(id string) error {
	user, err := s.Find(id)

	if err != nil {
		return err
	}

	user.NotificationConfig.Enabled = false

	return s.save(user)
}

func (s *Service) save(user *User) error {
	err := s.Saver.Save(user)

	if err != nil {
		return fmt.Errorf("unexpected error saving user: %w", err)
	}

	return nil
}
