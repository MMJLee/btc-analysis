package client

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type APIKeyClaims struct {
	*jwt.RegisteredClaims
	URI string `json:"uri"`
}

func BuildJWT(requestMethod, requestHost, requestPath string) (string, error) {
	uri := fmt.Sprintf("%s %s%s", requestMethod, requestHost, requestPath)
	key_secret := strings.ReplaceAll(os.Getenv("COINBASE_API_KEY_SECRET"), "\\n", "\n")

	block, _ := pem.Decode([]byte(key_secret))
	if block == nil {
		return "", fmt.Errorf("jwt: Could not decode private key")
	}

	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("jwt: %w", err)
	}

	key_name := os.Getenv("COINBASE_API_KEY_NAME")
	claims := &APIKeyClaims{
		&jwt.RegisteredClaims{
			Issuer:    "cdp",
			Subject:   key_name,
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Minute)),
		}, uri,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["typ"] = "JWT"
	token.Header["kid"] = key_name

	jwtString, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("jwt: %w", err)
	}
	return jwtString, nil
}
