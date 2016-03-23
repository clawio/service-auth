package service

import (
	"errors"
	"net/http"

	"github.com/NYTimes/gizmo/config"
	"github.com/NYTimes/gizmo/server"
	"github.com/Sirupsen/logrus"
	"github.com/clawio/codes"
	"github.com/clawio/service-auth/server/tokenstore"
	"github.com/clawio/service-auth/server/userstore"
)

type (

	// Service will implement server.Service and
	// handle all requests to the server.
	Service struct {
		Config     *Config
		UserStore  userstore.UserStore
		TokenStore tokenstore.TokenStore
	}

	// Config is a struct to contain all the needed
	// configuration for our Service
	Config struct {
		Server *config.Server
		Auth   *AuthConfig
	}

	// AuthConfig represents the configuration for this Auth server
	AuthConfig struct {
		JWTKey             string
		SQLUserStoreDriver string
		SQLUserStoreDSN    string
	}
)

// New will instantiate and return
// a new Service that implements server.Service.
func New(cfg *Config) (*Service, error) {
	//TODO(labkode) Add flexibility based on config to choose both user and token store.
	if cfg == nil {
		return nil, errors.New("config cannot be nil")
	}
	tokenStore := tokenstore.NewJWTTokenStore(cfg.Auth.JWTKey)
	userStore, err := userstore.NewSQLUserStore(cfg.Auth.SQLUserStoreDriver, cfg.Auth.SQLUserStoreDSN)
	if err != nil {
		return nil, err
	}
	return &Service{
		Config:     cfg,
		UserStore:  userStore,
		TokenStore: tokenStore,
	}, nil
}

// Prefix returns the string prefix used for all endpoints within
// this service.
func (s *Service) Prefix() string {
	return "/clawio/auth/v1"
}

// Middleware provides an http.Handler hook wrapped around all requests.
// In this implementation, we're using a GzipHandler middleware to
// compress our responses.
func (s *Service) Middleware(h http.Handler) http.Handler {
	return h
}

// JSONMiddleware provides a JSONEndpoint hook wrapped around all requests.
// In this implementation, we're using it to provide application logging and to check errors
// and provide generic responses.
func (s *Service) JSONMiddleware(j server.JSONEndpoint) server.JSONEndpoint {
	return func(r *http.Request) (int, interface{}, error) {

		status, res, err := j(r)
		// Convert non apiErr to apiErr with code 500 and message "unexpected error"
		if err != nil {
			switch err.(type) {
			case *codes.APIErr:
				// we do nothing
			default:
				server.LogWithFields(r).WithFields(logrus.Fields{
					"error": err,
				}).Error("unexpected error serving request")
				err = codes.NewAPIErr(codes.Internal)
			}
		}

		server.LogWithFields(r).Info("request served")
		return status, res, err
	}

}

// JSONEndpoints is a listing of all endpoints available in the Service.
func (s *Service) JSONEndpoints() map[string]map[string]server.JSONEndpoint {
	return map[string]map[string]server.JSONEndpoint{
		"/authenticate": map[string]server.JSONEndpoint{
			"POST": s.AuthenticateJSON,
		},
		"/verify/{token}": map[string]server.JSONEndpoint{
			"GET": s.VerifyJSON,
		},
	}

}
