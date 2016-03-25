package service

import (
	"testing"

	"github.com/NYTimes/gizmo/config"
	"github.com/NYTimes/gizmo/server"
	"github.com/clawio/service-auth/server/service/mock_tokenstore"
	"github.com/clawio/service-auth/server/service/mock_userstore"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
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
	require.Equal(suite.T(), suite.Service.Prefix(), "/clawio/auth/v1", "prefix must be equal")
}
