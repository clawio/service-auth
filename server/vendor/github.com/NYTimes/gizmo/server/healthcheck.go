package server

import (
	"io"
	"net/http"
)

// HealthCheckHandler is an interface used by SimpleServer and RPCServer
// to allow users to customize their service's health check. Start will be
// called just before server start up and the given ActivityMonitor should
// offer insite to the # of requests in flight, if needed.
// Stop will be called once the servers receive a kill signal.
type HealthCheckHandler interface {
	http.Handler
	Path() string
	Start(*ActivityMonitor) error
	Stop() error
}

// SimpleHealthCheck is a basic HealthCheckHandler implementation
// that _always_ returns with an "ok" status and shuts down immediately.
type SimpleHealthCheck struct {
	path string
}

// NewSimpleHealthCheck will return a new SimpleHealthCheck instance.
func NewSimpleHealthCheck(path string) *SimpleHealthCheck {
	return &SimpleHealthCheck{path: path}
}

// Path will return the configured status path to server on.
func (s *SimpleHealthCheck) Path() string {
	return s.path
}

// Start will do nothing.
func (s *SimpleHealthCheck) Start(monitor *ActivityMonitor) error {
	return nil
}

// Stop will do nothing and return nil.
func (s *SimpleHealthCheck) Stop() error {
	return nil
}

// ServeHTTP will always respond with "ok-"+server.Name.
func (s *SimpleHealthCheck) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := io.WriteString(w, "ok-"+Name); err != nil {
		LogWithFields(r).Warn("unable to write healthcheck response: ", err)
	}
}
