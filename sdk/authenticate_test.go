package sdk

import (
	"fmt"
	"net/http"

	"github.com/stretchr/testify/require"
)

func (suite *TestSuite) TestAuthenticate() {
	suite.Router.HandleFunc("/clawio/auth/v1/authenticate", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"token":"faketoken"}`)
	})
	token, err := suite.SDK.Authenticate("", "")
	require.Nil(suite.T(), err)
	require.Equal(suite.T(), token, "faketoken")
}

func (suite *TestSuite) TestAuthenticateInvalidJSONBody() {
	suite.Router.HandleFunc("/clawio/auth/v1/authenticate", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `thisisnotjson`)
	})
	_, err := suite.SDK.Authenticate("", "")
	require.NotNil(suite.T(), err)
}

func (suite *TestSuite) TestAuthenticateAPIError() {
	suite.Router.HandleFunc("/clawio/auth/v1/authenticate", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"message": "some api error", "code": 99}`)
	})
	_, err := suite.SDK.Authenticate("", "")
	require.NotNil(suite.T(), err)
}

func (suite *TestSuite) TestAuthenticateInvalidAPIError() {
	suite.Router.HandleFunc("/clawio/auth/v1/authenticate", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, ``)
	})
	_, err := suite.SDK.Authenticate("", "")
	require.NotNil(suite.T(), err)
}

func (suite *TestSuite) TestAuthenticateNetworkError() {
	suite.Server.Close()
	_, err := suite.SDK.Authenticate("", "")
	require.NotNil(suite.T(), err)
}
