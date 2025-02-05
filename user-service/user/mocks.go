package user

type MockRepository struct {
	FindCalls  []string
	FindResult User
	FindError  error

	SaveCalls []*User
	SaveError error
}

var _ Saver = (*MockRepository)(nil)
var _ Finder = (*MockRepository)(nil)

func (r *MockRepository) Find(id string) (*User, error) {
	r.FindCalls = append(r.FindCalls, id)
	return &r.FindResult, r.FindError
}

func (r *MockRepository) Save(user *User) error {
	r.SaveCalls = append(r.SaveCalls, user)
	return r.SaveError
}
