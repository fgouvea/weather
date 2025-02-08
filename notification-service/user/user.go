package user

import "errors"

var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID                 string
	Name               string
	NotificationConfig NotificationConfig
}

type NotificationConfig struct {
	Enabled bool
	Web     WebNotificationConfig
}

type WebNotificationConfig struct {
	Enabled bool
	ID      string
}
