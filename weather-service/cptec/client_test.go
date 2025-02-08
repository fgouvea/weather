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
			cptecResponse:     "<?xml version='1.0' encoding='ISO-8859-1'?><cidades><cidade><nome>Test City A</nome><uf>XY</uf><id>123</id></cidade><cidade><nome>Test City B</nome><uf>XY</uf><id>456</id></cidade></cidades>",
			expectedResult:    weather.City{},
			expectedError:     ErrMultipleCities,
		},
		{
			name:              "exact match among multiple cities",
			cptecResponseCode: 200,
			cptecResponse:     "<?xml version='1.0' encoding='ISO-8859-1'?><cidades><cidade><nome>Test City Other</nome><uf>XY</uf><id>456</id></cidade><cidade><nome>Test City</nome><uf>XY</uf><id>123</id></cidade></cidades>",
			expectedResult: weather.City{
				ID:    "123",
				Name:  "Test City",
				State: "XY",
			},
			expectedError: nil,
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
				UpdatedAt: "2025-02-07",
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

func TestClient_GetWaveForecast(t *testing.T) {
	tests := []struct {
		name              string
		cptecResponseCode int
		cptecResponse     string
		expectedResult    weather.CityWaveForecast
		expectedError     error
	}{
		{
			name:              "success",
			cptecResponseCode: 200,
			cptecResponse:     "<?xml version='1.0' encoding='ISO-8859-1'?><cidade><nome>Rio de Janeiro</nome><uf>RJ</uf><atualizacao>07-02-2025</atualizacao><manha><dia>07-02-2025 12h Z</dia><agitacao>Fraco</agitacao><altura>0.5</altura><direcao>SE</direcao><vento>4.1</vento><vento_dir>ENE</vento_dir></manha><tarde><dia>07-02-2025 18h Z</dia><agitacao>Fraco</agitacao><altura>0.5</altura><direcao>SE</direcao><vento>5.2</vento><vento_dir>E</vento_dir></tarde><noite><dia>07-02-2025 21h Z</dia><agitacao>Fraco</agitacao><altura>0.5</altura><direcao>SE</direcao><vento>5.3</vento><vento_dir>E</vento_dir></noite></cidade>",
			expectedResult: weather.CityWaveForecast{
				UpdatedAt: "2025-02-07",
				Date:      "2025-02-07",
				Morning: weather.WaveForecast{
					Swell:         "Fraco",
					Height:        0.5,
					Wind:          4.1,
					WaveDirection: "SE",
					WindDirection: "ENE",
				},
				Afternoon: weather.WaveForecast{
					Swell:         "Fraco",
					Height:        0.5,
					Wind:          5.2,
					WaveDirection: "SE",
					WindDirection: "E",
				},
				Evening: weather.WaveForecast{
					Swell:         "Fraco",
					Height:        0.5,
					Wind:          5.3,
					WaveDirection: "SE",
					WindDirection: "E",
				},
			},
			expectedError: nil,
		},
		{
			name:              "city not found",
			cptecResponseCode: 200,
			cptecResponse:     "<?xml version='1.0' encoding='ISO-8859-1'?><cidade><nome>undefined</nome><uf>undefined</uf><atualizacao>00/00/0000 00:00:00</atualizacao><manha><dia>00/00/0000 00:00:00</dia><agitacao>undefined</agitacao><altura>undefined</altura><direcao>undefined</direcao><vento>undefined</vento><vento_dir>undefined</vento_dir></manha><tarde><dia>00/00/0000 00:00:00</dia><agitacao>undefined</agitacao><altura>undefined</altura><direcao>undefined</direcao><vento>undefined</vento><vento_dir>undefined</vento_dir></tarde><noite><dia>00/00/0000 00:00:00</dia><agitacao>undefined</agitacao><altura>undefined</altura><direcao>undefined</direcao><vento>undefined</vento><vento_dir>undefined</vento_dir></noite></cidade>",
			expectedResult:    weather.CityWaveForecast{},
			expectedError:     ErrCityNotFound,
		},
		{
			name:              "cptec api returns error",
			cptecResponseCode: 500,
			cptecResponse:     "Internal Server Error",
			expectedResult:    weather.CityWaveForecast{},
			expectedError:     ErrFetchingResponse,
		},
		{
			name:              "malformed xml",
			cptecResponseCode: 200,
			cptecResponse:     "<?xml version='1.0' encoding='ISO-8859-1'?><cidades><cidade></cidades>",
			expectedResult:    weather.CityWaveForecast{},
			expectedError:     ErrReadingResponse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/XML/cidade/123/dia/0/ondas.xml", r.URL.String())
				w.WriteHeader(tt.cptecResponseCode)
				w.Write([]byte(tt.cptecResponse))
			}))

			defer server.Close()

			client := NewClient(server.Client(), server.URL)

			result, err := client.GetWaveForecast("123")

			assert.True(t, errors.Is(err, tt.expectedError), fmt.Sprintf("Expected: %s / Actual: %s", tt.expectedError, err))
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
