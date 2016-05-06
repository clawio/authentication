package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/stretchr/testify/require"
)

func (suite *TestSuite) TestAuthenticate() {
	suite.MockAuthenticationController.On("Authenticate").Once().Return("testtoken", nil)
	body := strings.NewReader(`{"username":"test", "password":"test"}`)
	r, err := http.NewRequest("POST", tokenURL, body)
	require.Nil(suite.T(), err)
	w := httptest.NewRecorder()
	suite.Server.ServeHTTP(w, r)
	require.Equal(suite.T(), http.StatusOK, w.Code)
	authNRes := &AuthenticateResponse{}
	err = json.NewDecoder(w.Body).Decode(authNRes)
	require.Nil(suite.T(), err)
	require.Equal(suite.T(), "testtoken", authNRes.AccessToken)
}
func (suite *TestSuite) TestAuthenticate_withNilBody() {
	r, err := http.NewRequest("POST", tokenURL, nil)
	require.Nil(suite.T(), err)
	w := httptest.NewRecorder()
	suite.Server.ServeHTTP(w, r)
	require.Equal(suite.T(), http.StatusInternalServerError, w.Code)
}
func (suite *TestSuite) TestAuthenticate_withInvalidJSON() {
	body := strings.NewReader("")
	r, err := http.NewRequest("POST", tokenURL, body)
	require.Nil(suite.T(), err)
	w := httptest.NewRecorder()
	suite.Server.ServeHTTP(w, r)
	require.Equal(suite.T(), http.StatusBadRequest, w.Code)
}
func (suite *TestSuite) TestAuthenticate_withAuthenticationControllerError() {
	suite.MockAuthenticationController.On("Authenticate").Once().Return("", errors.New("test error"))
	body := strings.NewReader(`{"username":"test", "password":"test"}`)
	r, err := http.NewRequest("POST", tokenURL, body)
	require.Nil(suite.T(), err)
	w := httptest.NewRecorder()
	suite.Server.ServeHTTP(w, r)
	require.Equal(suite.T(), http.StatusBadRequest, w.Code)
}
