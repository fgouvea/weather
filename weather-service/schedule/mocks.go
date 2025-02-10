package schedule

type userAndCity struct {
	userID   string
	cityName string
}

type serviceMock struct {
	validateCalls []userAndCity
	validateError error

	notifyCalls []userAndCity
	notifyError error

	saveCalls []Schedule
	saveError error
}

var _ Validator = (*serviceMock)(nil)
var _ ScheduleSaver = (*serviceMock)(nil)
var _ Notifier = (*serviceMock)(nil)

func (m *serviceMock) Validate(userID, cityName string) error {
	m.validateCalls = append(m.validateCalls, userAndCity{userID: userID, cityName: cityName})
	return m.validateError
}

func (m *serviceMock) NotifyUser(userID string, cityName string) error {
	m.notifyCalls = append(m.notifyCalls, userAndCity{userID: userID, cityName: cityName})
	return m.notifyError
}

func (m *serviceMock) Save(schedule Schedule) error {
	m.saveCalls = append(m.saveCalls, schedule)
	return m.saveError
}
