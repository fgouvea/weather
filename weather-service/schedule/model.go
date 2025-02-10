package schedule

import (
	"errors"
	"time"
)

var (
	ErrScheduleNotFound  = errors.New("schedule not found")
	ErrScheduleInThePast = errors.New("schedule time cannot be in the past")
	ErrFailedToSave      = errors.New("failed to save schedule")
	ErrFailedToProcess   = errors.New("failed to process schedule")
)

const (
	StatusActive     = "active"
	StatusProcessing = "processing"
	StatusCompleted  = "completed"
)

type Schedule struct {
	ID       string
	UserID   string
	CityName string
	Status   string
	Time     time.Time
}
