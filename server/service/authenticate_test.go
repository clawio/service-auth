package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/clawio/service-auth/server/spec"
	"github.com/stretchr/testify/require"
)

func (suite *TestSuite) TestAuthenticateJSON() {
	testIdentity := &spec.Identity{
		Username:    "test",
		Email:       "test@test.com",
		DisplayName: "Test",
	}
	suite.MockUserStore.On("FindByCredentials", "test", "test").Once().Return(testIdentity, nil)
	suite.MockTokenStore.On("Create", testIdentity).Once().Return("testtoken", nil)
	body := strings.NewReader(`{"username":"test", "password":"test"}`)
	r, err := http.NewRequest("POST", "/clawio/auth/v1/authenticate", body)
	require.Nil(suite.T(), err)
	w := httptest.NewRecorder()
	suite.Server.ServeHTTP(w, r)
	require.Equal(suite.T(), w.Code, 200)
	suite.MockUserStore.AssertExpectations(suite.T())
	suite.MockTokenStore.AssertExpectations(suite.T())
	authNRes := &spec.AuthNResponse{}
	err = json.NewDecoder(w.Body).Decode(authNRes)
	require.Nil(suite.T(), err)
	require.Equal(suite.T(), authNRes.Token, "testtoken")
}

func (suite *TestSuite) TestAuthenticateJSONNilBody() {
	r, err := http.NewRequest("POST", "/clawio/auth/v1/authenticate", nil)
	require.Nil(suite.T(), err)
	w := httptest.NewRecorder()
	suite.Server.ServeHTTP(w, r)
	require.Equal(suite.T(), w.Code, 500)
}

func (suite *TestSuite) TestAuthenticateJSONInvalidJSON() {
	body := strings.NewReader("")
	r, err := http.NewRequest("POST", "/clawio/auth/v1/authenticate", body)
	require.Nil(suite.T(), err)
	w := httptest.NewRecorder()
	suite.Server.ServeHTTP(w, r)
	require.Equal(suite.T(), w.Code, 400)
}

func (suite *TestSuite) TestAuthenticateJSONUserNotFound() {
	suite.MockUserStore.On("FindByCredentials", "unexistent", "unexistent").Return(&spec.Identity{}, errors.New("test error"))
	body := strings.NewReader(`{"username":"unexistent", "password":"unexistent"}`)
	r, err := http.NewRequest("POST", "/clawio/auth/v1/authenticate", body)
	require.Nil(suite.T(), err)
	w := httptest.NewRecorder()
	suite.Server.ServeHTTP(w, r)
	require.Equal(suite.T(), w.Code, 400)
}

func (suite *TestSuite) TestAuthenticateJSONTokenCreationError() {
	testIdentity := &spec.Identity{
		Username:    "test",
		Email:       "test@test.com",
		DisplayName: "Test",
	}
	suite.MockUserStore.On("FindByCredentials", "test", "test").Once().Return(testIdentity, nil)
	suite.MockTokenStore.On("Create", testIdentity).Once().Return("", errors.New("test error"))
	body := strings.NewReader(`{"username":"test", "password":"test"}`)
	r, err := http.NewRequest("POST", "/clawio/auth/v1/authenticate", body)
	require.Nil(suite.T(), err)
	w := httptest.NewRecorder()
	suite.Server.ServeHTTP(w, r)
	require.Equal(suite.T(), w.Code, 400)
	suite.MockUserStore.AssertExpectations(suite.T())
	suite.MockTokenStore.AssertExpectations(suite.T())
}