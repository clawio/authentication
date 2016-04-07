package authenticationcontroller

import "github.com/clawio/entities"

type AuthenticationController interface {
	Authenticate(username, password string) (string, error)
	Verify(token string) (entities.User, error)
	Invalidate(token string) error
}
