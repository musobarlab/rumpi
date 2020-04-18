package middleware

import (
	"github.com/musobarlab/rumpi/config"
	"github.com/musobarlab/rumpi/pkg/jwt"
)

// Middleware model
type Middleware struct {
	jwtService         jwt.JwtService
	username, password string
}

// NewMiddleware create new middleware instance
func NewMiddleware() *Middleware {
	return &Middleware{
		jwtService: jwt.NewJWT(config.Config.PublicKey, config.Config.PrivateKey),
		username:   config.Config.BasicAuthUsername,
		password:   config.Config.BasicAuthPassword,
	}
}
