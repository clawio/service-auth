package sdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/clawio/service-auth/service"
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

func (s *SDK) Authenticate(username, password string) (string, error) {
	authNRequest := &service.AuthNRequest{username, password}
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
		return "", errors.New(string(resBody))
	}

	authNResponse := &service.AuthNResponse{}
	err = json.Unmarshal(resBody, authNResponse)
	if err != nil {
		return "", err
	}
	return authNResponse.Token, nil
}
