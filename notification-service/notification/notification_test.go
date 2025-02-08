package notification

import (
	"errors"
	"fmt"
	"testing"

	"github.com/fgouvea/weather/notification-service/user"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestService_Process(t *testing.T) {
	tests := []struct {
		name                string
		senders             map[string]*senderMock
		userResult          user.User
		userError           error
		expectedError       error
		expectedCallSenders bool
	}{
		{
			name: "success single sender",
			senders: map[string]*senderMock{
				"sender-1": &senderMock{sendError: nil},
			},
			userResult:          user.User{ID: "USER-123", NotificationConfig: user.NotificationConfig{Enabled: true}},
			userError:           nil,
			expectedError:       nil,
			expectedCallSenders: true,
		},
		{
			name: "success multiple senders",
			senders: map[string]*senderMock{
				"sender-1": &senderMock{sendError: nil},
				"sender-2": &senderMock{sendError: fmt.Errorf("runtime error")},
				"sender-3": &senderMock{sendError: ErrUserOptOut},
			},
			userResult:          user.User{ID: "USER-123", NotificationConfig: user.NotificationConfig{Enabled: true}},
			userError:           nil,
			expectedError:       nil,
			expectedCallSenders: true,
		},
		{
			name: "user opted out",
			senders: map[string]*senderMock{
				"sender-1": &senderMock{sendError: nil},
			},
			userResult:          user.User{ID: "USER-123", NotificationConfig: user.NotificationConfig{Enabled: false}},
			userError:           nil,
			expectedError:       nil,
			expectedCallSenders: false,
		},
		{
			name: "error fetching user",
			senders: map[string]*senderMock{
				"sender-1": &senderMock{sendError: nil},
			},
			userResult:          user.User{ID: "USER-123", NotificationConfig: user.NotificationConfig{Enabled: false}},
			userError:           fmt.Errorf("runtime error"),
			expectedError:       ErrFailedToProcess,
			expectedCallSenders: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userFinderMock := &userFinderMock{
				findUserResult: tt.userResult,
				findUserError:  tt.userError,
			}

			logger, _ := zap.NewDevelopment()

			senders := map[string]Sender{}
			for senderName, sender := range tt.senders {
				senders[senderName] = sender
			}

			service := NewService(userFinderMock, senders, logger)

			err := service.Process(Notification{
				UserID:  "USER-123",
				Content: "test notification content",
			})

			assert.True(t, errors.Is(err, tt.expectedError), fmt.Sprintf("Expected: %s / Actual: %s", tt.expectedError, err))

			assert.Equal(t, []string{"USER-123"}, userFinderMock.findUserCalls)

			for _, sender := range tt.senders {
				if tt.expectedCallSenders {
					assert.Equal(t, []user.User{tt.userResult}, sender.sendCallsRecipient)
					assert.Equal(t, []string{"test notification content"}, sender.sendCallsContent)
				} else {
					assert.Equal(t, 0, len(sender.sendCallsRecipient))
				}
			}
		})
	}
}
