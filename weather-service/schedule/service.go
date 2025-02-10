package schedule

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Validator interface {
	Validate(userID, cityName string) error
}

type ScheduleSaver interface {
	Save(schedule Schedule) error
}

type Notifier interface {
	NotifyUser(userID, cityName string) error
}

type Service struct {
	Validator Validator
	Saver     ScheduleSaver
	Notifier  Notifier
}

func NewService(saver ScheduleSaver, validator Validator, notifier Notifier) *Service {
	return &Service{
		Validator: validator,
		Notifier:  notifier,
		Saver:     saver,
	}
}

func (s *Service) Schedule(userID, cityName string, scheduleTime time.Time) error {
	if scheduleTime.Before(time.Now()) {
		return ErrScheduleInThePast
	}

	err := s.Validator.Validate(userID, cityName)

	if err != nil {
		return err
	}

	schedule := Schedule{
		ID:       fmt.Sprintf("SCHEDULE-%s", uuid.New()),
		Status:   StatusActive,
		UserID:   userID,
		CityName: cityName,
		Time:     scheduleTime,
	}

	err = s.Saver.Save(schedule)

	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToSave, err)
	}

	return nil
}

func (s *Service) Process(schedule Schedule) error {
	err := s.Notifier.NotifyUser(schedule.UserID, schedule.CityName)

	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToProcess, err)
	}

	schedule.Status = StatusCompleted

	err = s.Saver.Save(schedule)

	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToSave, err)
	}

	return nil
}
