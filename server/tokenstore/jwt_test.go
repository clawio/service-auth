package tokenstore

import (
	"reflect"
	"testing"

	"github.com/clawio/service-auth/server/spec"
)

var identity = &spec.Identity{
	"test",
	"test@test.com",
	"Test",
}

func TestCreate(t *testing.T) {
	store := NewJWTTokenStore("")
	token, err := store.Create(identity)
	if err != nil {
		t.Fatal(err)
	}
	if token == "" {
		t.Fatal("token is empty")
	}
}

func TestVerify(t *testing.T) {
	store := NewJWTTokenStore("")
	token, err := store.Create(identity)
	if err != nil {
		t.Fatal(err)
	}
	if token == "" {
		t.Fatal("token is empty")
	}
	identityFromToken, err := store.Verify(token)
	if err != nil {
		t.Fatal(err)
	}
	if identityFromToken == nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(identity, identityFromToken) {
		t.Fatal("identity from token is different")
	}

	// Test for invalid token
	identityFromToken, err = store.Verify("faketoken")
	if err == nil {
		t.Fatal("verify must fail on invalid tokens")
	}
}
