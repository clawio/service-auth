package main

import (
	"github.com/clawio/service.auth/lib"
	pb "github.com/clawio/service.auth/proto"
	"github.com/dgrijalva/jwt-go"
	rus "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"time"
)

const (
	dirPerm = 0755
)

type newServerParams struct {
	dsn          string
	driver       string
	sharedSecret string
	signMethod   string
}

func newServer(p *newServerParams) *server {
	// TODO: connect to some backend with GORM
	return &server{p}
}

type server struct {
	p *newServerParams
}

func (s *server) Authenticate(ctx context.Context, r *pb.AuthRequest) (*pb.AuthResponse, error) {
	traceID := getGRPCTraceID(ctx)
	log := rus.WithField("trace", traceID).WithField("svc", serviceID)
	ctx = newGRPCTraceContext(ctx, traceID)

	log.Info("request started")

	// Time request
	reqStart := time.Now()

	defer func() {
		// Compute request duration
		reqDur := time.Since(reqStart)

		// Log access info
		log.WithFields(rus.Fields{
			"method":   "authenticate",
			"type":     "grpcaccess",
			"duration": reqDur.Seconds(),
		}).Info("request finished")

	}()

	if r.Username != "demo" || r.Password != "demo" {
		return nil, grpc.Errorf(codes.Unauthenticated, "%s not found", r.Username)
	}

	idt := &lib.Identity{}
	idt.Pid = r.Username
	idt.Idp = "localhost"
	idt.Email = "me@me.com"
	idt.DisplayName = "Demo User"

	token, err := s.createToken(idt)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	res := &pb.AuthResponse{}
	res.Token = token
	log.Infof("authentication successful for %s", idt)
	return res, nil
}

// It returns the JWT token or an error.
func (s *server) createToken(idt *lib.Identity) (string, error) {

	token := jwt.New(jwt.GetSigningMethod(s.p.signMethod))
	token.Claims["pid"] = idt.Pid
	token.Claims["idp"] = idt.Idp
	token.Claims["display_name"] = idt.DisplayName
	token.Claims["email"] = idt.Email
	token.Claims["iss"] = "localhost"
	token.Claims["exp"] = time.Now().Add(time.Second * 3600).UnixNano()

	tokenStr, err := token.SignedString([]byte(s.p.sharedSecret))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
