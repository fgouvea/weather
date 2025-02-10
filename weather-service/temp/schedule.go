package temp

import (
	"slices"
	"time"

	"github.com/fgouvea/weather/weather-service/schedule"
)

type InMemoryScheduleRepository struct {
	Schedules []schedule.Schedule
}

func (r *InMemoryScheduleRepository) Save(shcd schedule.Schedule) error {
	r.Schedules = append(r.Schedules, shcd)

	slices.SortFunc[[]schedule.Schedule](r.Schedules, func(a, b schedule.Schedule) int {
		return a.Time.Compare(b.Time)
	})

	return nil
}

func (r *InMemoryScheduleRepository) Find(id string) (schedule.Schedule, error) {
	for _, schdl := range r.Schedules {
		if schdl.ID == id {
			return schdl, nil
		}
	}

	return schedule.Schedule{}, schedule.ErrScheduleNotFound
}

func (r *InMemoryScheduleRepository) Delete(id string) error {
	newScheduleList := make([]schedule.Schedule, len(r.Schedules)-1)

	found := false

	for i, schd := range r.Schedules {
		if schd.ID == id {
			found = true
			continue
		}

		idx := i
		if found {
			idx = i - 1
		}

		newScheduleList[idx] = schd
	}

	if !found {
		return schedule.ErrScheduleNotFound
	}

	r.Schedules = newScheduleList

	return nil
}

func (r *InMemoryScheduleRepository) FindAllBefore(t time.Time) ([]schedule.Schedule, error) {
	index := r.countSchedulesBefore(t)

	if index == 0 {
		return []schedule.Schedule{}, nil
	}

	result := make([]schedule.Schedule, index)

	for i := 0; i < index; i++ {
		result[i] = r.Schedules[i]
	}

	return result, nil
}

func (r *InMemoryScheduleRepository) countSchedulesBefore(t time.Time) int {
	for i, s := range r.Schedules {
		if s.Time.After(t) {
			return i
		}
	}

	return len(r.Schedules)
}
