package sdk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/clawio/codes"
	"github.com/clawio/service-auth/server/spec"
)

func (s *SDK) Verify(token string) (*spec.Identity, error) {
	u := s.authServerURL
	u.Path = fmt.Sprintf("%s/verify/%s", u.Path, token)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		// Convert error to ApiErr
		apiErr := &codes.APIErr{}
		err := json.Unmarshal(resBody, apiErr)
		if err != nil {
			return nil, err
		}
		return nil, apiErr
	}

	verifyResponse := &spec.VerifyResponse{}
	err = json.Unmarshal(resBody, verifyResponse)
	if err != nil {
		return nil, err
	}
	return verifyResponse.Identity, nil
}
