package user

type UserTO struct {
	ID                 string               `json:"id"`
	Name               string               `json:"name"`
	NotificationConfig NotificationConfigTO `json:"notification"`
}

type NotificationConfigTO struct {
	Enabled bool                    `json:"enabled"`
	Web     WebNotificationConfigTO `json:"web"`
}

type WebNotificationConfigTO struct {
	Enabled bool   `json:"enabled"`
	ID      string `json:"id"`
}
