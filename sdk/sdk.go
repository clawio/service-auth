package sdk

import (
	"net/http"
)

// SDK provides a set of methods to interact with an authn server.
type SDK struct {
	BaseURL    string
	HTTPClient *http.Client
}

// New returns a new SDK. If a nil httpClient is provided, http.DefaultClient will be used.
func New(baseURL string, httpClient *http.Client) *SDK {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &SDK{baseURL, httpClient}
}
