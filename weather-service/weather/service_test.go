package weather

import (
	"errors"
	"fmt"
	"testing"

	"github.com/fgouvea/weather/weather-service/user"
	"github.com/stretchr/testify/assert"
)

func TestService_NotifyUser(t *testing.T) {
	tests := []struct {
		name                         string
		userResult                   user.User
		userError                    error
		cityResult                   City
		cityError                    error
		weatherResult                CityForecast
		weatherError                 error
		waveResult                   CityWaveForecast
		waveError                    error
		notifierError                error
		expectedError                error
		expectedFindUserCalls        []string
		expectedFindCityCalls        []string
		expectedGetForecastCalls     []string
		expectedGetWaveForecastCalls []string
		expectedNotifications        []string
	}{
		{
			name:                         "success with waves",
			userResult:                   testUser,
			cityResult:                   testCity,
			weatherResult:                testWeatherForecast,
			waveResult:                   testWavesForecast,
			expectedError:                nil,
			expectedFindUserCalls:        []string{"user-id"},
			expectedFindCityCalls:        []string{"test city"},
			expectedGetForecastCalls:     []string{"city-id"},
			expectedGetWaveForecastCalls: []string{"city-id"},
			expectedNotifications:        []string{"Fulano, aqui está a previsão do tempo para Test City\n\n07/02/2025: 1 - 31\n08/02/2025: 2 - 32\n09/02/2025: 3 - 33\n10/02/2025: 4 - 34\n\nOndas para o dia 07/02/2025:\nManhã: Fraca 0.10m\nTarde: Moderada 0.23m\nNoite: Forte 0.46m"},
		},
		{
			name:                         "success without waves",
			userResult:                   testUser,
			cityResult:                   testCity,
			weatherResult:                testWeatherForecast,
			waveError:                    ErrCityNotFound,
			expectedError:                nil,
			expectedFindUserCalls:        []string{"user-id"},
			expectedFindCityCalls:        []string{"test city"},
			expectedGetForecastCalls:     []string{"city-id"},
			expectedGetWaveForecastCalls: []string{"city-id"},
			expectedNotifications:        []string{"Fulano, aqui está a previsão do tempo para Test City\n\n07/02/2025: 1 - 31\n08/02/2025: 2 - 32\n09/02/2025: 3 - 33\n10/02/2025: 4 - 34"},
		},
		{
			name:                         "user not found",
			userError:                    user.ErrUserNotFound,
			expectedError:                user.ErrUserNotFound,
			expectedFindUserCalls:        []string{"user-id"},
			expectedFindCityCalls:        nil,
			expectedGetForecastCalls:     nil,
			expectedGetWaveForecastCalls: nil,
			expectedNotifications:        nil,
		},
		{
			name:                         "unexpected error getting user",
			userError:                    runtimeError,
			expectedError:                runtimeError,
			expectedFindUserCalls:        []string{"user-id"},
			expectedFindCityCalls:        nil,
			expectedGetForecastCalls:     nil,
			expectedGetWaveForecastCalls: nil,
			expectedNotifications:        nil,
		},
		{
			name:                         "city not found",
			userResult:                   testUser,
			cityError:                    ErrCityNotFound,
			expectedError:                ErrCityNotFound,
			expectedFindUserCalls:        []string{"user-id"},
			expectedFindCityCalls:        []string{"test city"},
			expectedGetForecastCalls:     nil,
			expectedGetWaveForecastCalls: nil,
			expectedNotifications:        nil,
		},
		{
			name:                         "unexpected error getting city",
			userResult:                   testUser,
			cityError:                    runtimeError,
			expectedError:                runtimeError,
			expectedFindUserCalls:        []string{"user-id"},
			expectedFindCityCalls:        []string{"test city"},
			expectedGetForecastCalls:     nil,
			expectedGetWaveForecastCalls: nil,
			expectedNotifications:        nil,
		},
		{
			name:                         "unexpected error getting weather forecast",
			userResult:                   testUser,
			cityResult:                   testCity,
			weatherError:                 runtimeError,
			expectedError:                runtimeError,
			expectedFindUserCalls:        []string{"user-id"},
			expectedFindCityCalls:        []string{"test city"},
			expectedGetForecastCalls:     []string{"city-id"},
			expectedGetWaveForecastCalls: nil,
			expectedNotifications:        nil,
		},
		{
			name:                         "unexpected error getting wave forecast",
			userResult:                   testUser,
			cityResult:                   testCity,
			weatherResult:                testWeatherForecast,
			waveError:                    runtimeError,
			expectedError:                runtimeError,
			expectedFindUserCalls:        []string{"user-id"},
			expectedFindCityCalls:        []string{"test city"},
			expectedGetForecastCalls:     []string{"city-id"},
			expectedGetWaveForecastCalls: []string{"city-id"},
			expectedNotifications:        nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockClient{
				findUserResult: tt.userResult,
				findUserError:  tt.userError,

				findCityResult: tt.cityResult,
				findCityError:  tt.cityError,

				getForecastResult: tt.weatherResult,
				getForecastError:  tt.weatherError,

				getWaveForecastResult: tt.waveResult,
				getWaveForecastError:  tt.waveError,
			}

			service := NewService(mock, mock, mock, mock, mock)

			err := service.NotifyUser("user-id", "test city")

			assert.True(t, errors.Is(err, tt.expectedError), fmt.Sprintf("Expected: %s / Actual: %s", tt.expectedError, err))

			assert.Equal(t, tt.expectedFindUserCalls, mock.findUserCalls)
			assert.Equal(t, tt.expectedFindCityCalls, mock.findCityCalls)
			assert.Equal(t, tt.expectedGetForecastCalls, mock.getForecastCalls)
			assert.Equal(t, tt.expectedGetWaveForecastCalls, mock.getWaveForecastCalls)
			assert.Equal(t, tt.expectedNotifications, mock.notifyCallsContent)

			for i, _ := range tt.expectedNotifications {
				assert.Equal(t, "user-id", mock.notifyCallsUserID[i])
			}
		})
	}
}
