package model

import (
	"time"

	"github.com/uptrace/bun"
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusBanned   UserStatus = "banned"
)

type Profile struct {
	bun.BaseModel `bun:"table:profiles,alias:prf"`

	ID     int64 `bun:"id,pk,autoincrement" form:"id" json:"id" xml:"id"`
	UserID int64 `bun:"user_id,notnull" form:"user_id" json:"user_id" xml:"user_id"`

	Photo    string         `bun:"photo,notnull" form:"photo" json:"photo" xml:"photo"`
	Name     string         `bun:"name,notnull" form:"name" json:"name" xml:"name"`
	Email    string         `bun:"email,unique,notnull" form:"email" json:"email" xml:"email"`
	Settings map[string]any `bun:",type:jsonb"` // PostgreSQL JSONB
	Status   UserStatus     `bun:",type:varchar(20),default:'active'"`

	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	// DeletedAt time.Time `bun:",soft_delete,nullzero"` // Soft delete support

	// Relationship
	User *User `bun:"rel:belongs-to,join:user_id=id"`
}
