package weather

type City struct {
	ID    string
	Name  string
	State string
}

type CityForecast struct {
	UpdatedAt string
	Forecast  []Forecast
}

type Forecast struct {
	Date           string
	Weather        string
	MaxTemperature int
	MinTemperature int
	IUV            float64
}

type CityWaveForecast struct {
	UpdatedAt string
	Date      string
	Morning   WaveForecast
	Afternoon WaveForecast
	Evening   WaveForecast
}

type WaveForecast struct {
	Swell         string
	Height        float64
	Wind          float64
	WaveDirection string
	WindDirection string
}
