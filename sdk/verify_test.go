package sdk

import (
	"fmt"
	"net/http"

	"github.com/stretchr/testify/require"
)

func (suite *TestSuite) TestVerify() {
	suite.Router.HandleFunc("/clawio/auth/v1/verify/{token}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"identity": {"username": "test", "email": "test@test.com", "display_name":"Test"}}`)
	})
	identity, err := suite.SDK.Verify("test")
	require.Nil(suite.T(), err)
	require.Equal(suite.T(), identity.Username, "test")
	require.Equal(suite.T(), identity.Email, "test@test.com")
	require.Equal(suite.T(), identity.DisplayName, "Test")
}

func (suite *TestSuite) TestVerifyInvalidJSONBody() {
	suite.Router.HandleFunc("/clawio/auth/v1/verify/{token}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, ``)
	})
	_, err := suite.SDK.Verify("testtoken")
	require.NotNil(suite.T(), err)
}
func (suite *TestSuite) TestVerifyAPIError() {
	suite.Router.HandleFunc("/clawio/auth/v1/verify/{token}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"message": "some api error", "code": 99}`)
	})
	_, err := suite.SDK.Verify("testtoken")
	require.NotNil(suite.T(), err)
}
func (suite *TestSuite) TestVerifyInvalidAPIError() {
	suite.Router.HandleFunc("/clawio/auth/v1/verify/{token}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{{}}`)
	})
	_, err := suite.SDK.Verify("testtoken")
	require.NotNil(suite.T(), err)
}

func (suite *TestSuite) TestVerifyNetworkError() {
	suite.Server.Close()
	_, err := suite.SDK.Verify("")
	require.NotNil(suite.T(), err)
}
