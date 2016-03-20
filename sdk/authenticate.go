package sdk

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/clawio/codes"
	"github.com/clawio/service-auth/server/spec"
)

func (s *SDK) Authenticate(username, password string) (string, error) {
	authNRequest := &spec.AuthNRequest{username, password}
	jsonBody, err := json.Marshal(authNRequest)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", s.authServerURL.String()+"/authenticate", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		// Convert error to ApiErr
		apiErr := &codes.APIErr{}
		err := json.Unmarshal(resBody, apiErr)
		if err != nil {
			return "", err
		}
		return "", apiErr
	}

	authNResponse := &spec.AuthNResponse{}
	err = json.Unmarshal(resBody, authNResponse)
	if err != nil {
		return "", err
	}
	return authNResponse.Token, nil
}
