package user

import (
	"errors"
)

var UserNotFound = errors.New("user not found")

type Repository interface {
	Save(user *User) error
	Find(id string) (*User, error)
}
