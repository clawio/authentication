package mock

import (
	"github.com/clawio/entities"
	"github.com/stretchr/testify/mock"
)

type MockAuthenticationController struct {
	mock.Mock
}

func (m *MockAuthenticationController) Authenticate(username, password string) (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockAuthenticationController) Verify(token string) (entities.User, error) {
	args := m.Called()
	return args.Get(0).(entities.User), args.Error(1)
}
func (m *MockAuthenticationController) Invalidate(token string) error {
	args := m.Called()
	return args.Error(0)
}
