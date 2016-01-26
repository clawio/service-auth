package main

import (
	"github.com/clawio/service.auth/lib"
	pb "github.com/clawio/service.auth/proto"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	rus "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"time"
)

const (
	dirPerm = 0755
)

// debugLogger satisfies Gorm's logger interface
// so that we can log SQL queries at Logrus' debug level
type debugLogger struct{}

func (*debugLogger) Print(msg ...interface{}) {
	rus.Debug(msg)
}

type newServerParams struct {
	dsn               string
	driver            string
	sharedSecret      string
	signMethod        string
	maxSqlIdle        int
	maxSqlConcurrency int
}

func newServer(p *newServerParams) (*server, error) {
	db, err := newDB(p.driver, p.dsn)
	if err != nil {
		rus.Error(err)
		return nil, err
	}

	db.LogMode(true)
	db.SetLogger(&debugLogger{})
	db.DB().SetMaxIdleConns(p.maxSqlIdle)
	db.DB().SetMaxOpenConns(p.maxSqlConcurrency)

	err = db.AutoMigrate(&identityRecord{}).Error
	if err != nil {
		rus.Error(err)
		return nil, err
	}

	rus.Infof("automigration applied")
	s := &server{}
	s.p = p
	s.db = db
	return s, nil
}

type server struct {
	p  *newServerParams
	db *gorm.DB
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

	rec := &identityRecord{}
	err := s.db.Where("pid=? AND password=?", r.Username, r.Password).First(rec).Error
	if err != nil {
		log.Error(err)
		return nil, grpc.Errorf(codes.Unauthenticated, "%s not found", r.Username)
	}

	idt := &lib.Identity{}
	idt.Pid = rec.Pid
	idt.Idp = rec.Idp
	idt.Email = rec.Email
	idt.DisplayName = rec.DisplayName

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
