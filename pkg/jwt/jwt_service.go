package jwt

import (
	"context"
	"crypto/rsa"
	"time"

	jwtgo "github.com/golang-jwt/jwt/v4"

	"github.com/musobarlab/rumpi/pkg/shared"
)

// JwtService represent jwt service
type JwtService interface {
	Generate(ctx context.Context, payload *Claim, expired time.Duration) shared.Result
	Validate(ctx context.Context, token string) shared.Result
}

// JWT implementation from JwtService
type JWT struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

// NewJWT constructor
func NewJWT(publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) *JWT {
	return &JWT{
		publicKey:  publicKey,
		privateKey: privateKey,
	}
}

// Generate token
func (r *JWT) Generate(ctx context.Context, payload *Claim, expired time.Duration) shared.Result {
	now := time.Now()
	exp := now.Add(expired)

	var key interface{}
	var token = new(jwtgo.Token)
	if payload.Alg == HS256 {
		token = jwtgo.New(jwtgo.SigningMethodHS256)
		key = []byte(HS256Key)
	} else {
		token = jwtgo.New(jwtgo.SigningMethodRS256)
		key = r.privateKey
	}

	claims := jwtgo.MapClaims{
		"iss": "rumpi",
		"exp": exp.Unix(),
		"iat": now.Unix(),
		"sub": payload.Subject,
		"aud": "97b33193-43ff-4e58-9124-b3a9b9f72c34",
	}

	if payload.User.ID != "" {
		claims["id"] = payload.User.ID
	}

	if payload.User.Email != "" {
		claims["email"] = payload.User.Email
	}

	if payload.User.FullName != "" {
		claims["fullName"] = payload.User.FullName
	}

	token.Claims = claims

	tokenString, err := token.SignedString(key)
	if err != nil {
		return shared.Result{Error: err}

	}

	return shared.Result{Data: tokenString}
}

// Validate token
func (r *JWT) Validate(ctx context.Context, tokenString string) shared.Result {
	tokenParse, err := jwtgo.Parse(tokenString, func(token *jwtgo.Token) (interface{}, error) {
		checkAlg, _ := shared.GetValueFromContext(ctx, shared.ContextKey("tokenAlg")).(Alg)
		if checkAlg == HS256 {
			return []byte(HS256Key), nil
		}
		return r.publicKey, nil
	})

	var errToken error
	switch ve := err.(type) {
	case *jwtgo.ValidationError:
		if ve.Errors == jwtgo.ValidationErrorExpired {
			errToken = shared.ErrTokenExpired
		} else {
			errToken = shared.ErrTokenFormat
		}
	}

	if errToken != nil {
		return shared.Result{Error: errToken}
	}

	if !tokenParse.Valid {
		return shared.Result{Error: shared.ErrTokenFormat}
	}

	mapClaims, _ := tokenParse.Claims.(jwtgo.MapClaims)

	var tokenClaim Claim
	tokenClaim.Issuer, _ = mapClaims["iss"].(string)
	tokenClaim.Audience, _ = mapClaims["aud"].(string)
	tokenClaim.IssuedAt, _ = mapClaims["iat"].(int64)
	tokenClaim.ExpiredAt, _ = mapClaims["exp"].(int64)
	tokenClaim.Subject, _ = mapClaims["sub"].(string)
	tokenClaim.User.ID, _ = mapClaims["id"].(string)
	tokenClaim.User.Email, _ = mapClaims["email"].(string)
	tokenClaim.User.FullName, _ = mapClaims["fullName"].(string)

	return shared.Result{Data: &tokenClaim}
}
