package api

import "github.com/fgouvea/weather/user-service/user"

type CreateUserRequestTO struct {
	Name              string `json:"name"`
	WebNotificationID string `json:"webNotificationId"`
}

type UserTO struct {
	Id                 string               `json:"id"`
	Name               string               `json:"name"`
	NotificationConfig NotificationConfigTO `json:"notification"`
}

type NotificationConfigTO struct {
	Enabled bool                    `json:"enabled"`
	Web     WebNotificationConfigTO `json:"web"`
}

type WebNotificationConfigTO struct {
	Enabled bool   `json:"enabled"`
	Id      string `json:"id"`
}

func buildUserTO(u *user.User) UserTO {
	return UserTO{
		Id:   u.Id,
		Name: u.Name,
		NotificationConfig: NotificationConfigTO{
			Enabled: u.NotificationConfig.Enabled,
			Web: WebNotificationConfigTO{
				Enabled: u.NotificationConfig.Web.Enabled,
				Id:      u.NotificationConfig.Web.Id,
			},
		},
	}
}
