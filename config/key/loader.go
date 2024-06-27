package key

import (
	"crypto/rsa"
	"io/ioutil"

	"github.com/golang-jwt/jwt/v4"
)

const (
	privateKeyPath = "config/key/app.key"
	publicKeyPath  = "config/key/app.key.pub"
)

var (
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
)

// LoadPublicKey return *rsa.PublicKey
func LoadPublicKey() (*rsa.PublicKey, error) {
	verifyBytes, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return nil, err
	}
	return publicKey, nil
}

// LoadPrivateKey return *rsa.PrivateKey
func LoadPrivateKey() (*rsa.PrivateKey, error) {
	signBytes, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}
