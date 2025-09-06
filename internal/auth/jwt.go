package auth

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWT ...
type JWT struct {
	private *rsa.PrivateKey
	keyFunc jwt.Keyfunc
	ttl     time.Duration
}

// NewJWT ...
func NewJWT(private *rsa.PrivateKey, public *rsa.PublicKey, ttl time.Duration) *JWT {
	return &JWT{
		private: private,
		keyFunc: func(token *jwt.Token) (any, error) {
			return public, nil
		},
		ttl: ttl,
	}
}

// Encode ...
func (j *JWT) Encode(p *Payload) (string, error) {
	signature, err := jwt.NewWithClaims(jwt.SigningMethodRS256, toClaims(p, j.ttl)).SignedString(j.private)
	if err != nil {
		return "", fmt.Errorf("jwt.SignedString: %w", err)
	}

	return signature, nil
}

// Decode ...
func (j *JWT) Decode(token string) (*Payload, error) {
	var customClaims claims

	if _, err := jwt.ParseWithClaims(token, &customClaims, j.keyFunc); err != nil {
		return nil, fmt.Errorf("jwt.ParseWithClaims: %w", err)
	}

	return fromClaims(customClaims), nil
}

type claims struct {
	Username string
	Roles    []string
	jwt.RegisteredClaims
}

func toClaims(p *Payload, ttl time.Duration) *claims {
	now := time.Now()
	return &claims{
		Username: p.Username,
		Roles:    p.Roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "rtc",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
}

func fromClaims(claims claims) *Payload {
	return &Payload{
		Username: claims.Username,
		Roles:    claims.Roles,
	}
}
