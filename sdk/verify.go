package sdk

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/clawio/codes"
	"github.com/clawio/service-auth/server/spec"
)

// Verify verifies if an issued authn token is valid. If it is valid returns
// the identity obtained from it.
func (s *SDK) Verify(token string) (*spec.Identity, error) {
	url := fmt.Sprintf("%s/verify/%s", s.BaseURL, token)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		// Convert error to ApiErr
		apiErr := &codes.APIErr{}
		if err := json.NewDecoder(resp.Body).Decode(apiErr); err != nil {
			return nil, err
		}
		return nil, apiErr
	}
	verifyResponse := &spec.VerifyResponse{}
	if err := json.NewDecoder(resp.Body).Decode(verifyResponse); err != nil {
		return nil, err
	}
	return verifyResponse.Identity, nil
}
