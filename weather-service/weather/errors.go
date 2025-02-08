package weather

import "errors"

var ErrCityNotFound = errors.New("city not found")
var ErrMultipleCities = errors.New("multiple cities found with name")
