package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/NYTimes/gizmo/server"
	"github.com/clawio/service-auth/codes"
	"golang.org/x/net/context"
)

// Authenticate authenticates the user using a username and password.
// Possible error codes: 215
func (s *RPCService) Authenticate(ctx context.Context, r *AuthNRequest) (*AuthNResponse, error) {
	var err error
	defer server.MonitorRPCRequest()(ctx, "Authenticate", err)
	res := &AuthNResponse{}
	if r.Username == r.Password && r.Username == "hugo" {
		res.Token = "mytoken"
		return res, nil
	}
	err = codes.NewAPIErr(codes.BadAuthenticationData)
	return res, err
}

// AuthenticateJSON authenticates a user with username and password.
// It delegates the logic to Authenticate.
// Possible HTTP codes: 500, 400
func (s *RPCService) AuthenticateJSON(r *http.Request) (int, interface{}, error) {
	authNRequest := &AuthNRequest{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	err = json.Unmarshal(body, authNRequest)
	if err != nil {
		return http.StatusBadRequest, nil, codes.NewAPIErr(codes.BadInputData)
	}
	res, err := s.Authenticate(
		context.Background(),
		authNRequest,
	)
	if err != nil {
		switch err := err.(type) {
		case *codes.APIErr:
			if err.Code == 215 {
				return http.StatusBadRequest, nil, err
			}
			return http.StatusInternalServerError, nil, err
		default:
			return http.StatusInternalServerError, nil, err
		}
	}
	return http.StatusOK, res, nil
}
