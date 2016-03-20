package service

import (
	"net/http"

	"github.com/NYTimes/gizmo/config"
	"github.com/NYTimes/gizmo/server"
	"github.com/Sirupsen/logrus"
	"github.com/clawio/service-auth/codes"
	"google.golang.org/grpc"
)

type (

	// RPCService will implement server.RPCService and
	// handle all requests to the server.
	RPCService struct {
		Config *Config
	}

	// Config is a struct to contain all the needed
	// configuration for our Service
	Config struct {
		Server *config.Server
	}
)

// NewRPCServicewill instantiate and return
// a new Service that implements server.Service.
func NewRPCService(cfg *Config) *RPCService {
	return &RPCService{Config: cfg}
}

// Prefix returns the string prefix used for all endpoints within
// this service.
func (s *RPCService) Prefix() string {
	return "/clawio/auth/v1"
}

// Service provides the RPCService with a description of the
// service to serve and the implementation.
func (s *RPCService) Service() (*grpc.ServiceDesc, interface{}) {
	return &_AuthN_serviceDesc, s
}

// Middleware provides an http.Handler hook wrapped around all requests.
// In this implementation, we're using a GzipHandler middleware to
// compress our responses.
func (s *RPCService) Middleware(h http.Handler) http.Handler {
	return h
}

// JSONMiddleware provides a JSONEndpoint hook wrapped around all requests.
// In this implementation, we're using it to provide application logging and to check errors
// and provide generic responses.
func (s *RPCService) JSONMiddleware(j server.JSONEndpoint) server.JSONEndpoint {
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

// JSONEndpoints is a listing of all endpoints available in the RPCService.
func (s *RPCService) JSONEndpoints() map[string]map[string]server.JSONEndpoint {
	return map[string]map[string]server.JSONEndpoint{
		"/authenticate": map[string]server.JSONEndpoint{
			"POST": s.AuthenticateJSON,
		},
		"/verify": map[string]server.JSONEndpoint{
			"POST": s.VerifyJSON,
		},
	}

}
