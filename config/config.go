package config

import (
	"crypto/rsa"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/musobarlab/rumpi/config/key"
)

type config struct {
	HTTPPort              uint16
	WebsocketKey          string
	AccessTokenExpired    time.Duration
	PasswordHashIteration uint16
	BasicAuthUsername     string
	BasicAuthPassword     string
	PublicKey             *rsa.PublicKey
	PrivateKey            *rsa.PrivateKey
}

// Config public type of 'env' will hold all environment's app
var Config config

// LoadConfig will help initialize app's environment
func LoadConfig() error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}

	var ok bool

	if privateKey, err := key.LoadPrivateKey(); err != nil {
		return err
	} else {
		Config.PrivateKey = privateKey
	}

	if publicKey, err := key.LoadPublicKey(); err != nil {
		return err
	} else {
		Config.PublicKey = publicKey
	}

	if httpPort, err := strconv.Atoi(os.Getenv("HTTP_PORT")); err != nil {
		return errors.New("required HTTP_PORT environment")
	} else {
		Config.HTTPPort = uint16(httpPort)
	}

	Config.WebsocketKey, ok = os.LookupEnv("WEBSOCKET_KEY")
	if !ok {
		return errors.New("required WEBSOCKET_KEY environment")
	}

	accessTokenExpired, ok := os.LookupEnv("ACCESS_TOKEN_EXPIRED")
	if !ok {
		return errors.New("required ACCESS_TOKEN_EXPIRED environment")
	}

	Config.AccessTokenExpired, err = time.ParseDuration(accessTokenExpired)
	if err != nil {
		return errors.New("invalid access token string")
	}

	if passwordHashIteration, err := strconv.Atoi(os.Getenv("PASSWORD_HASH_ITERATION")); err != nil {
		return errors.New("required PASSWORD_HASH_ITERATION environment")
	} else {
		Config.PasswordHashIteration = uint16(passwordHashIteration)
	}

	Config.BasicAuthUsername, ok = os.LookupEnv("BASIC_AUTH_USER")
	if !ok {
		return errors.New("required BASIC_AUTH_USER environment")
	}
	Config.BasicAuthPassword, ok = os.LookupEnv("BASIC_AUTH_PASS")
	if !ok {
		return errors.New("required BASIC_AUTH_PASS environment")
	}

	return nil
}
