package service

import (
	"encoding/json"
	"net/http"

	"github.com/NYTimes/gizmo/server"
	"github.com/Sirupsen/logrus"
	"github.com/clawio/codes"
	"github.com/clawio/service-auth/server/spec"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

// Verify verifies that an issued authn token is valid. If it is valid an identity is obtained from
// it and returned.
func (s *Service) Verify(ctx context.Context, r *spec.VerifyRequest) (*spec.VerifyResponse, error) {
	res := &spec.VerifyResponse{}
	identity, err := s.TokenStore.Verify(r.Token)
	if err != nil {
		server.Log.WithFields(logrus.Fields{
			"error": err,
		})
		return res, codes.NewAPIErr(codes.InvalidToken)
	}
	res.Identity = identity
	return res, nil
}

// VerifyJSON handles the JSON call and forwards the request to Authenticate.
func (s *Service) VerifyFunc(w http.ResponseWriter, r *http.Request) {
	res, err := s.Verify(
		context.Background(),
		&spec.VerifyRequest{
			mux.Vars(r)["token"],
		},
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		apiErr := codes.NewAPIErr(codes.BadInputData)
		json.NewEncoder(w).Encode(apiErr)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
