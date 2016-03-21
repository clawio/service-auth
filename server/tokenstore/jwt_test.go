package tokenstore

import (
	"reflect"
	"testing"

	"github.com/clawio/service-auth/server/spec"
	"github.com/dgrijalva/jwt-go"
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

	// test for nil identity
	_, err = store.Create(nil)
	if err == nil {
		t.Fatal("create must fail when identity is nil")
	}
}

func TestVerify(t *testing.T) {
	store := NewJWTTokenStore("")
	jwtStore := store.(*JWTTokenStore)

	token, err := store.Create(identity)
	if err != nil {
		t.Fatal(err)
	}
	if token == "" {
		t.Fatal("token cannot be empty")
	}
	identityFromToken, err := store.Verify(token)
	if err != nil {
		t.Fatal(err)
	}
	if identityFromToken == nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(identity, identityFromToken) {
		t.Fatal("identity from create and verify different")
	}

	// test for invalid token
	identityFromToken, err = store.Verify("faketoken")
	if err == nil {
		t.Fatal("verify must fail on invalid tokens")
	}

	// test for invalid  signing key
	jwtStore.JWTKey = "otherkey"
	_, err = store.Verify(token)
	if err == nil {
		t.Fatal("verify must fail on wrong signing key")
	}
	jwtStore.JWTKey = "" // restore to valid key

	// test for spurious token - username is not a string
	badToken := jwt.New(jwt.GetSigningMethod(jwtStore.SigningMethod))
	badToken.Claims["username"] = 123 // it should be a string
	hash, err := badToken.SignedString([]byte(jwtStore.JWTKey))
	if err != nil {
		t.Fatal(err)
	}
	_, err = store.Verify(hash)
	if err == nil {
		t.Fatal("verify should fail when token username is not a string or is not set")
	}

	// test for spurious token - email is not a string
	badToken = jwt.New(jwt.GetSigningMethod(jwtStore.SigningMethod))
	badToken.Claims["username"] = "test"
	badToken.Claims["email"] = 123 // it should be a string
	hash, err = badToken.SignedString([]byte(jwtStore.JWTKey))
	if err != nil {
		t.Fatal(err)
	}
	_, err = store.Verify(hash)
	if err == nil {
		t.Fatal("verify should fail when token email is not a string or is not set")
	}

	// test for spurious token - display_name is not a string
	badToken = jwt.New(jwt.GetSigningMethod(jwtStore.SigningMethod))
	badToken.Claims["username"] = "test"
	badToken.Claims["email"] = "test@test.com" // it should be a string
	badToken.Claims["display_name"] = 123
	hash, err = badToken.SignedString([]byte(jwtStore.JWTKey))
	if err != nil {
		t.Fatal(err)
	}
	_, err = store.Verify(hash)
	if err == nil {
		t.Fatal("verify should fail when token display_name is not a string or is not set")
	}

	// test for invalid signature
	// it will panic
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("create must panic when invalid signing method")
		}

	}()
	jwtStore.SigningMethod = "fakesigningmethod"
	store.Create(identity)
	jwtStore.SigningMethod = "H256" // restore to previous signing method

}
