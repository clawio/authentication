package lib

import (
	"errors"
	"github.com/clawio/entities"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const DefaultJWTKey = "secret"
const DefaultJWTSigningMethod = "HS256"

type Authenticator struct {
	JWTKey           string
	JWTSigningMethod string
}

func NewAuthenticator(key, method string) *Authenticator {
	if key == "" {
		key = DefaultJWTKey
	}
	if method == "" {
		method = DefaultJWTSigningMethod
	}
	return &Authenticator{JWTKey: key, JWTSigningMethod: method}
}

func (a *Authenticator) CreateToken(user *entities.User) (string, error) {
	if user == nil {
		return "", errors.New("user is nil")
	}
	token := jwt.New(jwt.GetSigningMethod(a.JWTSigningMethod))
	token.Claims["username"] = user.Username
	token.Claims["email"] = user.Email
	token.Claims["display_name"] = user.DisplayName
	token.Claims["exp"] = time.Now().Add(time.Second * 3600).UnixNano()
	return token.SignedString([]byte(a.JWTKey))
}

func (a *Authenticator) CreateUserFromToken(token string) (*entities.User, error) {
	rawToken, err := a.parseToken(token)
	if err != nil {
		return nil, err
	}
	return a.getUserFromRawToken(rawToken)
}

func (a *Authenticator) getUserFromRawToken(rawToken *jwt.Token) (*entities.User, error) {
	username, ok := rawToken.Claims["username"].(string)
	if !ok {
		return nil, errors.New("token username claim failed cast to string")
	}

	email, ok := rawToken.Claims["email"].(string)
	if !ok {
		return nil, errors.New("token email claim failed cast to string")
	}

	displayName, ok := rawToken.Claims["display_name"].(string)
	if !ok {
		return nil, errors.New("token display_name claim failed cast to string")
	}
	return &entities.User{
		Username:    username,
		Email:       email,
		DisplayName: displayName,
	}, nil
}
func (a *Authenticator) parseToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.JWTKey), nil
	})
}
