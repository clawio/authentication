package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/clawio/entities/mocks"
	"github.com/stretchr/testify/require"
)

func (suite *TestSuite) TestVerify() {
	testUser := &mocks.MockUser{Username: "test", Email: "test@test.com", DisplayName: "Tester"}
	suite.MockAuthenticationController.On("Verify").Once().Return(testUser, nil)
	r, err := http.NewRequest("GET", verifyURL+"testoken", nil)
	require.Nil(suite.T(), err)
	w := httptest.NewRecorder()
	suite.Server.ServeHTTP(w, r)
	require.Equal(suite.T(), http.StatusOK, w.Code)
	user := &mocks.MockUser{}
	err = json.NewDecoder(w.Body).Decode(user)
	require.Nil(suite.T(), err)
	require.Equal(suite.T(), testUser.Username, user.Username)
}
func (suite *TestSuite) TestVerify_witAuthenticationControllerError() {
	suite.MockAuthenticationController.On("Verify").Once().Return(&mocks.MockUser{}, errors.New("test error"))
	r, err := http.NewRequest("GET", verifyURL+"testoken", nil)
	require.Nil(suite.T(), err)
	w := httptest.NewRecorder()
	suite.Server.ServeHTTP(w, r)
	require.Equal(suite.T(), http.StatusBadRequest, w.Code)
}
