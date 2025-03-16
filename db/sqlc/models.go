// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Session struct {
	ID            uuid.UUID          `json:"id"`
	PublicAddress string             `json:"public_address"`
	RefreshToken  string             `json:"refresh_token"`
	UserAgent     string             `json:"user_agent"`
	ClientIp      string             `json:"client_ip"`
	IsBlocked     bool               `json:"is_blocked"`
	ExpiresAt     time.Time          `json:"expires_at"`
	CreatedAt     pgtype.Timestamptz `json:"created_at"`
}

type User struct {
	ID            int32            `json:"id"`
	PublicAddress string           `json:"public_address"`
	Nonce         string           `json:"nonce"`
	CreatedAt     pgtype.Timestamp `json:"created_at"`
	LastLogin     pgtype.Timestamp `json:"last_login"`
}
