package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fgouvea/weather/notification-service/notification"
	"github.com/fgouvea/weather/notification-service/user"
	"github.com/stretchr/testify/assert"
)

func TestClient_Send(t *testing.T) {
	tests := []struct {
		name            string
		enabled         bool
		apiResponseCode int
		expectedError   error
	}{
		{
			name:            "success",
			enabled:         true,
			apiResponseCode: 202,
			expectedError:   nil,
		},
		{
			name:            "api error",
			enabled:         true,
			apiResponseCode: 500,
			expectedError:   ErrFailedToSend,
		},
		{
			name:          "web notification disabled",
			enabled:       false,
			expectedError: notification.ErrUserOptOut,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/notification/send", r.URL.String())

				var body ExternalNotification

				err := json.NewDecoder(r.Body).Decode(&body)

				if err != nil {
					assert.Fail(t, fmt.Sprintf("failed to decode body: %s", err.Error()))
				}

				assert.Equal(t, "USER-123", body.ID)
				assert.Equal(t, "test notification content", body.Content)

				w.WriteHeader(tt.apiResponseCode)
			}))

			defer server.Close()

			client := NewClient(server.Client(), server.URL)

			recipient := user.User{
				NotificationConfig: user.NotificationConfig{
					Enabled: true,
					Web: user.WebNotificationConfig{
						Enabled: tt.enabled,
						ID:      "USER-123",
					},
				},
			}

			err := client.Send(recipient, "test notification content")

			assert.True(t, errors.Is(err, tt.expectedError), fmt.Sprintf("Expected: %s / Actual: %s", tt.expectedError, err))
		})
	}
}
