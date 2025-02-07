package cptec

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fgouvea/weather/weather-service/weather"
	"github.com/stretchr/testify/assert"
)

func TestClient_Find(t *testing.T) {
	tests := []struct {
		name              string
		cptecResponseCode int
		cptecResponse     string
		expectedResult    weather.City
		expectedError     error
	}{
		{
			name:              "success",
			cptecResponseCode: 200,
			cptecResponse:     "<?xml version='1.0' encoding='ISO-8859-1'?><cidades><cidade><nome>Test City</nome><uf>XY</uf><id>123</id></cidade></cidades>",
			expectedResult: weather.City{
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
			expectedResult:    weather.City{},
			expectedError:     ErrCityNotFound,
		},
		{
			name:              "api returns multiple cities",
			cptecResponseCode: 200,
			cptecResponse:     "<?xml version='1.0' encoding='ISO-8859-1'?><cidades><cidade><nome>Test City</nome><uf>XY</uf><id>123</id></cidade><cidade><nome>Another City</nome><uf>XY</uf><id>456</id></cidade></cidades>",
			expectedResult:    weather.City{},
			expectedError:     ErrMultipleCities,
		},
		{
			name:              "cptec api returns error",
			cptecResponseCode: 500,
			cptecResponse:     "Internal Server Error",
			expectedResult:    weather.City{},
			expectedError:     ErrFetchingResponse,
		},
		{
			name:              "malformed xml",
			cptecResponseCode: 200,
			cptecResponse:     "<?xml version='1.0' encoding='ISO-8859-1'?><cidades><cidade></cidades>",
			expectedResult:    weather.City{},
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

func TestClient_GetForecast(t *testing.T) {
	tests := []struct {
		name              string
		cptecResponseCode int
		cptecResponse     string
		expectedResult    weather.CityForecast
		expectedError     error
	}{
		{
			name:              "success",
			cptecResponseCode: 200,
			cptecResponse:     "<?xml version='1.0' encoding='ISO-8859-1'?><cidade><nome>Test City</nome><uf>XY</uf><atualizacao>2025-02-07</atualizacao><previsao><dia>2025-02-08</dia><tempo>pn</tempo><maxima>33</maxima><minima>21</minima><iuv>0.0</iuv></previsao><previsao><dia>2025-02-09</dia><tempo>c</tempo><maxima>33</maxima><minima>21</minima><iuv>0.0</iuv></previsao><previsao><dia>2025-02-10</dia><tempo>pn</tempo><maxima>35</maxima><minima>21</minima><iuv>0.0</iuv></previsao><previsao><dia>2025-02-11</dia><tempo>xyz</tempo><maxima>36</maxima><minima>22</minima><iuv>0.0</iuv></previsao></cidade>",
			expectedResult: weather.CityForecast{
				Name:  "Test City",
				State: "XY",
				Date:  "2025-02-07",
				Forecast: []weather.Forecast{
					{
						Date:           "2025-02-08",
						Weather:        "Parcialmente Nublado",
						MaxTemperature: 33,
						MinTemperature: 21,
						IUV:            0.0,
					},
					{
						Date:           "2025-02-09",
						Weather:        "Chuva",
						MaxTemperature: 33,
						MinTemperature: 21,
						IUV:            0.0,
					},
					{
						Date:           "2025-02-10",
						Weather:        "Parcialmente Nublado",
						MaxTemperature: 35,
						MinTemperature: 21,
						IUV:            0.0,
					},
					{
						Date:           "2025-02-11",
						Weather:        "Desconhecido",
						MaxTemperature: 36,
						MinTemperature: 22,
						IUV:            0.0,
					},
				},
			},
			expectedError: nil,
		},
		{
			name:              "city not found",
			cptecResponseCode: 200,
			cptecResponse:     "<?xml version='1.0' encoding='ISO-8859-1'?><cidade><nome>null</nome><uf>null</uf><atualizacao>null</atualizacao><previsao><dia>null</dia><tempo>null</tempo><maxima>null</maxima><minima>null</minima><iuv>0.0</iuv></previsao><previsao><dia>null</dia><tempo>null</tempo><maxima>null</maxima><minima>null</minima><iuv>0.0</iuv></previsao><previsao><dia>null</dia><tempo>null</tempo><maxima>null</maxima><minima>null</minima><iuv>0.0</iuv></previsao><previsao><dia>null</dia><tempo>null</tempo><maxima>null</maxima><minima>null</minima><iuv>0.0</iuv></previsao></cidade>",
			expectedResult:    weather.CityForecast{},
			expectedError:     ErrCityNotFound,
		},
		{
			name:              "cptec api returns error",
			cptecResponseCode: 500,
			cptecResponse:     "Internal Server Error",
			expectedResult:    weather.CityForecast{},
			expectedError:     ErrFetchingResponse,
		},
		{
			name:              "malformed xml",
			cptecResponseCode: 200,
			cptecResponse:     "<?xml version='1.0' encoding='ISO-8859-1'?><cidades><cidade></cidades>",
			expectedResult:    weather.CityForecast{},
			expectedError:     ErrReadingResponse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/XML/cidade/123/previsao.xml", r.URL.String())
				w.WriteHeader(tt.cptecResponseCode)
				w.Write([]byte(tt.cptecResponse))
			}))

			defer server.Close()

			client := NewClient(server.Client(), server.URL)

			result, err := client.GetForecast("123")

			assert.True(t, errors.Is(err, tt.expectedError), fmt.Sprintf("Expected: %s / Actual: %s", tt.expectedError, err))
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
