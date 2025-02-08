package user

import "errors"

var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID   string
	Name string
}
