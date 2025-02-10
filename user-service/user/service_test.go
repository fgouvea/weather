package user

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserService_Create_Success(t *testing.T) {
	tests := []struct {
		name                  string
		userName              string
		userWebNotificationId string
		expectedResult        User
	}{
		{
			name:                  "with web notification id",
			userName:              "Fulano Beltrano",
			userWebNotificationId: "123",
			expectedResult: User{
				Name: "Fulano Beltrano",
				NotificationConfig: NotificationConfig{
					Enabled: true,
					Web: WebNotificationConfig{
						Enabled: true,
						Id:      "123",
					},
				},
			},
		},
		{
			name:                  "without web notification id",
			userName:              "Fulano Beltrano",
			userWebNotificationId: "",
			expectedResult: User{
				Name: "Fulano Beltrano",
				NotificationConfig: NotificationConfig{
					Enabled: true,
					Web: WebNotificationConfig{
						Enabled: false,
						Id:      "",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repositoryMock := &MockRepository{
				SaveCalls: []*User{},
				SaveError: nil,
			}

			service := NewService(repositoryMock, repositoryMock)

			result, err := service.Create(tt.userName, tt.userWebNotificationId)

			assert.Nil(t, err)
			assert.NotNil(t, result)

			assert.Equal(t, tt.expectedResult.Name, result.Name)
			assert.Equal(t, tt.expectedResult.NotificationConfig, result.NotificationConfig)

			assert.Len(t, repositoryMock.SaveCalls, 1)
			assert.Equal(t, repositoryMock.SaveCalls[0], result)
		})
	}
}

func TestUserService_Create_Error(t *testing.T) {
	repositoryMock := &MockRepository{
		SaveCalls: []*User{},
		SaveError: fmt.Errorf("runtime error"),
	}

	service := NewService(repositoryMock, repositoryMock)

	result, err := service.Create("Fulano Beltrano", "123")

	assert.Nil(t, result)
	assert.EqualError(t, err, "unexpected error saving user: runtime error")
}

func TestUserService_Find_Success(t *testing.T) {
	savedUser := User{
		ID:   "USER-1",
		Name: "Fulano Beltrano",
		NotificationConfig: NotificationConfig{
			Enabled: true,
			Web: WebNotificationConfig{
				Enabled: true,
				Id:      "123",
			},
		},
	}

	repositoryMock := &MockRepository{
		FindCalls:  []string{},
		FindResult: savedUser,
		FindError:  nil,
	}

	service := NewService(repositoryMock, repositoryMock)

	result, err := service.Find("USER-1")

	assert.Nil(t, err)

	assert.Equal(t, result, &savedUser)

	assert.Len(t, repositoryMock.FindCalls, 1)
	assert.Equal(t, repositoryMock.FindCalls[0], "USER-1")
}

func TestUserService_Find_Error(t *testing.T) {
	tests := []struct {
		name          string
		findError     error
		expectedError error
	}{
		{
			name:          "user not found",
			findError:     ErrUserNotFound,
			expectedError: ErrUserNotFound,
		},
		{
			name:          "unexpected error",
			findError:     fmt.Errorf("error connecting to database"),
			expectedError: fmt.Errorf("unexpected error fetching user: error connecting to database"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repositoryMock := &MockRepository{
				FindCalls: []string{},
				FindError: tt.findError,
			}

			service := NewService(repositoryMock, repositoryMock)

			result, err := service.Find("USER-1")

			assert.Nil(t, result)

			assert.EqualError(t, err, tt.expectedError.Error())

			assert.Len(t, repositoryMock.FindCalls, 1)
			assert.Equal(t, repositoryMock.FindCalls[0], "USER-1")
		})
	}
}

func TestUserService_OptOutOf_Notifications(t *testing.T) {
	tests := []struct {
		name              string
		findError         error
		saveError         error
		expectedError     error
		expectedSaveCalls int
	}{
		{
			name:              "success",
			findError:         nil,
			saveError:         nil,
			expectedError:     nil,
			expectedSaveCalls: 1,
		},
		{
			name:              "error saving user",
			findError:         nil,
			saveError:         fmt.Errorf("error connecting to database"),
			expectedError:     fmt.Errorf("unexpected error saving user: error connecting to database"),
			expectedSaveCalls: 1,
		},
		{
			name:              "user not found",
			findError:         ErrUserNotFound,
			saveError:         nil,
			expectedError:     ErrUserNotFound,
			expectedSaveCalls: 0,
		},
		{
			name:              "error fetching user",
			findError:         fmt.Errorf("error connecting to database"),
			saveError:         nil,
			expectedError:     fmt.Errorf("unexpected error fetching user: error connecting to database"),
			expectedSaveCalls: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repositoryMock := &MockRepository{
				FindCalls: []string{},
				FindError: tt.findError,
				SaveCalls: []*User{},
				SaveError: tt.saveError,
			}

			if tt.findError == nil {
				repositoryMock.FindResult = User{
					ID:   "USER-1",
					Name: "Fulano Beltrano",
					NotificationConfig: NotificationConfig{
						Enabled: true,
						Web: WebNotificationConfig{
							Enabled: true,
							Id:      "123",
						},
					},
				}
			}

			service := NewService(repositoryMock, repositoryMock)

			err := service.OptOutOfNotifications("USER-1")

			assert.Equal(t, repositoryMock.FindCalls, []string{"USER-1"})

			assert.Equal(t, tt.expectedError == nil, err == nil)
			if err != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			}

			assert.Equal(t, tt.expectedSaveCalls, len(repositoryMock.SaveCalls))
			if tt.expectedSaveCalls > 0 {
				assert.Equal(t, false, repositoryMock.SaveCalls[0].NotificationConfig.Enabled)
			}
		})
	}
}
