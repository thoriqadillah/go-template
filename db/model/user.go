package model

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`
	Id            string     `json:"id" bun:"id,pk"`
	Email         string     `json:"email" bun:"email,unique"`
	Password      *string    `json:"-" bun:"password"`
	Name          string     `json:"name" bun:"name"`
	Source        string     `json:"source" bun:"source"`
	ResetToken    string     `json:"-" bun:"reset_token"`
	VerifiedAt    *time.Time `json:"verifiedAt" bun:"verified_at"`
	CreatedAt     time.Time  `json:"createdAt" bun:"created_at"`
	UpdatedAt     time.Time  `json:"updatedAt" bun:"updated_at"`
}
