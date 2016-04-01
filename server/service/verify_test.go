package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/clawio/service-auth/server/spec"
	"github.com/stretchr/testify/require"
)

func (suite *TestSuite) TestVerifyJSON() {
	testToken := "testtoken"
	testIdentity := &spec.Identity{
		Username:    "test",
		Email:       "test@test.com",
		DisplayName: "Test",
	}
	suite.MockTokenStore.On("Verify", testToken).Once().Return(testIdentity, nil)
	r, err := http.NewRequest("GET", verifyURL+testToken, nil)
	require.Nil(suite.T(), err)
	w := httptest.NewRecorder()
	suite.Server.ServeHTTP(w, r)
	require.Equal(suite.T(), 200, w.Code)
	suite.MockUserStore.AssertExpectations(suite.T())
	suite.MockTokenStore.AssertExpectations(suite.T())
	verifyRes := &spec.VerifyResponse{}
	err = json.NewDecoder(w.Body).Decode(verifyRes)
	require.Nil(suite.T(), err)
	require.Equal(suite.T(), testIdentity.Username, verifyRes.Identity.Username)
}

func (suite *TestSuite) TestVerifyJSONInvalidToken() {
	testToken := "testtoken"
	suite.MockTokenStore.On("Verify", testToken).Once().Return(&spec.Identity{}, errors.New("test error"))
	r, err := http.NewRequest("GET", verifyURL+testToken, nil)
	require.Nil(suite.T(), err)
	w := httptest.NewRecorder()
	suite.Server.ServeHTTP(w, r)
	require.Equal(suite.T(), 400, w.Code)
	suite.MockUserStore.AssertExpectations(suite.T())
	suite.MockTokenStore.AssertExpectations(suite.T())
}
