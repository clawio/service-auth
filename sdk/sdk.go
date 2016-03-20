package sdk

import (
	"net/url"
)

// SDK provides a set of methods to interact with an authn server.
type SDK struct {
	authServerURL *url.URL
}

// NewSDK returns a new SDK.
func NewSDK(serverAddr string) (*SDK, error) {
	u, err := url.Parse(serverAddr)
	if err != nil {
		return nil, err
	}
	return &SDK{u}, nil
}
