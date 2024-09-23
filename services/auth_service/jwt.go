package auth_service

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTClaims struct {
	Type      string `json:"type"`
	Issuer    string `json:"issuer"`
	IssuedAt  int64  `json:"issued_at"`
	ExpiresAt int64  `json:"expires_at"`
	UserId    string `json:"user_id"`
	Role      string `json:"role"`
}

func (c *JWTClaims) Valid() error {
	if c.ExpiresAt < time.Now().Unix() {
		return fmt.Errorf("token is expired")
	}
	return nil
}

func IsValidRole(role string, config JWTConfig) bool {
	for _, r := range config.Roles {
		if r == role {
			return true
		}
	}
	return false
}

func IsValidType(tokenType string, config JWTConfig) bool {
	for _, t := range config.Types {
		if t == tokenType {
			return true
		}
	}
	return false
}

func (c *JWTClaims) ValidateWithConfig(config JWTConfig) error {
	vErr := new(jwt.ValidationError)
	now := time.Now().Unix()

	if c.Issuer != config.Issuer {
		vErr.Inner = fmt.Errorf("invalid issuer")
		vErr.Errors |= jwt.ValidationErrorIssuer
	}
	if c.IssuedAt == 0 {
		vErr.Inner = fmt.Errorf("missing issuedAt")
		vErr.Errors |= jwt.ValidationErrorIssuedAt
	}
	if (c.IssuedAt - now) > 0 {
		vErr.Inner = fmt.Errorf("token used before issued")
		vErr.Errors |= jwt.ValidationErrorIssuedAt
	}
	if c.ExpiresAt == 0 {
		vErr.Inner = fmt.Errorf("missing expiresAt")
		vErr.Errors |= jwt.ValidationErrorExpired
	}
	if now > c.ExpiresAt {
		vErr.Inner = fmt.Errorf("token is expired")
		vErr.Errors |= jwt.ValidationErrorExpired
	}
	if !IsValidRole(c.Role, config) {
		vErr.Inner = fmt.Errorf("invalid role")
		vErr.Errors |= jwt.ValidationErrorClaimsInvalid
	}
	if !IsValidType(c.Type, config) {
		vErr.Inner = fmt.Errorf("invalid type")
		vErr.Errors |= jwt.ValidationErrorClaimsInvalid
	}

	if vErr.Errors > 0 {
		return vErr
	}

	return nil
}

func (s *AuthService) CreateJWTWithClaims(tokenType string, duration time.Duration, userId string, role string) (string, error) {
	claims := &JWTClaims{
		Type:      tokenType,
		Issuer:    s.JWTConfig.Issuer,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(duration).Unix(),
		UserId:    userId,
		Role:      role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(s.JWTConfig.Secret)

	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (s *AuthService) ParseJWT(tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.JWTConfig.Secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		err := claims.ValidateWithConfig(s.JWTConfig)
		if err != nil {
			return nil, err
		}
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
