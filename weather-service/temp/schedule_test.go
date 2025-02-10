package temp

import (
	"testing"
	"time"

	"github.com/fgouvea/weather/weather-service/schedule"
	"github.com/stretchr/testify/assert"
)

func TestTempRepository_GetAllBefore(t *testing.T) {
	time1, _ := time.Parse(time.RFC3339, "2025-02-09T10:00:00Z")
	time2, _ := time.Parse(time.RFC3339, "2025-02-09T11:00:00Z")
	time3, _ := time.Parse(time.RFC3339, "2025-02-09T12:00:00Z")

	schedule1 := schedule.Schedule{ID: "1", UserID: "USER-1", CityName: "rio de janeiro", Time: time1}
	schedule2 := schedule.Schedule{ID: "2", UserID: "USER-2", CityName: "volta redonda", Time: time2}
	schedule3 := schedule.Schedule{ID: "3", UserID: "USER-3", CityName: "arraial do cabo", Time: time3}

	repository := InMemoryScheduleRepository{}

	repository.Save(schedule2)
	repository.Save(schedule3)
	repository.Save(schedule1)

	assert.Equal(t, []schedule.Schedule{schedule1, schedule2, schedule3}, repository.Schedules)

	timeSearch, _ := time.Parse(time.RFC3339, "2025-02-09T11:50:00Z")

	result, _ := repository.FindAllBefore(timeSearch)

	assert.Equal(t, []schedule.Schedule{schedule1, schedule2}, result)

	repository.Delete("1")

	assert.Equal(t, []schedule.Schedule{schedule2, schedule3}, repository.Schedules)
}
