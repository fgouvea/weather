package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	getUserPath = "/user-service/user/%s"
)

var (
	ErrRequestAPI      = errors.New("error fetching user from api")
	ErrReadingResponse = errors.New("error reading api response")
)

type Client struct {
	Client *http.Client

	getUserURL string
}

func NewClient(httpClient *http.Client, basePath string) *Client {
	return &Client{
		Client: httpClient,

		getUserURL: fmt.Sprintf("%s%s", basePath, getUserPath),
	}
}

func (c *Client) FindUser(id string) (User, error) {
	response, err := c.Client.Get(fmt.Sprintf(c.getUserURL, id))

	if err != nil {
		return User{}, fmt.Errorf("%w: %w", ErrRequestAPI, err)
	}

	if response.StatusCode == http.StatusNotFound {
		return User{}, ErrUserNotFound
	}

	if response.StatusCode != 200 {
		return User{}, fmt.Errorf("%w: unexpected status code: %d", ErrRequestAPI, response.StatusCode)
	}

	var parsedResponse UserTO

	err = json.NewDecoder(response.Body).Decode(&parsedResponse)

	if err != nil {
		return User{}, fmt.Errorf("%w: %w", ErrReadingResponse, err)
	}

	return User{
		ID:   parsedResponse.ID,
		Name: parsedResponse.Name,
		NotificationConfig: NotificationConfig{
			Enabled: parsedResponse.NotificationConfig.Enabled,
			Web: WebNotificationConfig{
				Enabled: parsedResponse.NotificationConfig.Web.Enabled,
				ID:      parsedResponse.NotificationConfig.Web.ID,
			},
		},
	}, nil
}
