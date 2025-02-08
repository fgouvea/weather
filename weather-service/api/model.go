package api

type NotifyUserRequest struct {
	UserID string `json:"userId"`
	City   string `json:"city"`
}
