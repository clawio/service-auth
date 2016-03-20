package sdk

import (
	"net/url"
)

type SDK struct {
	authServerURL *url.URL
}

func NewSDK(serverAddr string) (*SDK, error) {
	u, err := url.Parse(serverAddr)
	if err != nil {
		return nil, err
	}
	return &SDK{u}, nil
}
