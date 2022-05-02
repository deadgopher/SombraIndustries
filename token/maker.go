package token

import (
	"errors"
	"germ/model"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
)

type TokenMaker interface {
	CreateToken(*model.User, time.Duration) (string, error)
	VerifyToken(string) (*Payload, error)
}

type JWTMaker struct {
	secret string
}

func (x *JWTMaker) CreateToken(data *model.User, exp time.Duration) (string, error) {

	payload, err := NewPayload(data, exp)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(x.secret))

}

func (x *JWTMaker) VerifyToken(t string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(x.secret), nil
	}

	token, err := jwt.ParseWithClaims(t, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	payload, ok := token.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}

func NewJWTMaker(secret string) *JWTMaker {
	return &JWTMaker{secret}
}
