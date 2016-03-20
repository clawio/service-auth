package tokenstore

import (
	"fmt"
	"os"
	"time"

	"github.com/clawio/service-auth/server/spec"
	"github.com/dgrijalva/jwt-go"
)

// JWTTokenStore implements the TokenStore using JWT Tokens.
// See http://jwt.io/
type JWTTokenStore struct {
	JWTKey string // the key for signing the token
}

// NewJWTTokenStore returns a TokenStore.
func NewJWTTokenStore(key string) TokenStore {
	return &JWTTokenStore{key}
}

// Create creates a JWT Token from an user Identity.
func (j *JWTTokenStore) Create(identity *spec.Identity) (string, error) {
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	host, _ := os.Hostname()
	token.Claims["username"] = identity.Username
	token.Claims["email"] = identity.Email
	token.Claims["display_name"] = identity.DisplayName
	token.Claims["iss"] = host
	token.Claims["exp"] = time.Now().Add(time.Second * 3600).UnixNano()
	hash, err := token.SignedString([]byte(j.JWTKey))
	if err != nil {
		return "", err
	}
	return hash, nil
}

// Verify checks id token is valid and creates an user identity from it.
func (j *JWTTokenStore) Verify(token string) (*spec.Identity, error) {
	t, err := j.parseToken(token)
	if err != nil {
		return nil, err
	}
	identity, err := j.createIdentityFromToken(t)
	if err != nil {
		return nil, err
	}
	return identity, nil
}

func (j *JWTTokenStore) parseToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (key interface{}, err error) {
		return []byte(j.JWTKey), nil
	})
}

func (j *JWTTokenStore) createIdentityFromToken(token *jwt.Token) (*spec.Identity, error) {
	username, ok := token.Claims["username"].(string)
	if !ok {
		return nil, fmt.Errorf("token username claim failed cast to string")
	}

	email, ok := token.Claims["email"].(string)
	if !ok {
		return nil, fmt.Errorf("token email claim failed cast to string")
	}

	displayName, ok := token.Claims["display_name"].(string)
	if !ok {
		return nil, fmt.Errorf("token display_name claim failed cast to string")
	}

	return &spec.Identity{
		Username:    username,
		Email:       email,
		DisplayName: displayName,
	}, nil
}
