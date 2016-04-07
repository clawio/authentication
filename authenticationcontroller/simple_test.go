package authenticationcontroller

import (
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	authenticationController       AuthenticationController
	simpleAuthenticationController *simpleAuthenticationController
}

func Test(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
func (suite *TestSuite) SetupTest() {
	opts := &SimpleAuthenticationControllerOptions{
		Driver:           "sqlite3",
		DSN:              "/tmp/userstore.db",
		JWTKey:           "secret",
		JWTSigningMethod: "HS256",
	}
	authenticationController, err := NewSimpleAuthenticationController(opts)
	require.Nil(suite.T(), err)
	suite.authenticationController = authenticationController
	suite.simpleAuthenticationController = suite.authenticationController.(*simpleAuthenticationController)
}
func (suite *TestSuite) TeardownTest() {
	os.RemoveAll("/tmp/t")
	os.RemoveAll("/tmp/userstore.db")
}
func (suite *TestSuite) TestNewSimpleAuthenticationController() {
	opts := &SimpleAuthenticationControllerOptions{
		Driver:           "sqlite3",
		DSN:              "/tmp/userstore.db",
		JWTKey:           "secret",
		JWTSigningMethod: "HS256",
	}
	_, err := NewSimpleAuthenticationController(opts)
	require.Nil(suite.T(), err)
}
func (suite *TestSuite) TestNewSimpleAuthenticationController_withBadDriver() {
	opts := &SimpleAuthenticationControllerOptions{
		Driver:           "thisnotexists",
		DSN:              "/tmp/userstore.db",
		JWTKey:           "secret",
		JWTSigningMethod: "HS256",
	}
	_, err := NewSimpleAuthenticationController(opts)
	require.NotNil(suite.T(), err)
}

func (suite *TestSuite) TestfindByCredentials() {
	db, err := sql.Open(suite.simpleAuthenticationController.driver, suite.simpleAuthenticationController.dsn)
	require.Nil(suite.T(), err)
	defer db.Close()
	sqlStmt := `insert into users values ("testFindByCredentials", "test@test.com", "Test", "testpwd")`
	_, err = db.Exec(sqlStmt)
	defer db.Exec("delete from users")
	require.Nil(suite.T(), err)
	user, err := suite.simpleAuthenticationController.findByCredentials("testFindByCredentials", "testpwd")
	require.Nil(suite.T(), err)
	require.Equal(suite.T(), "testFindByCredentials", user.GetUsername())
}
func (suite *TestSuite) TestfindByCredentials_withBadUser() {
	_, err := suite.simpleAuthenticationController.findByCredentials("", "")
	require.NotNil(suite.T(), err)
}
func (suite *TestSuite) TestcreateToken() {
	user := &userRecord{}
	_, err := suite.simpleAuthenticationController.createToken(user)
	require.Nil(suite.T(), err)
}
func (suite *TestSuite) TestcreateToken_withNilUser() {
	_, err := suite.simpleAuthenticationController.createToken(nil)
	require.NotNil(suite.T(), err)
}
func (suite *TestSuite) TestparseToken_withBadToken() {
	_, err := suite.simpleAuthenticationController.parseToken("")
	require.NotNil(suite.T(), err)
}
func (suite *TestSuite) TestparseToken() {
	user := &userRecord{}
	token, err := suite.simpleAuthenticationController.createToken(user)
	require.Nil(suite.T(), err)
	_, err = suite.simpleAuthenticationController.parseToken(token)
	require.Nil(suite.T(), err)
}
func (suite *TestSuite) TestcreateUserFromToken() {
	user := &userRecord{}
	token, err := suite.simpleAuthenticationController.createToken(user)
	require.Nil(suite.T(), err)
	jwtToken, err := suite.simpleAuthenticationController.parseToken(token)
	require.Nil(suite.T(), err)
	_, err = suite.simpleAuthenticationController.createUserFromToken(jwtToken)
	require.Nil(suite.T(), err)
}
func (suite *TestSuite) TestcreateUserFromToken_withBadUsername() {
	user := &userRecord{}
	token, err := suite.simpleAuthenticationController.createToken(user)
	require.Nil(suite.T(), err)
	jwtToken, err := suite.simpleAuthenticationController.parseToken(token)
	require.Nil(suite.T(), err)
	jwtToken.Claims["username"] = 0
	_, err = suite.simpleAuthenticationController.createUserFromToken(jwtToken)
	require.NotNil(suite.T(), err)
}
func (suite *TestSuite) TestcreateUserFromToken_withBadEmail() {
	user := &userRecord{}
	token, err := suite.simpleAuthenticationController.createToken(user)
	require.Nil(suite.T(), err)
	jwtToken, err := suite.simpleAuthenticationController.parseToken(token)
	require.Nil(suite.T(), err)
	jwtToken.Claims["email"] = 0
	_, err = suite.simpleAuthenticationController.createUserFromToken(jwtToken)
	require.NotNil(suite.T(), err)
}
func (suite *TestSuite) TestcreateUserFromToken_withBadDisplayName() {
	user := &userRecord{}
	token, err := suite.simpleAuthenticationController.createToken(user)
	require.Nil(suite.T(), err)
	jwtToken, err := suite.simpleAuthenticationController.parseToken(token)
	require.Nil(suite.T(), err)
	jwtToken.Claims["display_name"] = 0
	_, err = suite.simpleAuthenticationController.createUserFromToken(jwtToken)
	require.NotNil(suite.T(), err)
}
func (suite *TestSuite) TestAuthenticate() {
	db, err := sql.Open(suite.simpleAuthenticationController.driver, suite.simpleAuthenticationController.dsn)
	require.Nil(suite.T(), err)
	defer db.Close()
	sqlStmt := `insert into users values ("testAuthenticate", "test@test.com", "Test", "testpwd")`
	_, err = db.Exec(sqlStmt)
	require.Nil(suite.T(), err)
	defer db.Exec("delete from users where username=testAuthenticate")
	_, err = suite.simpleAuthenticationController.Authenticate("testAuthenticate", "testpwd")
	require.Nil(suite.T(), err)
}
func (suite *TestSuite) TestAuthenticate_withBadUser() {
	_, err := suite.simpleAuthenticationController.Authenticate("", "")
	require.NotNil(suite.T(), err)
}
func (suite *TestSuite) TestInvalidate() {
	err := suite.simpleAuthenticationController.Invalidate("")
	require.Nil(suite.T(), err)
}
func (suite *TestSuite) TestVerifyToken() {
	user := &userRecord{Username: "test"}
	token, err := suite.simpleAuthenticationController.createToken(user)
	require.Nil(suite.T(), err)
	givenUser, err := suite.simpleAuthenticationController.Verify(token)
	require.Nil(suite.T(), err)
	require.Equal(suite.T(), user.GetUsername(), givenUser.GetUsername())
}
func (suite *TestSuite) TestVerify_withBadToken() {
	_, err := suite.simpleAuthenticationController.Verify("")
	require.NotNil(suite.T(), err)
}
