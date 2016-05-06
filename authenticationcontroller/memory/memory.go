package memory

import (
	"errors"

	"github.com/clawio/authentication/authenticationcontroller"
	"github.com/clawio/authentication/lib"
	"github.com/clawio/entities"
)

type User struct {
	*entities.User
	Password string `json:"password"`
}

// Options  holds the configuration
// parameters used by the MemoryAuthenticationController.
type Options struct {
	Users                    []*User
	JWTKey, JWTSigningMethod string
}

// New returns an AuthenticationControler that
// stores users in memory. This controller is for testing purposes.
func New(opts *Options) authenticationcontroller.AuthenticationController {
	return &controller{
		users:         opts.Users,
		authenticator: lib.NewAuthenticator(opts.JWTKey, opts.JWTSigningMethod),
	}
}

func (c *controller) Authenticate(username, password string) (string, error) {
	for _, u := range c.users {
		if u.Username == username && u.Password == password {
			return c.authenticator.CreateToken(u.User)
		}
	}
	return "", errors.New("user not found")
}

type controller struct {
	users         []*User
	authenticator *lib.Authenticator
}
