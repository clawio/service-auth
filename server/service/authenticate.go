package service

import (
	"encoding/json"
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
		return res, codes.NewErr(codes.BadAuthenticationData, "user not found")
	}
	token, err := s.TokenStore.Create(identity)
	if err != nil {
		server.Log.WithFields(logrus.Fields{
			"error": err,
		}).Error("error creating authn token")
		return res, codes.NewErr(codes.BadAuthenticationData, "cannot create token")
	}
	res.Token = token
	return res, err
}

// AuthenticateJSON handles the JSON call and forwards the request to Authenticate.
// It delegates the logic to Authenticate.
func (s *Service) AuthenticateFunc(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	authNRequest := &spec.AuthNRequest{}
	if err := json.NewDecoder(r.Body).Decode(authNRequest); err != nil {
		e := codes.NewErr(codes.BadInputData, "")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(e)
		return
	}
	res, err := s.Authenticate(
		context.Background(),
		authNRequest,
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
