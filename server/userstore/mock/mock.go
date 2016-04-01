package mock

import (
	"github.com/clawio/service-auth/server/spec"
	"github.com/stretchr/testify/mock"
)

type MockUserStore struct {
	mock.Mock
}

func (m *MockUserStore) FindByCredentials(username, password string) (*spec.Identity, error) {
	args := m.Called(username, password)
	return args.Get(0).(*spec.Identity), args.Error(1)
}
