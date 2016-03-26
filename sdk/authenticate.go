package sdk

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/clawio/codes"
	"github.com/clawio/service-auth/server/spec"
)

// Authenticate authenticates a user using a username and a password.
func (s *SDK) Authenticate(username, password string) (string, error) {
	authNRequest := &spec.AuthNRequest{username, password}
	jsonBody, _ := json.Marshal(authNRequest)
	resp, err := http.Post(s.BaseURL+"/authenticate", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		apiErr := &codes.APIErr{}
		if err := json.NewDecoder(resp.Body).Decode(apiErr); err != nil {
			return "", err
		}
		return "", apiErr
	}
	authNResponse := &spec.AuthNResponse{}
	if err := json.NewDecoder(resp.Body).Decode(authNResponse); err != nil {
		return "", err
	}
	return authNResponse.Token, nil
}
