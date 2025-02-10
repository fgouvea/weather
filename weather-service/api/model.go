package api

type NotifyUserRequest struct {
	UserID string `json:"userId"`
	City   string `json:"city"`
}

type ScheduleRequest struct {
	UserID string `json:"userId"`
	City   string `json:"city"`
	Time   string `json:"time"`
}
