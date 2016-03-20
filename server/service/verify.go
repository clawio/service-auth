package service

import (
	"net/http"

	"github.com/NYTimes/gizmo/server"
	"github.com/Sirupsen/logrus"
	"github.com/clawio/codes"
	"github.com/clawio/service-auth/server/spec"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

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

func (s *Service) VerifyJSON(r *http.Request) (int, interface{}, error) {
	res, err := s.Verify(
		context.Background(),
		&spec.VerifyRequest{
			mux.Vars(r)["token"],
		},
	)
	if err != nil {
		switch err := err.(type) {
		case *codes.APIErr:
			if err.Code == codes.InvalidToken {
				return http.StatusBadRequest, nil, err
			}
			return http.StatusInternalServerError, nil, err
		default:
			return http.StatusInternalServerError, nil, err
		}
	}
	return http.StatusOK, res, nil
}

func isValidToken(ctx context.Context, token string) bool {
	if token == "mytoken" {
		return true
	}
	return false
}
