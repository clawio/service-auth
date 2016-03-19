package service

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/NYTimes/gizmo/server"
	"golang.org/x/net/context"
)

func (s *RPCService) Authenticate(ctx context.Context, r *AuthNRequest) (*AuthNResponse, error) {
	var err error
	defer server.MonitorRPCRequest()(ctx, "Authenticate", err)
	res := &AuthNResponse{}
	if r.Username == r.Password && r.Username == "hugo" {
		res.Token = "mytoken"
		return res, nil
	}
	return res, errors.New("username/password does not match")
}

func (s *RPCService) AuthenticateJSON(r *http.Request) (int, interface{}, error) {
	authNRequest := &AuthNRequest{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	err = json.Unmarshal(body, authNRequest)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	res, err := s.Authenticate(
		context.Background(),
		authNRequest,
	)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	return http.StatusOK, res.Token, nil
}
