package service

import (
	"net/http"

	"github.com/NYTimes/gizmo/server"
	"golang.org/x/net/context"
)

func (s *RPCService) Verify(ctx context.Context, r *VerifyRequest) (*VerifyResponse, error) {
	var err error
	defer server.MonitorRPCRequest()(ctx, "Verify", err)
	res := &VerifyResponse{}
	idt := &Identity{}
	idt.Username = "mytoken"
	res.Identity = idt
	return res, nil
}

func (s *RPCService) VerifyJSON(r *http.Request) (int, interface{}, error) {
	var err error
	res, err := s.Verify(
		context.Background(),
		&VerifyRequest{
			"mytoken",
		})
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	return http.StatusOK, res, nil
}
