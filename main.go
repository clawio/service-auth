package main

import (
	"fmt"
	pb "github.com/clawio/service-auth/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
	"runtime"
	"strconv"
)

const (
	serviceID              = "CLAWIO_AUTH"
	driverEnvar            = serviceID + "_DBDRIVER"
	dsnEnvar               = serviceID + "_DBDSN"
	signMethodEnvar        = serviceID + "_SIGNMETHOD"
	portEnvar              = serviceID + "_PORT"
	maxSqlIdleEnvar        = serviceID + "_MAXSQLIDLE"
	maxSqlConcurrencyEnvar = serviceID + "_MAXSQLCONCURRENCY"
	logLevelEnvar          = serviceID + "_LOGLEVEL"
	sharedSecretEnvar      = "CLAWIO_SHAREDSECRET"
)

type environ struct {
	dsn               string
	driver            string
	port              int
	signMethod        string
	maxSqlIdle        int
	maxSqlConcurrency int
	logLevel          string
	sharedSecret      string
}

func getEnviron() (*environ, error) {
	e := &environ{}
	e.driver = os.Getenv(driverEnvar)
	e.dsn = os.Getenv(dsnEnvar)
	e.signMethod = os.Getenv(signMethodEnvar)
	e.logLevel = os.Getenv(logLevelEnvar)
	port, err := strconv.Atoi(os.Getenv(portEnvar))
	if err != nil {
		return nil, err
	}
	e.port = port
	maxSqlIdle, err := strconv.Atoi(os.Getenv(maxSqlIdleEnvar))
	if err != nil {
		return nil, err
	}
	e.maxSqlIdle = maxSqlIdle

	maxSqlConcurrency, err := strconv.Atoi(os.Getenv(maxSqlConcurrencyEnvar))
	if err != nil {
		return nil, err
	}
	e.maxSqlConcurrency = maxSqlConcurrency
	e.sharedSecret = os.Getenv(sharedSecretEnvar)
	return e, nil
}

func printEnviron(e *environ) {
	log.Infof("%s=%s", driverEnvar, e.driver)
	log.Infof("%s=%s", dsnEnvar, e.dsn)
	log.Infof("%s=%s", signMethodEnvar, e.signMethod)
	log.Infof("%s=%d", portEnvar, e.port)
	log.Infof("%s=%d", maxSqlIdleEnvar, e.maxSqlIdle)
	log.Infof("%s=%d", maxSqlConcurrencyEnvar, e.maxSqlConcurrency)
	log.Infof("%s=%s", sharedSecretEnvar, "******")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	env, err := getEnviron()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	l, err := log.ParseLevel(env.logLevel)
	if err != nil {
		l = log.ErrorLevel
	}
	log.SetLevel(l)

	log.Infof("Service %s started", serviceID)

	printEnviron(env)

	p := &newServerParams{}
	p.dsn = env.dsn
	p.driver = env.driver
	p.sharedSecret = env.sharedSecret
	p.signMethod = env.signMethod
	p.maxSqlIdle = env.maxSqlIdle
	p.maxSqlConcurrency = env.maxSqlConcurrency

	srv, err := newServer(p)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", env.port))
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServer(grpcServer, srv)
	grpcServer.Serve(lis)
}
