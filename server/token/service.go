package token

import (
	"github.com/dgrijalva/jwt-go"
)

type Service interface {
	NewToken(claims jwt.Claims) (string, error)
	ParseToken(newToken string) (jwt.Claims, error)
}

type JwtService struct {
	algorithm  *jwt.SigningMethodECDSA
	privateKey interface{}
	publicKey  interface{}
}

type jwtValidationError struct {
	message string
}

func (e *jwtValidationError) Error() string {
	return e.message
}

func NewJwtService(privateKeyData []byte, publicKeyData []byte) (*JwtService, error) {
	privateKey, err := jwt.ParseECPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseECPublicKeyFromPEM(publicKeyData)
	if err != nil {
		return nil, err
	}

	return &JwtService{
		jwt.SigningMethodES256,
		privateKey,
		publicKey,
	}, nil
}

func (g *JwtService) NewToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(g.algorithm, claims)

	return token.SignedString(g.privateKey)
}

func (g *JwtService) ParseToken(newToken string) (jwt.Claims, error) {
	claims := jwt.StandardClaims{}
	parsedToken, err := jwt.ParseWithClaims(
		newToken,
		&claims,
		func(token *jwt.Token) (interface{}, error) {
			return g.publicKey, nil
		},
	)

	switch {
	case err == nil && parsedToken.Valid:
		return claims, nil
	case err != nil:
		return nil, err
	default:
		return claims, &jwtValidationError{"Invalid token claims"}
	}
}
