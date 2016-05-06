package lib

import (
	"testing"

	"github.com/clawio/entities"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var user = &entities.User{Username: "test"}

type TestSuite struct {
	suite.Suite
	authenticator *Authenticator
}

func Test(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
func (suite *TestSuite) SetupTest() {
	authenticator := NewAuthenticator("", "")
	suite.authenticator = authenticator
}

func (suite *TestSuite) TestNew() {
	authenticator := NewAuthenticator("", "")
	require.NotNil(suite.T(), authenticator)
}

func (suite *TestSuite) TestCreateToken() {
	_, err := suite.authenticator.CreateToken(user)
	require.Nil(suite.T(), err)
}
func (suite *TestSuite) TestCreateToken_withNilUser() {
	_, err := suite.authenticator.CreateToken(nil)
	require.NotNil(suite.T(), err)
}
func (suite *TestSuite) TestparseToken_withBadToken() {
	_, err := suite.authenticator.parseToken("")
	require.NotNil(suite.T(), err)
}
func (suite *TestSuite) TestparseToken() {
	token, err := suite.authenticator.CreateToken(user)
	require.Nil(suite.T(), err)
	_, err = suite.authenticator.parseToken(token)
	require.Nil(suite.T(), err)
}
func (suite *TestSuite) TestcreateUserFromToken() {
	token, err := suite.authenticator.CreateToken(user)
	require.Nil(suite.T(), err)
	_, err = suite.authenticator.CreateUserFromToken(token)
	require.Nil(suite.T(), err)
}
func (suite *TestSuite) TestcreateUserFromToken_withBadToken() {
	_, err := suite.authenticator.CreateUserFromToken("")
	require.NotNil(suite.T(), err)
}
func (suite *TestSuite) TestgetUserFromRawToken_withBadUsername() {
	token, err := suite.authenticator.CreateToken(user)
	require.Nil(suite.T(), err)
	jwtToken, err := suite.authenticator.parseToken(token)
	require.Nil(suite.T(), err)
	jwtToken.Claims["username"] = 0
	_, err = suite.authenticator.getUserFromRawToken(jwtToken)
	require.NotNil(suite.T(), err)
}
func (suite *TestSuite) TestgetUserFromRawToken_withBadEmail() {
	token, err := suite.authenticator.CreateToken(user)
	require.Nil(suite.T(), err)
	jwtToken, err := suite.authenticator.parseToken(token)
	require.Nil(suite.T(), err)
	jwtToken.Claims["email"] = 0
	_, err = suite.authenticator.getUserFromRawToken(jwtToken)
	require.NotNil(suite.T(), err)
}
func (suite *TestSuite) TestgetUserFromRawToken_withBadDisplayName() {
	token, err := suite.authenticator.CreateToken(user)
	require.Nil(suite.T(), err)
	jwtToken, err := suite.authenticator.parseToken(token)
	require.Nil(suite.T(), err)
	jwtToken.Claims["display_name"] = 0
	_, err = suite.authenticator.getUserFromRawToken(jwtToken)
	require.NotNil(suite.T(), err)
}
