package auth

import (
	"time"

	"github.com/google/uuid"
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Data      string    `json:"data"`
	IssuedAt  time.Time `json:"iat"`
	ExpiredAt time.Time `json:"exp"`
}

func (x *Payload) Valid() error {
	if time.Now().After(x.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

func NewPayload(data string, exp time.Duration) (*Payload, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	x := &Payload{
		id,
		data,
		time.Now(),
		time.Now().Add(exp),
	}
	return x, nil
}
