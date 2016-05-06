package service

import (
	"net/http"
	"net/http/httptest"
	"path"
	"testing"

	"github.com/NYTimes/gizmo/config"
	"github.com/NYTimes/gizmo/server"
	mock_authenticationcontroller "github.com/clawio/authentication/authenticationcontroller/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	tokenURL   string
	metricsURL string
)

type TestSuite struct {
	suite.Suite
	MockAuthenticationController *mock_authenticationcontroller.AuthenticationController
	Service                      *Service
	Server                       *server.SimpleServer
}

func Test(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) SetupTest() {
	mockAuthenticationController := &mock_authenticationcontroller.AuthenticationController{}

	svc := &Service{}
	svc.AuthenticationController = mockAuthenticationController
	cfg := &Config{
		General: &GeneralConfig{BaseURL: "/"},
	}
	svc.Config = cfg

	suite.Service = svc
	suite.MockAuthenticationController = mockAuthenticationController

	serverCfg := &config.Server{}

	serv := server.NewSimpleServer(serverCfg)
	serv.Register(suite.Service)
	suite.Server = serv

	// set testing urls
	tokenURL = path.Join(svc.Config.General.BaseURL, "/token")
	metricsURL = path.Join(svc.Config.General.BaseURL, "/metrics")

}

func (suite *TestSuite) TestNew_withSimple() {
	authCfg := &AuthenticationControllerConfig{
		Type:                   "simple",
		SimpleDriver:           "sqlite3",
		SimpleDSN:              "/tmp/userstore.db",
		SimpleJWTKey:           "secret",
		SimpleJWTSigningMethod: "HS256",
	}
	cfg := &Config{
		General:                  &GeneralConfig{},
		AuthenticationController: authCfg,
	}
	_, err := New(cfg)
	require.Nil(suite.T(), err)
}

func (suite *TestSuite) TestNew_withMemory() {
	authCfg := &AuthenticationControllerConfig{
		Type: "memory",
	}
	cfg := &Config{
		General:                  &GeneralConfig{},
		AuthenticationController: authCfg,
	}
	_, err := New(cfg)
	require.Nil(suite.T(), err)
}
func (suite *TestSuite) TestNew_withBadController() {
	authCfg := &AuthenticationControllerConfig{
		Type: "notfound",
	}
	cfg := &Config{
		General:                  &GeneralConfig{},
		AuthenticationController: authCfg,
	}
	_, err := New(cfg)
	require.NotNil(suite.T(), err)
}
func (suite *TestSuite) TestNew_withNilConfig() {
	_, err := New(nil)
	require.NotNil(suite.T(), err)
}
func (suite *TestSuite) TestNew_withNilGeneralConfig() {
	cfg := &Config{
		General: &GeneralConfig{},
	}
	_, err := New(cfg)
	require.NotNil(suite.T(), err)
}
func (suite *TestSuite) TestNew_withNilAuthenticationControllerConfig() {
	cfg := &Config{
		Server: nil,
		AuthenticationController: nil,
	}
	_, err := New(cfg)
	require.NotNil(suite.T(), err)
}
func (suite *TestSuite) TestNew_withNilWithBadDSN() {
	authCfg := &AuthenticationControllerConfig{
		Type:                   "simple",
		SimpleDriver:           "sqlite3",
		SimpleDSN:              "/this/does/not/exists/userstore.db",
		SimpleJWTKey:           "secret",
		SimpleJWTSigningMethod: "HS256",
	}
	cfg := &Config{
		Server: nil,
		AuthenticationController: authCfg,
	}
	_, err := New(cfg)
	require.NotNil(suite.T(), err)
}
func (suite *TestSuite) TestMetrics() {
	r, err := http.NewRequest("GET", metricsURL, nil)
	require.Nil(suite.T(), err)
	w := httptest.NewRecorder()
	suite.Server.ServeHTTP(w, r)
	require.Equal(suite.T(), http.StatusOK, w.Code)
}
