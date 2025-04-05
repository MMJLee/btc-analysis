package api

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mmjlee/btc-analysis/internal/util"
)

type APIKeyClaims struct {
	*jwt.RegisteredClaims
	URI string `json:"uri"`
}

func BuildJWT(requestMethod, requestHost, requestPath string) (string, error) {
	uri := fmt.Sprintf("%s %s%s", requestMethod, requestHost, requestPath)

	// block, _ := pem.Decode([]byte(os.Getenv("COINBASE_API_KEY_SECRET")))
	block, _ := pem.Decode([]byte(util.COINBASE_API_KEY_SECRET))
	if block == nil {
		return "", fmt.Errorf("jwt: Could not decode private key")
	}

	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("jwt: %w", err)
	}

	claims := &APIKeyClaims{
		&jwt.RegisteredClaims{
			Issuer: "cdp",
			// Subject:   os.Getenv("COINBASE_API_KEY_NAME"),
			Subject:   util.COINBASE_API_KEY_NAME,
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Minute)),
		}, uri,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["typ"] = "JWT"
	// token.Header["kid"] = os.Getenv("COINBASE_API_KEY_NAME")
	token.Header["kid"] = util.COINBASE_API_KEY_NAME

	jwtString, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("jwt: %w", err)
	}
	return jwtString, nil
}
