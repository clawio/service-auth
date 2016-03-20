package userstore

import "github.com/clawio/service-auth/server/spec"

type UserStore interface {
	FindByCredentials(username, password string) (*spec.Identity, error)
}
