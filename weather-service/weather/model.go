package weather

type City struct {
	ID    string
	Name  string
	State string
}

type CityForecast struct {
	Name     string
	State    string
	Date     string
	Forecast []Forecast
}

type Forecast struct {
	Date           string
	Weather        string
	MaxTemperature int
	MinTemperature int
	IUV            float64
}
