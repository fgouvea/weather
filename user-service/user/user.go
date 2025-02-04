package user

type User struct {
	Id                 string
	Name               string
	NotificationConfig NotificationConfig
}

type NotificationConfig struct {
	Enabled bool
	Web     WebNotificationConfig
}

type WebNotificationConfig struct {
	Enabled bool
	Id      string
}
