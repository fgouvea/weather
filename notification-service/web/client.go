package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/fgouvea/weather/notification-service/notification"
	"github.com/fgouvea/weather/notification-service/user"
)

const (
	sendNotificationPath = "/notification/send"
)

var (
	ErrFailedToSend = errors.New("failed to send notification to web api")
)

type ExternalNotification struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

type Client struct {
	Client *http.Client

	sendNotificationURL string
}

func NewClient(httpClient *http.Client, host string) *Client {
	return &Client{
		Client: httpClient,

		sendNotificationURL: fmt.Sprintf("%s%s", host, sendNotificationPath),
	}
}

func (c *Client) Send(recipient user.User, content string) error {
	if !recipient.NotificationConfig.Web.Enabled {
		return notification.ErrUserOptOut
	}

	notification := ExternalNotification{
		ID:      recipient.NotificationConfig.Web.ID,
		Content: content,
	}

	body, err := json.Marshal(notification)

	if err != nil {
		return fmt.Errorf("failed to encode notification: %w", err)
	}

	response, err := c.Client.Post(c.sendNotificationURL, "application/json", bytes.NewReader(body))

	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToSend, err)
	}

	if response.StatusCode != http.StatusAccepted {
		return fmt.Errorf("%w: unexpected status code: %d", ErrFailedToSend, response.StatusCode)
	}

	return nil
}
