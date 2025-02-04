package test

import "github.com/fgouvea/weather/user-service/user"

type MockRepository struct {
	users map[string]*user.User
}

func NewMockRepository() user.Repository {
	return &MockRepository{users: map[string]*user.User{}}
}

func (r *MockRepository) Find(id string) (*user.User, error) {
	if user, exists := r.users[id]; exists {
		return user, nil
	}

	return nil, user.UserNotFound
}

func (r *MockRepository) Save(user *user.User) error {
	r.users[user.Id] = user
	return nil
}
