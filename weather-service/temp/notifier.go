package temp

import (
	"fmt"

	"github.com/fgouvea/weather/weather-service/weather"
)

type TempNotifier struct{}

var _ weather.Notifier = (*TempNotifier)(nil)

func (*TempNotifier) Notify(userID string, content string) error {
	fmt.Printf("User: %s\n", userID)
	fmt.Println("---------")
	fmt.Println(content)
	fmt.Println("---------")
	return nil
}
