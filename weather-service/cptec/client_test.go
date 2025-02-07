package cptec

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
		name              string
		cptecResponseCode int
		cptecResponse     string
		expectedResult    City
		expectedError     error
	}{
		{
			name:              "success",
			cptecResponseCode: 200,
			cptecResponse:     "<?xml version='1.0' encoding='ISO-8859-1'?><cidades><cidade><nome>Test City</nome><uf>XY</uf><id>123</id></cidade></cidades>",
			expectedResult: City{
				ID:    "123",
				Name:  "Test City",
				State: "XY",
			},
			expectedError: nil,
		},
		{
			name:              "city not found",
			cptecResponseCode: 200,
			cptecResponse:     "<?xml version='1.0' encoding='ISO-8859-1'?><cidades></cidades>",
			expectedResult:    City{},
			expectedError:     ErrCityNotFound,
		},
		{
			name:              "api returns multiple cities",
			cptecResponseCode: 200,
			cptecResponse:     "<?xml version='1.0' encoding='ISO-8859-1'?><cidades><cidade><nome>Test City</nome><uf>XY</uf><id>123</id></cidade><cidade><nome>Another City</nome><uf>XY</uf><id>456</id></cidade></cidades>",
			expectedResult:    City{},
			expectedError:     ErrMultipleCities,
		},
		{
			name:              "cptec api returns error",
			cptecResponseCode: 500,
			cptecResponse:     "Internal Server Error",
			expectedResult:    City{},
			expectedError:     ErrFetchingCities,
		},
		{
			name:              "malformed xml",
			cptecResponseCode: 200,
			cptecResponse:     "<?xml version='1.0' encoding='ISO-8859-1'?><cidades><cidade></cidades>",
			expectedResult:    City{},
			expectedError:     ErrReadingResponse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/XML/listaCidades?city=test+city", r.URL.String())
				w.WriteHeader(tt.cptecResponseCode)
				w.Write([]byte(tt.cptecResponse))
			}))

			defer server.Close()

			client := NewClient(server.Client(), server.URL)

			result, err := client.FindCity("test city")

			assert.True(t, errors.Is(err, tt.expectedError), fmt.Sprintf("Expected: %s / Actual: %s", tt.expectedError, err))
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
