package userstore

import (
	"github.com/clawio/service-auth/server/spec"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type SQLUserStore struct {
	driver, dsn string
	db          *gorm.DB
}

func NewSQLUserStore(driver, dsn string) (UserStore, error) {
	db, err := gorm.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&userEntity{}).Error
	if err != nil {
		return nil, err
	}
	return &SQLUserStore{
		driver: driver,
		dsn:    dsn,
		db:     &db,
	}, nil
}

func (s *SQLUserStore) FindByCredentials(username, password string) (*spec.Identity, error) {
	rec := &userEntity{}
	err := s.db.Where("username=? AND password=?", username, password).First(rec).Error
	if err != nil {
		return nil, err
	}
	identity := &spec.Identity{
		Username:    rec.Username,
		Email:       rec.Email,
		DisplayName: rec.DisplayName,
	}
	return identity, nil
}

// TODO(labkode) set collation for table and column to utf8. The default is swedish
type userEntity struct {
	Username    string `gorm:"primary_key"`
	Email       string
	DisplayName string
	Password    string
}
