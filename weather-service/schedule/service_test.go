package schedule

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/fgouvea/weather/weather-service/user"
	"github.com/stretchr/testify/assert"
)

func TestService_Schedule(t *testing.T) {
	tests := []struct {
		name                  string
		scheduleTime          string
		validateError         error
		saveError             error
		expectedError         error
		expectedValidateCalls int
		expectedSaveCalls     int
	}{
		{
			name:                  "success",
			scheduleTime:          "2028-02-01T10:00:00Z",
			validateError:         nil,
			saveError:             nil,
			expectedError:         nil,
			expectedValidateCalls: 1,
			expectedSaveCalls:     1,
		},
		{
			name:                  "schedule in the past",
			scheduleTime:          "2005-02-01T10:00:00Z",
			validateError:         nil,
			saveError:             nil,
			expectedError:         ErrScheduleInThePast,
			expectedValidateCalls: 0,
			expectedSaveCalls:     0,
		},
		{
			name:                  "error validating",
			scheduleTime:          "2028-02-01T10:00:00Z",
			validateError:         user.ErrUserNotFound,
			saveError:             nil,
			expectedError:         user.ErrUserNotFound,
			expectedValidateCalls: 1,
			expectedSaveCalls:     0,
		},
		{
			name:                  "error saving schedule",
			scheduleTime:          "2028-02-01T10:00:00Z",
			validateError:         nil,
			saveError:             errors.New("failed to connect to db"),
			expectedError:         ErrFailedToSave,
			expectedValidateCalls: 1,
			expectedSaveCalls:     1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &serviceMock{
				validateError: tt.validateError,
				saveError:     tt.saveError,
			}

			service := NewService(mock, mock, mock)

			scheduleTime, _ := time.Parse(time.RFC3339, tt.scheduleTime)

			err := service.Schedule("USER-ID", "city name", scheduleTime)

			assert.True(t, errors.Is(err, tt.expectedError), fmt.Sprintf("Expected: %s / Actual: %s", tt.expectedError, err))

			assert.Equal(t, tt.expectedValidateCalls, len(mock.validateCalls))
			for i := 0; i < tt.expectedValidateCalls; i++ {
				assert.Equal(t, userAndCity{userID: "USER-ID", cityName: "city name"}, mock.validateCalls[i])
			}

			assert.Equal(t, tt.expectedSaveCalls, len(mock.saveCalls))
			for i := 0; i < tt.expectedSaveCalls; i++ {
				assert.Equal(t, "USER-ID", mock.saveCalls[i].UserID)
				assert.Equal(t, StatusActive, mock.saveCalls[i].Status)
				assert.Equal(t, "city name", mock.saveCalls[i].CityName)
				assert.Equal(t, scheduleTime, mock.saveCalls[i].Time)
			}
		})
	}
}

func TestService_Process(t *testing.T) {
	tests := []struct {
		name                string
		notifyError         error
		saveError           error
		expectedError       error
		expectedNotifyCalls int
		expectedSaveCalls   int
	}{
		{
			name:                "success",
			notifyError:         nil,
			saveError:           nil,
			expectedError:       nil,
			expectedNotifyCalls: 1,
			expectedSaveCalls:   1,
		},
		{
			name:                "error notifying",
			notifyError:         errors.New("runtime error"),
			saveError:           nil,
			expectedError:       ErrFailedToProcess,
			expectedNotifyCalls: 1,
			expectedSaveCalls:   0,
		},
		{
			name:                "error saving schedule",
			notifyError:         nil,
			saveError:           errors.New("failed to connect to db"),
			expectedError:       ErrFailedToSave,
			expectedNotifyCalls: 1,
			expectedSaveCalls:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &serviceMock{
				notifyError: tt.notifyError,
				saveError:   tt.saveError,
			}

			service := NewService(mock, mock, mock)

			schedule := Schedule{
				UserID:   "USER-ID",
				CityName: "city name",
				Status:   StatusActive,
			}

			err := service.Process(schedule)

			assert.True(t, errors.Is(err, tt.expectedError), fmt.Sprintf("Expected: %s / Actual: %s", tt.expectedError, err))

			assert.Equal(t, tt.expectedNotifyCalls, len(mock.notifyCalls))
			for i := 0; i < tt.expectedNotifyCalls; i++ {
				assert.Equal(t, userAndCity{userID: "USER-ID", cityName: "city name"}, mock.notifyCalls[i])
			}

			assert.Equal(t, tt.expectedSaveCalls, len(mock.saveCalls))
			for i := 0; i < tt.expectedSaveCalls; i++ {
				assert.Equal(t, "USER-ID", mock.saveCalls[i].UserID)
				assert.Equal(t, StatusCompleted, mock.saveCalls[i].Status)
				assert.Equal(t, "city name", mock.saveCalls[i].CityName)
			}
		})
	}
}
