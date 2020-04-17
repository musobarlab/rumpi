package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type env struct {
	HTTPPort     uint16
	WebsocketKey string
}

// Env public type of 'env' will hold all environment's app
var Env env

// LoadEnv will help initialize app's environment
func LoadEnv() error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}

	var ok bool

	if httpPort, err := strconv.Atoi(os.Getenv("HTTP_PORT")); err != nil {
		return errors.New("required HTTP_PORT environment")
	} else {
		Env.HTTPPort = uint16(httpPort)
	}

	Env.WebsocketKey, ok = os.LookupEnv("WEBSOCKET_KEY")
	if !ok {
		return errors.New("required WEBSOCKET_KEY environment")
	}

	return nil
}
