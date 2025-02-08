package user

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_Find(t *testing.T) {
	tests := []struct {
		name            string
		apiResponseCode int
		apiResponse     string
		expectedResult  User
		expectedError   error
	}{
		{
			name:            "success",
			apiResponseCode: 200,
			apiResponse:     `{"id":"USER-123","name":"Fulano","notification":{"enabled":true,"web":{"enabled":false,"id":""}}}`,
			expectedResult: User{
				ID:   "USER-123",
				Name: "Fulano",
			},
			expectedError: nil,
		},
		{
			name:            "user not found",
			apiResponseCode: 404,
			apiResponse:     "",
			expectedResult:  User{},
			expectedError:   ErrUserNotFound,
		},
		{
			name:            "error calling api",
			apiResponseCode: 500,
			apiResponse:     "",
			expectedResult:  User{},
			expectedError:   ErrRequestAPI,
		},
		{
			name:            "malformed json",
			apiResponseCode: 200,
			apiResponse:     `{"id":"USER-123","name":"Fulano","notification":{"enabled":true,"`,
			expectedResult:  User{},
			expectedError:   ErrReadingResponse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/user-service/user/USER-123", r.URL.String())
				w.WriteHeader(tt.apiResponseCode)
				w.Write([]byte(tt.apiResponse))
			}))

			defer server.Close()

			client := NewClient(server.Client(), server.URL)

			result, err := client.FindUser("USER-123")

			assert.True(t, errors.Is(err, tt.expectedError), fmt.Sprintf("Expected: %s / Actual: %s", tt.expectedError, err))
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
