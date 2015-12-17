package main

import (
	"fmt"
	pb "github.com/clawio/service.auth/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
	"runtime"
	"strconv"
)

const (
	serviceID         = "CLAWIO_AUTH"
	driverEnvar       = serviceID + "_DBDRIVER"
	dsnEnvar          = serviceID + "_DBDSN"
	signMethodEnvar   = serviceID + "_SIGNMETHOD"
	portEnvar         = serviceID + "_PORT"
	sharedSecretEnvar = "CLAWIO_SHAREDSECRET"
)

type environ struct {
	dsn          string
	driver       string
	port         int
	signMethod   string
	sharedSecret string
}

func getEnviron() (*environ, error) {
	e := &environ{}
	e.dsn = os.Getenv(dsnEnvar)
	e.signMethod = os.Getenv(signMethodEnvar)
	port, err := strconv.Atoi(os.Getenv(portEnvar))
	if err != nil {
		return nil, err
	}
	e.port = port
	e.sharedSecret = os.Getenv(sharedSecretEnvar)
	return e, nil
}

func printEnviron(e *environ) {
	log.Infof("%s=%s", dsnEnvar, e.dsn)
	log.Infof("%s=%s", signMethodEnvar, e.signMethod)
	log.Infof("%s=%d", portEnvar, e.port)
	log.Infof("%s=%s", sharedSecretEnvar, "******")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Infof("Service %s started", serviceID)

	env, err := getEnviron()
	printEnviron(env)

	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	p := &newServerParams{}
	p.dsn = env.dsn
	p.driver = env.driver
	p.sharedSecret = env.sharedSecret
	p.signMethod = env.signMethod
	srv := newServer(p)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", env.port))
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServer(grpcServer, srv)
	grpcServer.Serve(lis)
}
