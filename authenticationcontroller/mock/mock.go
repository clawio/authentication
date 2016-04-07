package mock

import (
	"github.com/clawio/entities"
	"github.com/stretchr/testify/mock"
)

// AuthenticationController mocks an AuthenticationController for testing purposes.
type AuthenticationController struct {
	mock.Mock
}

// Authenticate mocks the Authenticate call.
func (m *AuthenticationController) Authenticate(username, password string) (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

// Verify mocks the Verify call.
func (m *AuthenticationController) Verify(token string) (entities.User, error) {
	args := m.Called()
	return args.Get(0).(entities.User), args.Error(1)
}

// Invalidate mocks the Invalidate call.
func (m *AuthenticationController) Invalidate(token string) error {
	args := m.Called()
	return args.Error(0)
}
