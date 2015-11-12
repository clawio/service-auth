package main

import (
	"github.com/clawio/service.auth/lib"
	pb "github.com/clawio/service.auth/proto"
	"github.com/dgrijalva/jwt-go"
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
	if r.Username != "demo" || r.Password != "demo" {
		return nil, grpc.Errorf(codes.Unauthenticated, "entity %s not found", r.Username)
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
		log.Error(err)
		return "", err
	}

	return tokenStr, nil
}
