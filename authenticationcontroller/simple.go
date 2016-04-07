package authenticationcontroller

import (
	"errors"
	"os"
	"time"

	"github.com/clawio/entities"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql" // enable mysql driver
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"           // enable postgresql driver
	_ "github.com/mattn/go-sqlite3" // enable sqlite3 driver
)

type simpleAuthenticationController struct {
	driver, dsn      string
	db               *gorm.DB
	jwtKey           string // the key to sign the token
	jwtSigningMethod string // the algo to sign the token
}

type SimpleAuthenticationControllerOptions struct {
	Driver, DSN              string
	JWTKey, JWTSigningMethod string
}

func NewSimpleAuthenticationController(opts *SimpleAuthenticationControllerOptions) (AuthenticationController, error) {
	db, err := gorm.Open(opts.Driver, opts.DSN)
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&userRecord{}).Error
	if err != nil {
		return nil, err
	}
	return &simpleAuthenticationController{
		driver:           opts.Driver,
		dsn:              opts.DSN,
		db:               db,
		jwtKey:           opts.JWTKey,
		jwtSigningMethod: opts.JWTSigningMethod,
	}, nil
}

func (c *simpleAuthenticationController) Authenticate(username, password string) (string, error) {
	rec, err := c.findByCredentials(username, password)
	if err != nil {
		return "", err
	}
	return c.createToken(rec)
}

// Verify checks id token is valid and creates an user user from it.
func (c *simpleAuthenticationController) Verify(token string) (entities.User, error) {
	t, err := c.parseToken(token)
	if err != nil {
		return nil, err
	}
	return c.createUserFromToken(t)
}

func (c *simpleAuthenticationController) Invalidate(token string) error {
	return nil
}

// findByCredentials finds an user given an username and a password.
func (c *simpleAuthenticationController) findByCredentials(username, password string) (*userRecord, error) {
	rec := &userRecord{}
	err := c.db.Where("username=? AND password=?", username, password).First(rec).Error
	return rec, err
}

// TODO(labkode) set collation for table and column to utf8. The default is swedish
type userRecord struct {
	Username    string `gorm:"primary_key" json:"username"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Password    string `json:"-"`
}

func (u userRecord) TableName() string {
	return "users"
}

func (u *userRecord) GetUsername() string    { return u.Username }
func (u *userRecord) GetEmail() string       { return u.Email }
func (u *userRecord) GetDisplayName() string { return u.DisplayName }

func (c *simpleAuthenticationController) createToken(user entities.User) (string, error) {
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

func (c *simpleAuthenticationController) parseToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (key interface{}, err error) {
		return []byte(c.jwtKey), nil
	})
}

func (c *simpleAuthenticationController) createUserFromToken(token *jwt.Token) (entities.User, error) {
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

	return &userRecord{
		Username:    username,
		Email:       email,
		DisplayName: displayName,
	}, nil
}
