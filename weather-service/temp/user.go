package temp

import (
	"github.com/fgouvea/weather/weather-service/user"
	"github.com/fgouvea/weather/weather-service/weather"
)

type TempUserClient struct{}

var _ weather.UserFinder = (*TempUserClient)(nil)

func (*TempUserClient) FindUser(id string) (user.User, error) {
	return user.User{
		ID:   "USER-123456",
		Name: "Fernando",
	}, nil
}
