package service

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NYTimes/gizmo/config"
	"github.com/NYTimes/gizmo/server"
	mock_tokenstore "github.com/clawio/service-auth/server/tokenstore/mock"
	mock_userstore "github.com/clawio/service-auth/server/userstore/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	baseURL         = "/clawio/v1/auth/"
	authenticateURL = baseURL + "authenticate"
	verifyURL       = baseURL + "verify/"
)

type TestSuite struct {
	suite.Suite
	MockUserStore  *mock_userstore.MockUserStore
	MockTokenStore *mock_tokenstore.MockTokenStore
	Service        *Service
	Server         *server.SimpleServer
}

func Test(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) SetupTest() {
	mockUserStore := &mock_userstore.MockUserStore{}
	mockTokenStore := &mock_tokenstore.MockTokenStore{}
	svc := &Service{}
	svc.TokenStore = mockTokenStore
	svc.UserStore = mockUserStore
	suite.Service = svc
	suite.MockUserStore = mockUserStore
	suite.MockTokenStore = mockTokenStore
	cfg := &config.Server{}
	cfg.LogLevel = "fatal"
	serv := server.NewSimpleServer(cfg)
	serv.Register(suite.Service)
	suite.Server = serv
}

func (suite *TestSuite) TestNew() {
	authCfg := &AuthConfig{
		JWTKey:             "",
		SQLUserStoreDriver: "sqlite3",
		SQLUserStoreDSN:    "/tmp/sqliteuserstore.db",
	}
	cfg := &Config{
		Server: nil,
		Auth:   authCfg,
	}
	_, err := New(cfg)
	require.Nil(suite.T(), err)
}
func (suite *TestSuite) TestNewNilConfig() {
	_, err := New(nil)
	require.NotNil(suite.T(), err)
}

func (suite *TestSuite) TestNewInvalidUserStore() {
	authCfg := &AuthConfig{
		JWTKey:             "",
		SQLUserStoreDriver: "",
		SQLUserStoreDSN:    "",
	}
	cfg := &Config{
		Server: nil,
		Auth:   authCfg,
	}
	_, err := New(cfg)
	require.NotNil(suite.T(), err)
}

func (suite *TestSuite) TestPrefix() {
	require.Equal(suite.T(), suite.Service.Prefix(), "/clawio/v1/auth", "prefix must be equal")
}

func (suite *TestSuite) TestMetrics() {
	r, err := http.NewRequest("GET", "/clawio/v1/auth/metrics", nil)
	require.Nil(suite.T(), err)
	w := httptest.NewRecorder()
	suite.Server.ServeHTTP(w, r)
	require.Equal(suite.T(), 200, w.Code)
}
