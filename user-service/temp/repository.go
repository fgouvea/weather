package temp

import "github.com/fgouvea/weather/user-service/user"

type InMemoryUserRepository struct {
	users map[string]*user.User
}

func NewMockRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{users: map[string]*user.User{}}
}

var _ user.Saver = (*InMemoryUserRepository)(nil)
var _ user.Finder = (*InMemoryUserRepository)(nil)

func (r *InMemoryUserRepository) Find(id string) (*user.User, error) {
	if user, exists := r.users[id]; exists {
		return user, nil
	}

	return nil, user.ErrUserNotFound
}

func (r *InMemoryUserRepository) Save(user *user.User) error {
	r.users[user.Id] = user
	return nil
}
