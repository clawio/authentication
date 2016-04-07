package authenticationcontroller

import "github.com/clawio/entities"

// AuthenticationController defines an interface to authenticate users,
// verify their tokens and invalidate them.
type AuthenticationController interface {
	Authenticate(username, password string) (string, error)
	Verify(token string) (entities.User, error)
	Invalidate(token string) error
}
