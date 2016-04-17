package memory

import (
	"errors"
	"os"
	"time"

	"github.com/clawio/authentication/authenticationcontroller"
	"github.com/clawio/entities"
	"github.com/dgrijalva/jwt-go"
)

// User implements github.com/clawio/entities.User interface.
// It is exported because it helps to decode users from a JSON file.
type User struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Password    string `json:"password"`
}

// GetUsername returns the username.
func (u *User) GetUsername() string { return u.Username }

// GetEmail returns the email.
func (u *User) GetEmail() string { return u.Email }

// GetDisplayName returns the DisplayName.
func (u *User) GetDisplayName() string { return u.DisplayName }

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
		users:            opts.Users,
		jwtKey:           opts.JWTKey,
		jwtSigningMethod: opts.JWTSigningMethod,
	}
}

func (c *controller) Authenticate(username, password string) (string, error) {
	for _, u := range c.users {
		if u.Username == username && u.Password == password {
			return c.createToken(u)
		}
	}
	return "", errors.New("user not found")
}

// Verify checks id token is valid and creates an user user from it.
func (c *controller) Verify(token string) (entities.User, error) {
	t, err := c.parseToken(token)
	if err != nil {
		return nil, err
	}
	return c.createUserFromToken(t)
}

func (c *controller) Invalidate(token string) error {
	return nil
}

type controller struct {
	users            []*User
	jwtKey           string // the key to sign the token
	jwtSigningMethod string // the algo to sign the token
}

func (c *controller) createToken(user entities.User) (string, error) {
	if user == nil {
		return "", errors.New("user is nil")
	}
	token := jwt.New(jwt.GetSigningMethod(c.jwtSigningMethod))
	host, _ := os.Hostname()
	token.Claims["username"] = user.GetUsername()
	token.Claims["email"] = user.GetEmail()
	token.Claims["display_name"] = user.GetDisplayName()
	token.Claims["iss"] = host
	token.Claims["exp"] = time.Now().Add(time.Second * 3600).UnixNano()
	return token.SignedString([]byte(c.jwtKey))
}

func (c *controller) parseToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (key interface{}, err error) {
		return []byte(c.jwtKey), nil
	})
}

func (c *controller) createUserFromToken(token *jwt.Token) (entities.User, error) {
	username, ok := token.Claims["username"].(string)
	if !ok {
		return nil, errors.New("token username claim failed cast to string")
	}

	email, ok := token.Claims["email"].(string)
	if !ok {
		return nil, errors.New("token email claim failed cast to string")
	}

	displayName, ok := token.Claims["display_name"].(string)
	if !ok {
		return nil, errors.New("token display_name claim failed cast to string")
	}

	return &User{
		Username:    username,
		Email:       email,
		DisplayName: displayName,
	}, nil
}
