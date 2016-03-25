package service

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/NYTimes/gizmo/server"
	"github.com/Sirupsen/logrus"
	"github.com/clawio/codes"
	"github.com/clawio/service-auth/server/spec"
	"golang.org/x/net/context"
)

// Authenticate authenticates a user using a username and password.
func (s *Service) Authenticate(ctx context.Context, r *spec.AuthNRequest) (*spec.AuthNResponse, error) {
	res := &spec.AuthNResponse{}
	identity, err := s.UserStore.FindByCredentials(r.Username, r.Password)
	if err != nil {
		server.Log.WithFields(logrus.Fields{
			"error": err,
		}).Error("error finding user")
		return res, codes.NewAPIErr(codes.BadAuthenticationData)
	}
	token, err := s.TokenStore.Create(identity)
	if err != nil {
		server.Log.WithFields(logrus.Fields{
			"error": err,
		}).Error("error creating authn token")
		return res, codes.NewAPIErr(codes.BadAuthenticationData)
	}
	res.Token = token
	return res, err
}

// AuthenticateJSON handles the JSON call and forwards the request to Authenticate.
// It delegates the logic to Authenticate.
func (s *Service) AuthenticateJSON(r *http.Request) (int, interface{}, error) {
	authNRequest := &spec.AuthNRequest{}
	if r.Body == nil {
		return http.StatusInternalServerError, nil, errors.New("body cannot be nil")
	}
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
		return http.StatusBadRequest, nil, err
	}
	return http.StatusOK, res, nil
}
