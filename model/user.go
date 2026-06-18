package model

import (
	"time"

	"github.com/uptrace/bun"
)

type UserRole string

const (
	UserRoleAdmin   UserRole = "admin"
	UserRoleTeacher UserRole = "teacher"
	UserRoleStudent UserRole = "student"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:usr"`

	ID        int64     `bun:"id,pk,autoincrement" form:"id" json:"id" xml:"id"`
	Name      string    `bun:"name,unique,notnull" form:"name" json:"name" xml:"name"`
	Role      UserRole  `bun:"role,notnull" form:"role" json:"role" xml:"role"`
	Password  string    `bun:"password,notnull" form:"password" json:"password" xml:"password"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" form:"created_at" json:"created_at" xml:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" form:"updated_at" json:"updated_at" xml:"updated_at"`
}
