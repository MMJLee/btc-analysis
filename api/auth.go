package api

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type APIKeyClaims struct {
	*jwt.RegisteredClaims
	URI string `json:"uri"`
}

func BuildJWT(requestMethod, requestHost, requestPath string) (string, error) {
	uri := fmt.Sprintf("%s %s%s", requestMethod, requestHost, requestPath)

	block, _ := pem.Decode([]byte(keySecret))
	if block == nil {
		return "", fmt.Errorf("jwt: Could not decode private key")
	}

	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("jwt: %w", err)
	}

	claims := &APIKeyClaims{
		&jwt.RegisteredClaims{
			Issuer:    "cdp",
			Subject:   keyName,
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Minute)),
		}, uri,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["typ"] = "JWT"
	token.Header["kid"] = keyName

	jwtString, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("jwt: %w", err)
	}
	return jwtString, nil
}
