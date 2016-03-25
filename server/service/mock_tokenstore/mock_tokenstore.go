package mock_tokenstore

import (
	"github.com/clawio/service-auth/server/spec"
	"github.com/stretchr/testify/mock"
)

type MockTokenStore struct {
	mock.Mock
}

func (m *MockTokenStore) Create(identity *spec.Identity) (string, error) {
	args := m.Called(identity)
	return args.String(0), args.Error(1)
}
func (m *MockTokenStore) Verify(token string) (*spec.Identity, error) {
	args := m.Called(token)
	return args.Get(0).(*spec.Identity), args.Error(1)
}
