package token

import (
	"errors"
	"time"
)

// Different types of error returned by the VerifyToken function
var (
	ErrInvalidToken         = errors.New("token is invalid")
	ErrExpiredToken         = errors.New("token has expired")
	AuthorizationPayloadKey = "auth_payload"
	TokenTypeAuth           = "auth"
	TokenTypeRefresh        = "refresh"
	TokenTypeResetPassword  = "reset_password"
	TokenTypeVerifyEmail    = "verify_email"
)

// Payload contains the payload data of the token
type Payload struct {
	UserID    int64     `json:"user_id"`
	ClientID  int64     `json:"client_id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Type      string    `json:"token_type"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewPayload creates a new auth token payload
func NewPayload(clientID, userID int64, tokenType, email string, role string, duration time.Duration) (*Payload, error) {
	payload := &Payload{
		ClientID:  clientID,
		UserID:    userID,
		Email:     email,
		Role:      role,
		Type:      tokenType,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

// Valid checks if the token payload is valid or not
func (payload *Payload) Valid(tokenType string) error {
	if payload.Type != tokenType {
		return ErrInvalidToken
	}
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
