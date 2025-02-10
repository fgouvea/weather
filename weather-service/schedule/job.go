package schedule

import (
	"time"

	"go.uber.org/zap"
)

type ScheduleSearcher interface {
	FindAllBefore(t time.Time) ([]Schedule, error)
}

type Job struct {
	Interval  time.Duration
	Publisher *Publisher
	Searcher  ScheduleSearcher
	Logger    *zap.Logger
}

func NewJob(interval time.Duration, publisher *Publisher, searcher ScheduleSearcher, logger *zap.Logger) *Job {
	return &Job{
		Interval:  interval,
		Publisher: publisher,
		Searcher:  searcher,
		Logger:    logger,
	}
}

func (j *Job) Start() {
	ticker := time.NewTicker(j.Interval)

	go func() {
		j.Logger.Info("starting schedule job")

		for currentTime := range ticker.C {
			// TODO: Lock job to avoid multiple instances running the same job together

			threshold := time.Now().Add(j.Interval)

			schedules, err := j.Searcher.FindAllBefore(threshold)

			if err != nil {
				j.Logger.Error("error querying schedules within threshold", zap.Error(err))
				continue
			}

			j.Logger.Info("running schedule job", zap.Int("schedulesFound", len(schedules)))

			for _, schedule := range schedules {
				timer := time.NewTimer(schedule.Time.Sub(currentTime))

				go func() {
					_ = <-timer.C
					j.Publisher.Publish(schedule)
				}()
			}
		}

		j.Logger.Info("stopping schedule job")
	}()
}
