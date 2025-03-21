package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type Payload struct {
	ID            uuid.UUID `json:"id"`
	PublicAddress string    `json:"public_address"`
	IssueedAt     time.Time `json:"issued_at"`
	ExpiredAt     time.Time `json:"expired_at"`
}

func NewPayload(publicaddress string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:            tokenID,
		PublicAddress: publicaddress,
		IssueedAt:     time.Now(),
		ExpiredAt:     time.Now().Add(duration),
	}

	return payload, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
