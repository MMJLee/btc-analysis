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

func buildJWT(requestMethod, requestHost, requestPath string) (string, error) {
	uri := fmt.Sprintf("%s %s%s", requestMethod, requestHost, requestPath)
	keySecret := strings.ReplaceAll(os.Getenv("COINBASE_API_KEY_SECRET"), "\\n", "\n")

	block, _ := pem.Decode([]byte(keySecret))
	if block == nil {
		return "", fmt.Errorf("BuildJWT-Could not decode private key")
	}

	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("BuildJWT-%w", err)
	}

	keyName := os.Getenv("COINBASE_API_KEY_NAME")
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
		return "", fmt.Errorf("BuildJWT-%w", err)
	}
	return jwtString, nil
}
