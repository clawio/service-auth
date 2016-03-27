package service

import (
	"errors"
	"net/http"

	"github.com/NYTimes/gizmo/config"
	"github.com/clawio/service-auth/server/tokenstore"
	"github.com/clawio/service-auth/server/userstore"
	"github.com/prometheus/client_golang/prometheus"
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

// Endpoints is a listing of all endpoints available in the MixedService.
func (s *Service) Endpoints() map[string]map[string]http.HandlerFunc {
	return map[string]map[string]http.HandlerFunc{
		"/metrics": map[string]http.HandlerFunc{
			"GET": func(w http.ResponseWriter, r *http.Request) {
				prometheus.Handler().ServeHTTP(w, r)
			},
		},
		"/authenticate": map[string]http.HandlerFunc{
			"POST": prometheus.InstrumentHandlerFunc("/authenticate", s.AuthenticateFunc),
		},
		"/verify/{token}": map[string]http.HandlerFunc{
			"GET": prometheus.InstrumentHandlerFunc("/verify", s.VerifyFunc),
		},
	}
}
