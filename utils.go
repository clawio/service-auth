package main

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"github.com/jinzhu/gorm"
	"golang.org/x/net/context"
	metadata "google.golang.org/grpc/metadata"
)

// TODO(labkode) set collation for table and column to utf8. The default is swedish
type identityRecord struct {
	Pid         string `gorm:"primary_key"`
	Idp         string
	Email       string
	DisplayName string
	Password    string
}

func (r *identityRecord) String() string {
	return fmt.Sprintf("pid=%s idp=%s email=%s display_name=%s password=***",
		r.Pid, r.Idp, r.Email, r.DisplayName)
}

func newDB(driver, dsn string) (*gorm.DB, error) {

	db, err := gorm.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&identityRecord{})

	return &db, nil
}

func newGRPCTraceContext(ctx context.Context, trace string) context.Context {
	md := metadata.Pairs("trace", trace)
	ctx = metadata.NewContext(ctx, md)
	return ctx

}

func getGRPCTraceID(ctx context.Context) string {

	md, ok := metadata.FromContext(ctx)
	if !ok {
		return uuid.New()

	}

	tokens := md["trace"]
	if len(tokens) == 0 {
		return uuid.New()

	}

	if tokens[0] != "" {
		return tokens[0]

	}

	return uuid.New()

}
