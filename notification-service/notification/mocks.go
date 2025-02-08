package notification

import "github.com/fgouvea/weather/notification-service/user"

type userFinderMock struct {
	findUserCalls  []string
	findUserResult user.User
	findUserError  error
}

func (m *userFinderMock) FindUser(id string) (user.User, error) {
	m.findUserCalls = append(m.findUserCalls, id)
	return m.findUserResult, m.findUserError
}

type senderMock struct {
	sendCallsRecipient []user.User
	sendCallsContent   []string
	sendError          error
}

func (m *senderMock) Send(recipient user.User, content string) error {
	m.sendCallsRecipient = append(m.sendCallsRecipient, recipient)
	m.sendCallsContent = append(m.sendCallsContent, content)
	return m.sendError
}
