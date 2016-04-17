package memory

import (
	"testing"

	"github.com/clawio/authentication/authenticationcontroller"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var users = []*User{
	{Username: "test", Password: "test"},
	{Username: "hugo", Password: "hugo"},
}

type TestSuite struct {
	suite.Suite
	authenticationController authenticationcontroller.AuthenticationController
	controller               *controller
}

func Test(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
func (suite *TestSuite) SetupTest() {
	opts := &Options{
		Users:            users,
		JWTKey:           "secret",
		JWTSigningMethod: "HS256",
	}
	authenticationController := New(opts)
	require.NotNil(suite.T(), authenticationController)
	suite.authenticationController = authenticationController
	suite.controller = suite.authenticationController.(*controller)
}

func (suite *TestSuite) TestNew() {
	opts := &Options{
		Users:            users,
		JWTKey:           "secret",
		JWTSigningMethod: "HS256",
	}
	c := New(opts)
	require.NotNil(suite.T(), c)
}

func (suite *TestSuite) TestAuthenticate() {
	_, err := suite.authenticationController.Authenticate("test", "test")
	require.Nil(suite.T(), err)
}
func (suite *TestSuite) TestAuthenticate_withBadUser() {
	_, err := suite.authenticationController.Authenticate("notfound", "notfound")
	require.NotNil(suite.T(), err)
}
func (suite *TestSuite) TestcreateToken() {
	User := &User{}
	_, err := suite.controller.createToken(User)
	require.Nil(suite.T(), err)
}
func (suite *TestSuite) TestcreateToken_withNilUser() {
	_, err := suite.controller.createToken(nil)
	require.NotNil(suite.T(), err)
}
func (suite *TestSuite) TestparseToken_withBadToken() {
	_, err := suite.controller.parseToken("")
	require.NotNil(suite.T(), err)
}
func (suite *TestSuite) TestparseToken() {
	User := &User{}
	token, err := suite.controller.createToken(User)
	require.Nil(suite.T(), err)
	_, err = suite.controller.parseToken(token)
	require.Nil(suite.T(), err)
}
func (suite *TestSuite) TestcreateUserFromToken() {
	User := &User{}
	token, err := suite.controller.createToken(User)
	require.Nil(suite.T(), err)
	jwtToken, err := suite.controller.parseToken(token)
	require.Nil(suite.T(), err)
	_, err = suite.controller.createUserFromToken(jwtToken)
	require.Nil(suite.T(), err)
}
func (suite *TestSuite) TestcreateUserFromToken_withBadUsername() {
	User := &User{}
	token, err := suite.controller.createToken(User)
	require.Nil(suite.T(), err)
	jwtToken, err := suite.controller.parseToken(token)
	require.Nil(suite.T(), err)
	jwtToken.Claims["username"] = 0
	_, err = suite.controller.createUserFromToken(jwtToken)
	require.NotNil(suite.T(), err)
}
func (suite *TestSuite) TestcreateUserFromToken_withBadEmail() {
	User := &User{}
	token, err := suite.controller.createToken(User)
	require.Nil(suite.T(), err)
	jwtToken, err := suite.controller.parseToken(token)
	require.Nil(suite.T(), err)
	jwtToken.Claims["email"] = 0
	_, err = suite.controller.createUserFromToken(jwtToken)
	require.NotNil(suite.T(), err)
}
func (suite *TestSuite) TestcreateUserFromToken_withBadDisplayName() {
	User := &User{}
	token, err := suite.controller.createToken(User)
	require.Nil(suite.T(), err)
	jwtToken, err := suite.controller.parseToken(token)
	require.Nil(suite.T(), err)
	jwtToken.Claims["display_name"] = 0
	_, err = suite.controller.createUserFromToken(jwtToken)
	require.NotNil(suite.T(), err)
}
func (suite *TestSuite) TestInvalidate() {
	err := suite.controller.Invalidate("")
	require.Nil(suite.T(), err)
}
func (suite *TestSuite) TestVerifyToken() {
	User := &User{Username: "test"}
	token, err := suite.controller.createToken(User)
	require.Nil(suite.T(), err)
	givenUser, err := suite.controller.Verify(token)
	require.Nil(suite.T(), err)
	require.Equal(suite.T(), User.GetUsername(), givenUser.GetUsername())
}
func (suite *TestSuite) TestVerify_withBadToken() {
	_, err := suite.controller.Verify("")
	require.NotNil(suite.T(), err)
}
